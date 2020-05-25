package main

import (
	"longtrail-api/config"
	"longtrail-api/dbinterface"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/flick-web/dispatch"
)

var db dbinterface.DatabaseModifier
var conf *config.LongtrailConfig
var api *dispatch.API

func main() {
	conf = config.NewFromEnv()
	db = dbinterface.NewModifier(conf)

	api = &dispatch.API{}
	api.AddEndpoint("POST/practices/new", createPractice)
	api.AddEndpoint("GET/practice/{id}", getPractice)
	api.AddEndpoint("GET/practices", getPractices)
	api.AddEndpoint("POST/practice/{id}", setPractice)
	api.AddEndpoint("DELETE/practice/{id}", deletePractice)

	lambda.Start(api.LambdaProxy)
}
