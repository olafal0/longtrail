package main

import (
	"errors"
	"longtrail-api/dbinterface"
	"time"

	"github.com/flick-web/dispatch"
)

var errorMissingID = errors.New("no id specified")

func createPractice(practice *dbinterface.Practice, ctx *dispatch.Context) (id string, err error) {
	userID := ctx.LambdaRequest.RequestContext.Identity.User
	practice.UserID = userID
	return db.CreatePractice(*practice)
}

func getPractice(ctx *dispatch.Context) (*dbinterface.Practice, error) {
	id, ok := ctx.PathVars["id"]
	if !ok {
		return nil, errorMissingID
	}
	userID := ctx.LambdaRequest.RequestContext.Identity.User
	return db.GetPractice(id, userID)
}

// getPractices expects query parameters "start" and "end" as RFC3339 timestamps.
func getPractices(ctx *dispatch.Context) ([]dbinterface.Practice, error) {
	var windowStart, windowEnd time.Time
	start, ok := ctx.LambdaRequest.QueryStringParameters["start"]
	if ok {
		var err error
		windowStart, err = time.Parse(time.RFC3339, start)
		if err != nil {
			return nil, err
		}
	}

	end, ok := ctx.LambdaRequest.QueryStringParameters["end"]
	if ok {
		var err error
		windowEnd, err = time.Parse(time.RFC3339, end)
		if err != nil {
			return nil, err
		}
	}

	userID := ctx.LambdaRequest.RequestContext.Identity.User
	return db.GetPractices(userID, windowStart, windowEnd)
}

func setPractice(practice dbinterface.Practice, ctx *dispatch.Context) error {
	userID := ctx.LambdaRequest.RequestContext.Identity.User
	practice.UserID = userID
	return db.SetPractice(practice)
}

func deletePractice(ctx *dispatch.Context) error {
	id, ok := ctx.PathVars["id"]
	if !ok {
		return errorMissingID
	}
	userID := ctx.LambdaRequest.RequestContext.Identity.User
	return db.DeletePractice(id, userID)
}
