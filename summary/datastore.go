package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var datastore = dynamodb.New(session.New(), aws.NewConfig().WithRegion(region))
var summaryTableName = "profile-api-" + stage + "-summary"

func getSectionItem(section string) (*summary, error) {

	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(summaryTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Section": {
				S: aws.String(section),
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

	//summ short for summary
	summ := new(summary)
	err = dynamodbattribute.UnmarshalMap(getItemResult.Item, summ)
	if err != nil {
		return nil, err
	}

	return summ, nil
}

func addSectionItem(summ *summary) error {
	// summ = short for summary
	newSection := &dynamodb.PutItemInput{
		TableName: aws.String(summaryTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"Section": {
				S: aws.String(summ.Section),
			},
			"Content": {
				S: aws.String(summ.Content),
			},
		},
	}

	_, err := datastore.PutItem(newSection)
	return err
}
