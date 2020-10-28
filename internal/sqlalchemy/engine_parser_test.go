package sqlalchemy_test

import (
	"reflect"
	"testing"

	"github.com/c-bata/goptuna/internal/sqlalchemy"
)

func TestParseDatabaseURL(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		option *sqlalchemy.EngineOption

		// Go DSN (Data Source Name)
		// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		wantDialect string
		wantDsn     string
	}{
		// SQLite3
		// sqlite://<nohostname>/<path>
		{
			name:        "sqlite3 simple",
			url:         "sqlite:///example.db",
			wantDialect: "sqlite3",
			wantDsn:     "example.db",
		},
		{
			name:        "sqlite3 simple2",
			url:         "sqlite:///db.sqlite3",
			wantDialect: "sqlite3",
			wantDsn:     "db.sqlite3",
		},
		// MySQL
		{
			name:        "mysql",
			url:         "mysql://scott:tiger@localhost/foo",
			wantDialect: "mysql",
			wantDsn:     "scott:tiger@tcp(localhost)/foo",
		},
		{
			name:        "mysql (with driver)",
			url:         "mysql+pymysql://user:pass@localhost:6000/bar",
			wantDialect: "mysql",
			wantDsn:     "user:pass@tcp(localhost:6000)/bar",
		},
		{
			name:        "mysql (with unix domain socket)",
			url:         "mysql+pymysql://username:password@localhost/foo?unix_socket=/var/lib/mysql/mysql.sock",
			wantDialect: "mysql",
			wantDsn:     "username:password@unix(/var/lib/mysql/mysql.sock)/foo",
		},
		{
			name: "mysql (with parsetime option)",
			url:  "mysql+mysqldb://user:pass@localhost:6000/bar",
			option: &sqlalchemy.EngineOption{
				ParseTime: true,
			},
			wantDialect: "mysql",
			wantDsn:     "user:pass@tcp(localhost:6000)/bar?parseTime=true",
		},
		// Postgres
		{
			name:        "postgres",
			url:         "postgresql://scott:tiger@localhost/mydatabase",
			wantDialect: "postgres",
			wantDsn:     "user=scott password=tiger dbname=mydatabase",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDialect, gotArgs, err := sqlalchemy.ParseDatabaseURL(tt.url, tt.option)
			if err != nil {
				t.Errorf("ParseDatabaseURL() err = %s, want nil", err)
			}
			if gotDialect != tt.wantDialect {
				t.Errorf("ParseDatabaseURL() gotDialect = %v, want %v", gotDialect, tt.wantDialect)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantDsn) {
				t.Errorf("ParseDatabaseURL() gotArgs = %v, want %v", gotArgs, tt.wantDsn)
			}
		})
	}
}
