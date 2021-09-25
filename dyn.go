package main

import (
	"fmt"
	"strconv"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func readDB(seedData []string) {
    sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(Timezone),
	})
	svc := dynamodb.New(sess)
	for index, entry := range seedData {
		input := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Link": {
					S: aws.String(entry),
				},
			},
			TableName: aws.String(TableName),
		}
		result, err := svc.GetItem(input)
		if err != nil {
			fmt.Println("there was an error: ", err)
		}
		if(len(result.Item) == 0 && entry != "") {
			// alert me
			fmt.Println("new entry found: ", entry)
			alertMe(entry)
			updateDb(entry, strconv.Itoa(index))
		}
		fmt.Println("here is your result: ", result.Item)
	}
	
}

func updateDb(entry string, anIndex string) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(Timezone),
	})
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Link": {
				S: aws.String(entry),
			},
			"TheIndex": {
				S: aws.String(anIndex),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName: aws.String(TableName),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling UpdateItem: ", err)
	}
}

func alertMe(entry string) {
	fmt.Println("work in progress")
}
	
// func seedDB(seedData []string) {
// 	sess, _ := session.NewSession(&aws.Config{
// 		Region: aws.String(Timezone),
// 	})
// 	svc := dynamodb.New(sess)
// 	seedData = seedData[1:]
// 	for index, entry := range seedData {
// 		anIndex := strconv.Itoa(index)
// 		fmt.Println("here is the index stringified: ", anIndex)
// 		input := &dynamodb.PutItemInput{
// 			Item: map[string]*dynamodb.AttributeValue{
// 				"Link": {
// 					S: aws.String(entry),
// 				},
// 				"TheIndex": {
// 					S: aws.String(anIndex),
// 				},
// 			},
// 			ReturnConsumedCapacity: aws.String("TOTAL"),
// 			TableName: aws.String(TableName),
// 		}
// 		result, err := svc.PutItem(input)
// 		if err != nil {
// 			fmt.Println("I have an error", err)
// 			os.Exit(1)
// 		}
// 		fmt.Println(result)
// 	}
// }
