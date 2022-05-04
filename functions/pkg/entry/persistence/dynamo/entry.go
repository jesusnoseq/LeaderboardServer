package dynamo

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/models"
)

type EntryRepository struct {
	tableName string
	dbclient  *dynamodb.Client
	timeout   time.Duration
}

func NewEntryRepository(dbclient *dynamodb.Client, defaultTimeout time.Duration) *EntryRepository {
	e := EntryRepository{
		tableName: "LeaderboardTable",
		timeout:   defaultTimeout,
		dbclient:  dbclient,
	}
	return &e
}

// read one
func (e *EntryRepository) GetEntry(ctx context.Context, uuid string) (models.Entry, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: uuid},
		},
		TableName:            aws.String(e.tableName),
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String("id, name, score"),
	}
	getItemResponse, err := e.dbclient.GetItem(ctx, getItemInput)
	if err != nil {
		log.Fatalf("get work unit failed, %v", err)
	}
	if getItemResponse.Item == nil {
		log.Fatalf("item not found")
	}
	entry := models.Entry{}
	err = attributevalue.UnmarshalMap(getItemResponse.Item, &entry)
	if err != nil {
		log.Fatalf("unmarshal failed, %v", err)
	}

	return entry, nil
}

func (e *EntryRepository) GetEntries(ctx context.Context) ([]models.Entry, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()
	out, err := e.dbclient.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(e.tableName),
	})
	if err != nil {
		return nil, err
	}
	items := make([]models.Entry, out.Count)
	err = attributevalue.UnmarshalListOfMaps(out.Items, &items)
	return items, err
}

// read with filters
// https://aws.github.io/aws-sdk-go-v2/docs/code-examples/dynamodb/scanitems/#source-code

func (e *EntryRepository) CreateEntry(ctx context.Context, entry models.Entry) (models.Entry, error) {
	ctx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	scoreStr := strconv.FormatUint(uint64(entry.Score), 10)
	id, _ := uuid.New().MarshalBinary()
	out, err := e.dbclient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(e.tableName),
		Item: map[string]types.AttributeValue{
			"id":    &types.AttributeValueMemberB{Value: id},
			"name":  &types.AttributeValueMemberS{Value: entry.Name},
			"score": &types.AttributeValueMemberN{Value: scoreStr},
		},
	})
	if err != nil {
		return entry, err
	}
	entry = models.Entry{}
	err = attributevalue.UnmarshalMap(out.Attributes, &entry)

	return entry, err
}
