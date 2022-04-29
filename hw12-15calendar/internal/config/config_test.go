package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	wantConf := Config{InMem: false,
		DB:     DBConf{"172.17.0.2", 5432, "user-event-db", "123456", "event-db"},
		HTTP:   HTTPConf{"127.0.0.1", 8080},
		GRPC:   GRPCConf{"localhost", 55555},
		Logger: LoggerConf{"path/to/file", "info"},
	}

	conf, err := Parse("./testdata/config.yml")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if *conf != wantConf {
		t.Errorf("configuration mismatch:\n\thave: %v\n\twant: %v", *conf, wantConf)
	}
}
