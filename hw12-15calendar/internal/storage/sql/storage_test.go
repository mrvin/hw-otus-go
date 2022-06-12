package sqlstorage

import (
	"context"
	"testing"

	"github.com/mrvin/hw-otus-go/hw12-15calendar/cmd/calendar/config"
	"github.com/mrvin/hw-otus-go/hw12-15calendar/internal/storage"
)

var ctx = context.Background()

func TestUserCRUD(t *testing.T) {
	conf := config.DBConf{"localhost", 5432, "event-db", "event-db", "event-db"}
	st, err := New(ctx, &conf)
	if err != nil {
		t.Fatalf("db: %v", err)
	}

	storage.TestUserCRUD(ctx, t, st)

	if err := st.DropSchemaDB(ctx); err != nil {
		t.Fatalf("DropSchemaDB: %v", err)
	}
	st.Close()
}

func TestEventCRUD(t *testing.T) {
	conf := config.DBConf{"localhost", 5432, "event-db", "event-db", "event-db"}
	st, err := New(ctx, &conf)
	if err != nil {
		t.Fatalf("db: %v", err)
	}

	storage.TestEventCRUD(ctx, t, st)

	if err := st.DropSchemaDB(ctx); err != nil {
		t.Fatalf("DropSchemaDB: %v", err)
	}
	st.Close()
}
