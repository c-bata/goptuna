package sqlalchemy

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

var (
	// ErrInvalidDatabaseURL means invalid as a SQLAlchemy's Engine URL format.
	ErrInvalidDatabaseURL = errors.New("invalid database url")
	// ErrUnsupportedDialect means the given dialect is unsupported.
	ErrUnsupportedDialect = errors.New("unsupported dialect")
)

// This regex pattern is based on following code:
// https://github.com/zzzeek/sqlalchemy/blob/c6554ac52/lib/sqlalchemy/engine/url.py#L234-L292
var engineURLPattern = regexp.MustCompile(
	`(?P<dialect>[\w]+)` +
		`(?:\+(?P<driver>[\w]+))?://` +
		`(?:` +
		`(?P<username>[^:/]*)` +
		`(?::(?P<password>.*))?` +
		`@)?` +
		`(?:` +
		`(?:` +
		`\[(?P<ipv6host>[^/]+)\] |` +
		`(?P<ipv4host>[^/:]+)` +
		`)?` +
		`(?::(?P<port>[^/]*))?` +
		`)?` +
		`(?:/(?P<database>[^?]*))?` +
		`(?:\?(?P<query>.*))?`)

// EngineOption to set the DSN option
type EngineOption struct {
	ParseTime bool
}

// ParseDatabaseURL parse SQLAlchemy's Engine URL format and returns Go's dialect and args.
func ParseDatabaseURL(url string, opt *EngineOption) (string, string, error) {
	// https://docs.sqlalchemy.org/en/13/core/engines.html
	// dialect+driver://username:password@host:port/database
	submatch := engineURLPattern.FindStringSubmatch(url)
	if submatch == nil {
		return "", "", ErrInvalidDatabaseURL
	}
	parsed := make(map[string]string, 8)
	for i, name := range engineURLPattern.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		parsed[name] = submatch[i]
	}

	var godialect string
	var dsn string
	var err error

	// TODO(c-bata): Support postgres
	switch parsed["dialect"] {
	case "sqlite":
		godialect = "sqlite3"
		dsn = parsed["database"]
	case "mysql":
		godialect = "mysql"
		dsn, err = buildMySQLArgs(parsed, opt)
	default:
		return "", "", ErrUnsupportedDialect
	}
	if err != nil {
		return "", "", err
	}

	return godialect, dsn, nil
}

func buildMySQLArgs(parsed map[string]string, opt *EngineOption) (string, error) {
	var dsn, unixpass, dbname string
	var query url.Values
	var err error

	dbname = parsed["database"]
	query, err = url.ParseQuery(parsed["query"])
	if err != nil {
		return "", err
	}

	protocol := "tcp"
	if parsed["driver"] == "pymysql" && query.Get("unix_socket") != "" {
		protocol = "unix"
		unixpass = query.Get("unix_socket")
	}

	dsn = parsed["username"]
	if parsed["password"] != "" {
		dsn += ":" + parsed["password"]
	}

	switch protocol {
	case "tcp":
		if parsed["port"] == "" {
			dsn += fmt.Sprintf("@tcp(%s)", parsed["ipv4host"])
		} else {
			dsn += fmt.Sprintf("@tcp(%s:%s)",
				parsed["ipv4host"], parsed["port"])
		}
	case "unix":
		dsn += fmt.Sprintf("@unix(%s)", unixpass)
	}
	dsn += "/" + dbname

	if opt != nil {
		if opt.ParseTime {
			dsn += "?parseTime=true"
		}
	}

	return dsn, nil
}
