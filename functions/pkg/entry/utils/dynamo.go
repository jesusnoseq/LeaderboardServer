package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/config"
)

func NewDynamoClient() *dynamodb.Client {
	dyn := dynamodb.NewFromConfig(config.DefaultAwsSession())
	return dyn
}

func NewLocalDynamoClient() *dynamodb.Client {
	return dynamodb.NewFromConfig(CreateLocalAWSConf())
}

func MustCreateTable(
	ctx context.Context,
	dbClient *dynamodb.Client,
	tableName string) *dynamodb.CreateTableOutput {
	out, err := dbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeB,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(tableName),
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		panic(fmt.Errorf("failed to create table %q: %w", tableName, err))
	}
	WaitForTable(ctx, dbClient, tableName)
	return out
}

func MustCreateTableSync(
	ctx context.Context,
	dbClient *dynamodb.Client,
	tableName string) *dynamodb.CreateTableOutput {
	out := MustCreateTable(ctx, dbClient, tableName)
	WaitForTable(ctx, dbClient, tableName)
	return out
}

func WaitForTable(ctx context.Context, db *dynamodb.Client, tableName string) error {
	w := dynamodb.NewTableExistsWaiter(db)
	err := w.Wait(ctx,
		&dynamodb.DescribeTableInput{
			TableName: aws.String(tableName),
		},
		2*time.Minute,
		func(o *dynamodb.TableExistsWaiterOptions) {
			o.MaxDelay = 5 * time.Second
			o.MinDelay = 5 * time.Second
		})
	if err != nil {
		return fmt.Errorf("timed out while waiting for table to become active: %w", err)
	}

	return err
}

func ListTables(ctx context.Context, dbClient *dynamodb.Client) []string {
	pages := dynamodb.NewListTablesPaginator(dbClient, nil, func(o *dynamodb.ListTablesPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})
	tables := make([]string, 0)
	for pages.HasMorePages() {
		out, err := pages.NextPage(ctx)
		if err != nil {
			panic(err)
		}

		for _, tn := range out.TableNames {
			tables = append(tables, tn)
			fmt.Println(tn)
		}
	}
	return tables
}
