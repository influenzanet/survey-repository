package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/influenzanet/survey-repository/pkg/surveys"
	"github.com/spf13/cobra"
)

//var file string

func init() {
	rootCmd.AddCommand(LoadCmd)
	//LoadCmd.Flags().StringVar(&file, "file", "", "File to load")
}

var LoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load File",
	Long:  `All software has versions. This is Hugo's`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(args[0])
		if err != nil {
			log.Fatalf("Unable to load file %s", err)
		}
		survey, err := surveys.ExtractSurveyMetadata(data)
		if err != nil {
			log.Fatalf("Unable to load file %s", err)
		}
		b, err := json.Marshal(survey)
		if err != nil {
			log.Fatalf("Unable to serialize to json %s", err)
		}
		fmt.Println(string(b))
	},
}
