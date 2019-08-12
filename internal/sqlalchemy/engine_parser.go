package sqlalchemy

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	// ErrInvalidDatabaseURL means invalid as a SQLAlchemy's Engine URL format.
	ErrInvalidDatabaseURL = errors.New("invalid database url")
	// ErrUnsupportedDialect means the given dialect is unsupported.
	ErrUnsupportedDialect = errors.New("unsupported dialect")
)

// ParseDatabaseURL parse SQLAlchemy's Engine URL format and returns Go's dialect and args.
func ParseDatabaseURL(url string) (string, []interface{}, error) {
	// https://docs.sqlalchemy.org/en/13/core/engines.html
	// dialect+driver://username:password@host:port/database
	x := strings.SplitN(url, "://", 2)
	if len(x) != 2 {
		return "", nil, ErrInvalidDatabaseURL
	}

	var pydialect, pydriver string
	if strings.Contains(x[0], "+") {
		y := strings.SplitN(x[0], "+", 2)
		pydialect = y[0]
		pydriver = y[1]
	} else {
		pydialect = x[0]
	}

	var godialect string
	var dbargs []interface{}
	var err error
	switch pydialect {
	case "sqlite":
		godialect = "sqlite3"
		dbargs, err = parseSQLiteArgs(x[1])
	case "mysql":
		godialect = "mysql"
		dbargs, err = parseMySQLArgs(pydriver, x[1])
	default:
		return "", nil, ErrUnsupportedDialect
	}
	if err != nil {
		return "", nil, err
	}

	return godialect, dbargs, nil
}

func parseSQLiteArgs(pyargs string) ([]interface{}, error) {
	database := strings.TrimLeft(pyargs, "/")
	return []interface{}{
		database,
	}, nil
}

func parseMySQLArgs(pydriver string, pyargs string) ([]interface{}, error) {
	var godsn string
	var username, password string
	var database string
	var query url.Values
	var unixpass string
	protocol := "tcp"
	hostname := "localhost:3306"

	if strings.Contains(pyargs, "@") {
		x := strings.SplitN(pyargs, "@", 2)
		userpass := x[0]
		hostInfoAndOption := x[1]

		if strings.Contains(userpass, ":") {
			y := strings.SplitN(userpass, ":", 2)
			username = y[0]
			password = y[1]
		} else {
			username = userpass
		}

		var hostinfo string
		if strings.Contains(hostInfoAndOption, "?") {
			y := strings.SplitN(hostInfoAndOption, "?", 2)

			var err error
			query, err = url.ParseQuery(y[1])
			if err != nil {
				return nil, err
			}
			hostinfo = y[0]
		} else {
			hostinfo = hostInfoAndOption
		}

		z := strings.SplitN(hostinfo, "/", 2)
		if len(z) != 2 {
			return nil, errors.New("cannot extract database name")
		}
		hostname = z[0]
		database = z[1]
	}

	if query != nil && pydriver == "pymysql" {
		protocol = "unix"
		unixpass = query.Get("unix_socket")
	}

	switch protocol {
	case "tcp":
		godsn = fmt.Sprintf("tcp(%s)/%s", hostname, database)
	case "unix":
		godsn = fmt.Sprintf("unix(%s)/%s", unixpass, database)
	}

	if username != "" {
		if password != "" {
			godsn = fmt.Sprintf("%s:%s@", username, password) + godsn
		} else {
			godsn = fmt.Sprintf("%s@", username) + godsn
		}
	}

	return []interface{}{
		godsn,
	}, nil
}
