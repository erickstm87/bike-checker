package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func readDB() {
    sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := dynamodb.New(sess)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Link": {
				S: aws.String("https://www.pinkbike.com/buysell/3150797/"),
			},
		},
		TableName: aws.String("bike-availability"),
	}
	result, err := svc.GetItem(input)
	if err != nil {
		fmt.Println("there was an error")
	}
	fmt.Println("here is your result: ", result.Item)
}
	
