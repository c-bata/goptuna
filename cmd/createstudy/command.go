package createstudy

import (
	"fmt"
	"os"

	"github.com/c-bata/goptuna/internal/sqlalchemy"
	"github.com/c-bata/goptuna/rdbstorage"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func GetCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "create-study",
		Short: "create a study",
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
			storage, err := rdbstorage.NewStorage(db)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			db.AutoMigrate(&rdbstorage.StudyModel{})
			db.AutoMigrate(&rdbstorage.StudyUserAttributeModel{})
			db.AutoMigrate(&rdbstorage.StudySystemAttributeModel{})
			db.AutoMigrate(&rdbstorage.TrialModel{})
			db.AutoMigrate(&rdbstorage.TrialUserAttributeModel{})
			db.AutoMigrate(&rdbstorage.TrialSystemAttributeModel{})
			db.AutoMigrate(&rdbstorage.TrialParamModel{})
			db.AutoMigrate(&rdbstorage.TrialValueModel{})

			studyName, err := cmd.Flags().GetString("study-name")
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			_, err = storage.CreateNewStudyID(studyName)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			fmt.Println(studyName)
		},
	}
	command.Flags().StringP(
		"storage", "", "", "DB URL. (e.g. sqlite:///example.db)")
	command.Flags().StringP(
		"study-name", "", "",
		"A human-readable name of a study to distinguish it from others.")
	return command
}
