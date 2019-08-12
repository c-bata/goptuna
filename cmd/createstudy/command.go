package createstudy

import (
	"fmt"
	"os"

	"github.com/c-bata/goptuna/internal/sqlalchemy"
	"github.com/c-bata/goptuna/rdbstorage"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// GetCommand returns the cobra's command for create-study sub-command.
func GetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "create-study",
		Short:   "Create a study in your relational database storage.",
		Example: "  goptuna create-study --storage sqlite:///example.db --study-name test-study",
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

			dialect, dbargs, err := sqlalchemy.ParseDatabaseURL(storageURL)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			db, err := gorm.Open(dialect, dbargs...)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			defer db.Close()

			withoutMigrate, err := cmd.Flags().GetBool("without-migration")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			if !withoutMigrate {
				rdbstorage.RunAutoMigrate(db)
			}

			studyName, err := cmd.Flags().GetString("study-name")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			storage := rdbstorage.NewStorage(db)
			_, err = storage.CreateNewStudyID(studyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			fmt.Println(studyName)
		},
	}
	command.Flags().StringP(
		"storage", "", "", "DB URL specified in Engine Database URL format of SQLAlchemy (e.g. sqlite:///example.db). See https://docs.sqlalchemy.org/en/13/core/engines.html for more details.")
	command.Flags().StringP(
		"study-name", "", "",
		"A human-readable name of a study to distinguish it from others.")
	// http://gorm.io/docs/migration.html
	command.Flags().BoolP(
		"without-migration", "", false,
		"Run create-study without running Auto-Migration. Auto-Migration will ONLY create tables, missing columns and missing indexes, and WON’T change existing column’s type or delete unused columns to protect your data.")
	return command
}
