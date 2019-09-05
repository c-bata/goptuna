package createstudy

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/goptuna"

	"github.com/c-bata/goptuna/internal/sqlalchemy"
	"github.com/c-bata/goptuna/rdb"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// GetCommand returns the cobra's command for create-study sub-command.
func GetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "create-study",
		Short:   "Create a study in your relational database storage.",
		Example: "  goptuna create-study --storage sqlite:///example.db",
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

			dialect, dbargs, err := sqlalchemy.ParseDatabaseURL(storageURL, nil)
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
				rdb.RunAutoMigrate(db)
			}

			studyName, err := cmd.Flags().GetString("study")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			direction := goptuna.StudyDirectionMinimize
			directionStr, err := cmd.Flags().GetString("direction")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			if strings.ToLower(directionStr) == "maximize" {
				direction = goptuna.StudyDirectionMaximize
			}

			storage := rdb.NewStorage(db)
			studyID, err := storage.CreateNewStudy(studyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			studyName, err = storage.GetStudyNameFromID(studyID)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			err = storage.SetStudyDirection(studyID, direction)
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
		"study", "", "",
		"A human-readable name of a study to distinguish it from others.")
	command.Flags().StringP(
		"direction", "", "minimize",
		"Set study direction.")
	// http://gorm.io/docs/migration.html
	command.Flags().BoolP(
		"without-migration", "", false,
		"Run create-study without running Auto-Migration. Auto-Migration will ONLY create tables, missing columns and missing indexes, and WON’T change existing column’s type or delete unused columns to protect your data.")
	return command
}
