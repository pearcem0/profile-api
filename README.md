# profile-api

An experimental and contrived preview of **Golang**. A very simple event driven API to fetch information about the Author, Michael Pearce.

A 'Services' approach - Separate lambda functions for each category of information contact, summary, interests etc. Including all of the basic actions - get, put, delete etc.

## Example calls

`curl "https://api.michaelpearce.info/contact?channel=Email"`
`curl "https://api.michaelpearce.info/contact?channel=Facebook"`
`curl "https://api.michaelpearce.info/contact?channel=Instagram"`
`curl "https://api.michaelpearce.info/contact?channel=LinkedIn"`
`curl "https://api.michaelpearce.info/contact?channel=Website"`

## Infrastructure

Build on services provided by Amazon Web Services (AWS) cloud.
* api gateway, for a scalable backend
* lambda, to run 'serverless' code functions in response to events from API Gateway
* dynamodb, for data store
* cloudfront, for CDN
* route 53 for domain name registration and hosting, and CNAME record set
* cloudwatch for logging and cloudwatch/sns for alarms

## Prerequisites

### AWS Account & Infrastructure

The project is build on AWS Infrastructure and assumes the infrastructure mentioned in the above 'Infrastructure' section is already deployed and configured to run the code. Along with the correct permissions!

### GO packages

`go get github.com/aws/aws-lambda-go/lambda`
`go get github.com/aws/aws-sdk-go`
`go get github.com/aws/aws-lambda-go/events`

## Build


### Env Vars & Go build

Set Go OS and Architecture variables so that it is built to run on the Linux servers that AWS Lambda uses. Build it generally wherever you want, but the temp folder keeps things tidy in the long run. 

`env GOOS=linux GOARCH=amd64 go build -o /tmp/contact contact`

### Package it up for Lambda

Package it up as a zip file for Lambda.
`zip -j /tmp/main.zip /tmp/contact`


### Deploy new code version on Lambda

Deploy to Lambda, however you prefer. For a quick proof of concept you can use the CLI. 

`aws lambda update-function-code --function-name contact --zip-file fileb:///tmp/main.zip`

#### Workaround for update-function-code command timeout

If the aws command times out because the package is too big (and your network bandwith is low), it is easier to copy to S3 first, then deploy from there. 

`aws s3 cp /tmp/main.zip s3://<YOUR-S3-BUCKET-NAME>/`

`aws lambda update-function-code --function-name contact --s3-bucket <YOUR-S3-BUCKET-NAME> --s3-key main.zip`
