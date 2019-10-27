package deletestudy

import (
	"os"

	"github.com/c-bata/goptuna/internal/sqlalchemy"
	"github.com/c-bata/goptuna/rdb"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// GetCommand returns the cobra's command for create-study sub-command.
func GetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     "delete-study",
		Short:   "Delete a study in your relational database storage.",
		Example: "  goptuna delete-study --storage sqlite:///example.db --study study",
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

			studyName, err := cmd.Flags().GetString("study")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			storage := rdb.NewStorage(db)
			studyID, err := storage.GetStudyIDFromName(studyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			err = storage.DeleteStudy(studyID)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		},
	}
	command.Flags().StringP(
		"storage", "", "", "DB URL specified in Engine Database URL format of SQLAlchemy (e.g. sqlite:///example.db). See https://docs.sqlalchemy.org/en/13/core/engines.html for more details.")
	command.Flags().StringP(
		"study", "", "",
		"A human-readable name of a study to distinguish it from others.")
	return command
}
