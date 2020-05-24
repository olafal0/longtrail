package dbinterface

import (
	"longtrail-api/config"
)

// A DatabaseModifier is an object that handles database operations.
type DatabaseModifier interface {
	PracticeModifier
}

// DatabaseHandler is an object that fulfills the DatabaseModifier interface.
type DatabaseHandler struct {
	Practices
}

// NewModifier creates a new object that satisfies the DatabaseModifier interface.
func NewModifier(cfg *config.LongtrailConfig) DatabaseModifier {
	db := DatabaseHandler{
		Practices: Practices{
			Dynamo: DynamoHandler{
				TableName:  cfg.EventsTableName,
				PrimaryKey: "id",
				SortKey:    "userId",
				Indexes: map[string]DynamoIndex{
					"userId-start-index": DynamoIndex{
						PrimaryKey: "userId",
						SortKey:    "start",
					},
				},
			},
			UserTimeIndexName: "userId-start-index",
		},
	}

	return &db
}
