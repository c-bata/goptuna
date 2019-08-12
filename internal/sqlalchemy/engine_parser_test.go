package sqlalchemy

import (
	"reflect"
	"testing"
)

func TestParseDatabaseURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantDialect string
		wantArgs    []interface{}
	}{
		{
			name:        "sqlite3 simple",
			url:         "sqlite:///example.db",
			wantDialect: "sqlite3",
			wantArgs: []interface{}{
				"example.db",
			},
		},
		{
			name:        "sqlite3 simple2",
			url:         "sqlite:///db.sqlite3",
			wantDialect: "sqlite3",
			wantArgs: []interface{}{
				"db.sqlite3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDialect, gotArgs, err := ParseDatabaseURL(tt.url)
			if err != nil {
				t.Errorf("ParseDatabaseURL() err = %s, want nil", err)
			}
			if gotDialect != tt.wantDialect {
				t.Errorf("ParseDatabaseURL() gotDialect = %v, want %v", gotDialect, tt.wantDialect)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("ParseDatabaseURL() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
