package dblayer

import (
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

func TestNew(t *testing.T) {
	client, err := mongo.NewClient("mongodb://localhost:80")
	if err != nil || client == nil {
		t.Errorf("failed to create a mongo client: %s", err.Error())
		return
	}

	db, err := New(client, "test", time.Duration(1)*time.Second)
	if err != nil || db == nil {
		t.Errorf("failed to create a db layer: %s", err.Error())
		return
	}

	if db.mongoClient == nil {
		t.Error("expected to have a valid mongo client in db.mongoClient, but got nil")
	}

	if db.dbName != "test" {
		t.Errorf("db.dbName was expected to be 'test', but got: '%s'", db.dbName)
	}

	if db.timeout != time.Duration(1)*time.Second {
		t.Errorf("db.timeout was expected to be '1s', but got: '%s'", db.timeout)
	}
}
