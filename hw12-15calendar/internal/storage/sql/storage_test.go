package sqlstorage

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/cmd/calendar/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

var ctx = context.Background()

func initDB(st *Storage, t *testing.T) {
	conf := config.DBConf{"localhost", 5432, "event-db", "event-db", "event-db"}

	if err := st.Connect(ctx, &conf); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	if err := st.CreateSchemaDB(ctx); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) {
			if pgerr.Code != "42P07" {
				t.Errorf("CreateSchemaDB: %v", err)
			}
		}
	}

	if err := st.PrepareQuery(ctx); err != nil {
		t.Errorf("PrepareQuery: %v", err)
	}
}

func TestUserCRUD(t *testing.T) {
	var st Storage

	initDB(&st, t)

	storage.TestUserCRUD(ctx, t, &st)

	st.Close()
}

func TestEventCRUD(t *testing.T) {
	var st Storage

	initDB(&st, t)

	storage.TestEventCRUD(ctx, t, &st)

	st.Close()
}
