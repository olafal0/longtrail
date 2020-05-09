package dbinterface_test

import (
	"longtrail-api/config"
	database "longtrail-api/dbinterface"
	"os"
	"testing"
	"time"
)

// Since these tests depend on the state of an AWS stack, they should not be run
// unless specified. Set TEST_DYNAMO=true to run these tests.
func TestMain(m *testing.M) {
	testEnv := os.Getenv("TEST_DYNAMO")
	if testEnv == "1" || testEnv == "true" {
		os.Exit(m.Run())
	} else {
		os.Exit(0)
	}
}

func TestPracticesInterface(t *testing.T) {
	conf := &config.LongtrailConfig{}
	conf.EventsTableName = "longtrail-testing-events"

	now := time.Now()

	practice := database.Practice{
		UserID:    "userId",
		StartTime: now,
		EndTime:   now.Add(60 * time.Minute),
	}

	db := database.NewModifier(conf)
	practiceID, err := db.CreatePractice(practice)
	if err != nil {
		t.Error(err)
	}

	resultPractice, err := db.GetPractice(practiceID, "userId")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if !resultPractice.StartTime.Equal(practice.StartTime) || !resultPractice.EndTime.Equal(practice.EndTime) {
		t.Errorf(
			"Practice times did not match: %v to %v (expected %v to %v)\n",
			resultPractice.StartTime, resultPractice.EndTime,
			practice.StartTime, practice.EndTime,
		)
	}

	practices, err := db.GetPractices(
		"userId",
		time.Now().Add(-360*time.Minute),
		time.Now().Add(360*time.Minute),
	)
	if err != nil {
		t.Error(err)
	}
	if len(practices) != 1 || practices[0].ID != practiceID {
		t.Errorf("GetPracticesByRoom mismatch: got %v, expected %v\n", practices, practice)
	}

	err = db.DeletePractice(practiceID, "userId")
	if err != nil {
		t.Error(err)
	}
}
