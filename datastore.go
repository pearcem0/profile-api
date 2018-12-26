package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// @TODO - parametize this using Region
var datastore = dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-west-2"))

func getChannelItem(channel string) (*contact, error) {
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String("Contact"),
		Key: map[string]*dynamodb.AttributeValue{
			"Channel": {
				S: aws.String(channel),
			},
		},
	}

	getItemResult, err := datastore.GetItem(getItemInput)
	if err != nil {
		return nil, err
	}
	if getItemResult.Item == nil {
		return nil, nil
	}

	//cont short for contact
	cont := new(contact)
	err = dynamodbattribute.UnmarshalMap(getItemResult.Item, cont)
	if err != nil {
		return nil, err
	}

	return cont, nil
}
