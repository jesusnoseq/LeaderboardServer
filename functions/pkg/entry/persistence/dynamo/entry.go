package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func NewDBSession(awss *aws.Session) *dynamodb.DynamoDB {
	return dynamodb.New(awss)
}

func NewAWSSession(region string) (*aws.Session, error) {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	return awsSession, err
}
