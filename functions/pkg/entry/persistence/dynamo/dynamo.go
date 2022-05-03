package dynamo

type Entry struct {
}

//LEADERBOARD_TABLE_NAME
type DynamoDAO struct {
	dbClient dynamodbiface.DynamoDBAPI
}

// sess := session.Must(session.NewSessionWithOptions(session.Options{
//     SharedConfigState: session.SharedConfigEnable,
// }))

// // Create DynamoDB client
// svc := dynamodb.New(sess)
