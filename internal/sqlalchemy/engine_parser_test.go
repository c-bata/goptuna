package sqlalchemy_test

import (
	"reflect"
	"testing"

	"github.com/c-bata/goptuna/internal/sqlalchemy"
)

func TestParseDatabaseURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantDialect string
		wantArgs    []interface{}
	}{
		// Python Engine URL format.
		// sqlite://<nohostname>/<path>
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
		// Go DSN (Data Source Name)
		// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		{
			name:        "mysql",
			url:         "mysql://scott:tiger@localhost/foo",
			wantDialect: "mysql",
			wantArgs: []interface{}{
				"scott:tiger@tcp(localhost)/foo",
			},
		},
		{
			name:        "mysql (with driver)",
			url:         "mysql+pymysql://user:pass@localhost:6000/bar",
			wantDialect: "mysql",
			wantArgs: []interface{}{
				"user:pass@tcp(localhost:6000)/bar",
			},
		},
		{
			name:        "mysql (with unix domain socket)",
			url:         "mysql+pymysql://username:password@localhost/foo?unix_socket=/var/lib/mysql/mysql.sock",
			wantDialect: "mysql",
			wantArgs: []interface{}{
				"username:password@unix(/var/lib/mysql/mysql.sock)/foo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDialect, gotArgs, err := sqlalchemy.ParseDatabaseURL(tt.url)
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
