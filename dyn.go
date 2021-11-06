package main

import (
	"fmt"
	"strconv"
	"time"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/rest/api/v2010"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Item struct {
    TimeStamp string
    Interested string
    Model string
    Link string
	TheIndex string
}

func newSession() *dynamodb.DynamoDB{
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("TIMEZONE")),
	})
	svc := dynamodb.New(sess)
	return svc
}

func readDB(seedData []AvailableBike) {
	itemInterested := map[string]Item{}
    
	svc := newSession()
	for index, entry := range seedData {
        if(entry.link != "" || entry.link != " ") {
			input := &dynamodb.GetItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"Link": {
						S: aws.String(entry.link),
					},
				},
				TableName: aws.String(os.Getenv("TABLE_NAME")),
			}
			result, err := svc.GetItem(input)
			item := Item{}

			err = dynamodbattribute.UnmarshalMap(result.Item, &item)
			if err != nil {
				fmt.Println("Failed to unmarshal Record, ", err)
			}

			if(len(result.Item) == 0 && entry.link != "") {
				// alert me
				fmt.Println("new entry found: ", entry)
				alertMe(entry.link, entry.model, "new bike")
				updateDb(entry.link, entry.model, strconv.Itoa(index))
			} else if(item.Interested == "Yes") {
				itemInterested[item.Link] = item
			} 
		}
	}
	checkInterested(seedData)
}

func checkInterested(seedData []AvailableBike) {
	bikeInterested := map[string]Item{}
	seedConvertedToMap := map[string]AvailableBike{}
	for _, entry := range seedData {
		seedConvertedToMap[entry.link] = entry
	}
	Interested := "Yes"
	svc := newSession()

	filt := expression.Name("Interested").Equal(expression.Value(Interested))
	proj := expression.NamesList(expression.Name("Interested"), expression.Name("Link"), expression.Name("Model"), expression.Name("TimeStamp"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		fmt.Println("there was an error!!", err)
	}
	params := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	result, err := svc.Scan(params)
	if err != nil {
		fmt.Println("There was an error:", err)
	}

	for _, entry := range result.Items {
		item := Item{}
		err = dynamodbattribute.UnmarshalMap(entry, &item)
		if err != nil {
			fmt.Println("there was an error: ", err)
		}

		if item.Interested == "Yes" {
			fmt.Println("you are on the right track", item.TimeStamp, item.Model, item.Link)
			bikeInterested[item.Link] = item
		}
	}
	for _, entry := range bikeInterested {
		_, found := seedConvertedToMap[entry.Link]
		if !found {
			fmt.Println("it's not here!!!!", entry)
			input := &dynamodb.DeleteItemInput{
				Key: map[string]*dynamodb.AttributeValue{
					"Link": {
						S: aws.String(entry.Link),
					},
				},
				TableName: aws.String(os.Getenv("TABLE_NAME")),
			}
			
			_, err := svc.DeleteItem(input)
			if err != nil {
				fmt.Println("Got error calling DeleteItem: ", err)
			}
			
			fmt.Println("Deleted " + entry.Link + " from table bike-availability")
			alertMe(entry.Link, entry.Model, "bike sold")
		}
	}
}

func updateDb(link string, model string, anIndex string) {
	loc, _ := time.LoadLocation("MST")
    now := time.Now().In(loc)
	
	svc := newSession()
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Link": {
				S: aws.String(link),
			},
			"TheIndex": {
				S: aws.String(anIndex),
			},
			"TimeStamp": {
				S: aws.String(now.String()),
			},
			"Model": {
				S: aws.String(model),
			},
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err := svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling UpdateItem: ", err)
	}
}

func alertMe(link string, model string, signature string) {
	message := ""
	TWILIO_ACCOUNT_SID := os.Getenv("TWILIO_ACCOUNT_SID")
	TWILIO_AUTH_TOKEN := os.Getenv("TWILIO_AUTH_TOKEN")
	if(TWILIO_ACCOUNT_SID == "" || TWILIO_AUTH_TOKEN == "") {
		fmt.Println("you need either an account sid or auth token")
		os.Exit(1)
	}
	client := twilio.NewRestClient()
	params := &openapi.CreateMessageParams{}
	if(signature == "new bike") {
		message = "New " + model + " available! " + link 
	} else if(signature == "bike sold") {
		message = model + " has been sold and deleted from the table " + link
	}
	params.SetTo(os.Getenv("PHONE_NUMBER"))
	params.SetFrom(os.Getenv("TWILIO_PHONE"))
	params.SetBody(message)
	
	_, err := client.ApiV2010.CreateMessage(params)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("SMS sent successfully!")
    }
}
	
func seedDB(seedData []AvailableBike) {
	loc, _ := time.LoadLocation("MST")
    now := time.Now().In(loc)

	svc := newSession()
	// seedData = seedData[1:]
	for index, entry := range seedData {
		anIndex := strconv.Itoa(index)
		fmt.Println("here is the index stringified: ", anIndex)
		if(entry.link == "" || entry.link == " ") {
			continue
		}
		input := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"Link": {
					S: aws.String(entry.link),
				},
				"TheIndex": {
					S: aws.String(anIndex),
				},
				"TimeStamp": {
					S: aws.String(now.String()),
				},
				"Model": {
					S: aws.String(entry.model),
				},
			},
			ReturnConsumedCapacity: aws.String("TOTAL"),
			TableName: aws.String(os.Getenv("TABLE_NAME")),
		}
		result, err := svc.PutItem(input)
		if err != nil {
			fmt.Println("I have an error", err)
			os.Exit(1)
		}
		fmt.Println(result)
	}
}
