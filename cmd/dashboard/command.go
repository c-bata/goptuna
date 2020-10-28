package dashboard

import (
	"fmt"
	"net/http"
	"os"

	"github.com/c-bata/goptuna/dashboard"
	"github.com/c-bata/goptuna/internal/sqlalchemy"
	"github.com/c-bata/goptuna/rdb.v2"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// GetCommand returns the cobra's command for create-study sub-command.
func GetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "dashboard",
		Short:   "Launch web dashboard",
		Example: "  goptuna dashboard --storage sqlite:///example.db --host 127.0.0.1 --port 8000",
		Run: func(cmd *cobra.Command, args []string) {
			storageURL, err := cmd.Flags().GetString("storage")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			if storageURL == "" {
				cmd.PrintErrln("Storage URL is specified neither in config file nor --storage option.")
				os.Exit(1)
			}

			dialect, dsn, err := sqlalchemy.ParseDatabaseURL(storageURL, nil)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			var db *gorm.DB
			if dialect == "sqlite3" {
				db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
			} else if dialect == "mysql" {
				db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			} else if dialect == "postgres" {
				db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			} else {
				cmd.Println("unsupported dialect")
				os.Exit(1)
			}
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			storage := rdb.NewStorage(db)

			server, err := dashboard.NewServer(storage)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			hostName, err := cmd.Flags().GetString("host")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			addr := fmt.Sprintf("%s:%s", hostName, port)

			cmd.Printf("Started to serve at http://%s\n", addr)
			if err := http.ListenAndServe(addr, server); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		},
	}
	command.Flags().StringP(
		"storage", "", "", "DB URL specified in Engine Database URL format of SQLAlchemy (e.g. sqlite:///example.db). See https://docs.sqlalchemy.org/en/13/core/engines.html for more details.")
	command.Flags().StringP("host", "", "127.0.0.1", "hostname")
	command.Flags().StringP("port", "p", "8000", "port number")
	return command
}
