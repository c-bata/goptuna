package sqlalchemy

import (
	"errors"
	"strings"
)

var (
	ErrInvalidDatabaseURL = errors.New("invalid database url")
	ErrUnsupportedDialect = errors.New("unsupported dialect")
)

func ParseDatabaseURL(url string) (dialect string, args []interface{}, err error) {
	// https://docs.sqlalchemy.org/en/13/core/engines.html
	// dialect+driver://username:password@host:port/database
	x := strings.Split(url, "://")
	if len(x) != 2 {
		return "", nil, ErrInvalidDatabaseURL
	}
	dialect, err = getGoDialect(x[0])
	if err != nil {
		return "", nil, err
	}
	y := strings.Split(x[1], "/")
	if len(y) != 2 {
		return "", nil, ErrInvalidDatabaseURL
	}
	database := y[1]

	return dialect, []interface{}{database}, nil
}

func getGoDialect(alchemyDialect string) (string, error) {
	if strings.Contains(alchemyDialect, "+") {
		alchemyDialect = strings.Split(alchemyDialect, "+")[0]
	}

	switch alchemyDialect {
	case "sqlite":
		return "sqlite3", nil
	default:
		return "", ErrUnsupportedDialect
	}
}
