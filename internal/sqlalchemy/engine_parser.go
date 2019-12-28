package sqlalchemy

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	// ErrInvalidDatabaseURL means invalid as a SQLAlchemy's Engine URL format.
	ErrInvalidDatabaseURL = errors.New("invalid database url")
	// ErrUnsupportedDialect means the given dialect is unsupported.
	ErrUnsupportedDialect = errors.New("unsupported dialect")
)

// https://github.com/zzzeek/sqlalchemy/blob/c6554ac52bfb7ce9ecd30ec777ce90adfe7861d2/lib/sqlalchemy/engine/url.py#L234-L292
var rfc1738pattern = regexp.MustCompile(
	`(?P<name>[\w\+]+)://` +
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
		`(?:/(?P<database>.*))?`)

// EngineOption to set the DSN option
type EngineOption struct {
	ParseTime bool
}

// ParseDatabaseURL parse SQLAlchemy's Engine URL format and returns Go's dialect and args.
func ParseDatabaseURL(url string, opt *EngineOption) (string, []interface{}, error) {
	// https://docs.sqlalchemy.org/en/13/core/engines.html
	// dialect+driver://username:password@host:port/database
	submatch := rfc1738pattern.FindStringSubmatch(url)
	if submatch == nil {
		return "", nil, ErrInvalidDatabaseURL
	}
	parsed := make(map[string]string, 8)
	for i, name := range rfc1738pattern.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		parsed[name] = submatch[i]
	}

	var pydialect, pydriver string
	if strings.Contains(parsed["name"], "+") {
		x := strings.SplitN(parsed["name"], "+", 2)
		pydialect = x[0]
		pydriver = x[1]
	} else {
		pydialect = parsed["name"]
	}

	var godialect string
	var dbargs []interface{}
	var err error
	switch pydialect {
	case "sqlite":
		godialect = "sqlite3"
		dbargs = []interface{}{
			parsed["database"],
		}
	case "mysql":
		godialect = "mysql"
		dbargs, err = buildMySQLArgs(pydriver, parsed, opt)
	default:
		return "", nil, ErrUnsupportedDialect
	}
	if err != nil {
		return "", nil, err
	}

	return godialect, dbargs, nil
}

func buildMySQLArgs(pydriver string, parsed map[string]string, opt *EngineOption) ([]interface{}, error) {
	var godsn, unixpass, database string
	var query url.Values
	var err error

	x := strings.SplitN(parsed["database"], "?", 2)
	database = x[0]
	if len(x) == 2 {
		query, err = url.ParseQuery(x[1])
		if err != nil {
			return nil, err
		}
	}

	protocol := "tcp"
	if pydriver == "pymysql" && query.Get("unix_socket") != "" {
		protocol = "unix"
		unixpass = query.Get("unix_socket")
	}

	godsn = parsed["username"]
	if parsed["password"] != "" {
		godsn += ":" + parsed["password"]
	}

	switch protocol {
	case "tcp":
		if parsed["port"] == "" {
			godsn += fmt.Sprintf("@tcp(%s)", parsed["ipv4host"])
		} else {
			godsn += fmt.Sprintf("@tcp(%s:%s)",
				parsed["ipv4host"], parsed["port"])
		}
	case "unix":
		godsn += fmt.Sprintf("@unix(%s)", unixpass)
	}
	godsn += "/" + database

	if opt != nil {
		if opt.ParseTime {
			godsn += "?parseTime=true"
		}
	}

	return []interface{}{
		godsn,
	}, nil
}
