package main

import (
	"errors"
	"fmt"
	"longtrail-api/dbinterface"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/flick-web/dispatch"
)

var errorMissingID = errors.New("no id specified")
var errorUnauthorized = errors.New(http.StatusText(http.StatusUnauthorized))
var errorInternal = errors.New(http.StatusText(http.StatusInternalServerError))

func ctxToUserID(ctx events.APIGatewayProxyRequestContext) (id string, err error) {
	// AWS wraps claims.sub inside of an interface{}, inside of a map[string]interface{},
	// inside of another map[string]interface{}. This means there are four different
	// points at which this function might panic if we don't explicitly error check.
	// This is ridiculous, so just assume everything will go fine and recover in
	// case of panic.
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("ctxToUserID panic: %v\n", r)
			err = errorInternal
			return
		}
	}()
	claims := ctx.Authorizer["claims"]
	claimsMap := claims.(map[string]interface{})
	userID := claimsMap["sub"]
	id = userID.(string)
	return id, nil
}

func createPractice(practice *dbinterface.Practice, ctx *dispatch.Context) (id string, err error) {
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return "", err
	}
	practice.UserID = userID
	return db.CreatePractice(*practice)
}

func getPractice(ctx *dispatch.Context) (*dbinterface.Practice, error) {
	id, ok := ctx.PathVars["id"]
	if !ok {
		return nil, errorMissingID
	}
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return nil, err
	}
	return db.GetPractice(id, userID)
}

// getPractices expects query parameters "start" and "end" as RFC3339 timestamps.
func getPractices(ctx *dispatch.Context) ([]dbinterface.Practice, error) {
	var windowStart, windowEnd time.Time
	start, ok := ctx.LambdaRequest.QueryStringParameters["start"]
	if ok {
		var err error
		windowStart, err = time.Parse(time.RFC3339, start)
		windowStart = windowStart.UTC()
		if err != nil {
			return nil, err
		}
	}

	end, ok := ctx.LambdaRequest.QueryStringParameters["end"]
	if ok {
		var err error
		windowEnd, err = time.Parse(time.RFC3339, end)
		windowEnd = windowEnd.UTC()
		if err != nil {
			return nil, err
		}
	}

	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return nil, err
	}
	return db.GetPractices(userID, windowStart, windowEnd)
}

func setPractice(practice dbinterface.Practice, ctx *dispatch.Context) error {
	id, ok := ctx.PathVars["id"]
	if !ok || id == "undefined" || id == "null" {
		return errorMissingID
	}
	practice.ID = id
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return err
	}
	practice.UserID = userID
	return db.SetPractice(practice)
}

func deletePractice(ctx *dispatch.Context) error {
	id, ok := ctx.PathVars["id"]
	if !ok {
		return errorMissingID
	}
	userID, err := ctxToUserID(ctx.LambdaRequest.RequestContext)
	if err != nil {
		return err
	}
	return db.DeletePractice(id, userID)
}
