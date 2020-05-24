package dbinterface

import (
	"errors"
	"fmt"
	"time"

	"github.com/guregu/dynamo"
)

// A PracticeModifier is an object that handles database operations relating to
// practice sessions.
type PracticeModifier interface {
	CreatePractice(Practice) (string, error)
	GetPractice(id, userID string) (*Practice, error)
	GetPractices(userID string, windowStart, windowEnd time.Time) ([]Practice, error)
	SetPractice(Practice) error
	DeletePractice(id, userID string) error
}

// A Practice is a logged practice session, associated with a user.
type Practice struct {
	ID             string    `dynamo:"id" json:"id"`
	UserID         string    `dynamo:"userId" json:"userId"`
	Start          time.Time `dynamo:"start" json:"start"`
	End            time.Time `dynamo:"end" json:"end"`
	AdditionalData string    `dynamo:"additionalData,omitempty" json:"additionalData,omitempty"`
}

// HashKey fulfills the dynamo.Keyed interface.
func (l *Practice) HashKey() interface{} {
	return l.ID
}

// RangeKey fulfills the dynamo.Keyed interface.
func (l *Practice) RangeKey() interface{} {
	return l.UserID
}

// Practices is a PracticeModifier for use with a Dynamo backend.
type Practices struct {
	Dynamo            DynamoHandler
	UserTimeIndexName string
}

// CreatePractice adds a practice with the specified timeframe and any
// additional data provided. This should always be used to create Practices,
// even though it is just a wrapper around SetPractice, as it is where IDs are
// created and set.
func (p *Practices) CreatePractice(practice Practice) (string, error) {
	practice.ID = p.Dynamo.newID()
	return practice.ID, p.SetPractice(practice)
}

// GetPractice returns information about a practice by ID.
func (p *Practices) GetPractice(id, userID string) (*Practice, error) {
	db, err := p.Dynamo.getClient()
	if err != nil {
		return nil, err
	}

	result := Practice{}
	table := db.Table(p.Dynamo.TableName)
	err = table.Get(p.Dynamo.PrimaryKey, id).Range(p.Dynamo.SortKey, dynamo.Equal, userID).One(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPractices searches for practices for a given student, filtering by timeframe.
func (p *Practices) GetPractices(userID string, start, end time.Time) ([]Practice, error) {
	db, err := p.Dynamo.getClient()
	if err != nil {
		return nil, err
	}

	index, ok := p.Dynamo.Indexes[p.UserTimeIndexName]
	if !ok {
		return nil, fmt.Errorf("Index %s not found in index configuration", p.UserTimeIndexName)
	}

	result := []Practice{}

	table := db.Table(p.Dynamo.TableName)
	query := table.Get(index.PrimaryKey, userID).Index(p.UserTimeIndexName)

	if !start.IsZero() {
		query = query.Range(index.SortKey, dynamo.GreaterOrEqual, start)
	}

	if !end.IsZero() {
		query = query.Range(index.SortKey, dynamo.Less, end)
	}

	err = query.All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// SetPractice writes the given practice to the database. Unlike CreatePractice,
// it requires the passed argument to already have an ID set.
func (p *Practices) SetPractice(practice Practice) error {
	db, err := p.Dynamo.getClient()
	if err != nil {
		return err
	}

	// Confirm that the practice object has an ID set
	if practice.ID == "" {
		return errors.New("ID not set in SetPractice request")
	}

	table := db.Table(p.Dynamo.TableName)
	err = table.Put(practice).Run()
	return err
}

// DeletePractice removes a practice specified by ID.
func (p *Practices) DeletePractice(id, userID string) error {
	db, err := p.Dynamo.getClient()
	if err != nil {
		return err
	}

	table := db.Table(p.Dynamo.TableName)
	err = table.Delete(p.Dynamo.PrimaryKey, id).Range(p.Dynamo.SortKey, userID).Run()
	return err
}
