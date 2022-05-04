package dynamo

import (
	"github.com/jesusnoseq/LeaderboardServer/functions/pkg/entry/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func DefaulDynamoClient() *dynamodb.Client {
	dyn := dynamodb.NewFromConfig(config.DefaultAwsSession())
	return dyn
}
