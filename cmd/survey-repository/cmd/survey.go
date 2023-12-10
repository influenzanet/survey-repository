package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(surveyCmd)
	surveyCmd.AddCommand(surveyShowCmd)
}

var surveyCmd = &cobra.Command{
	Use:   "survey",
	Short: "Survey management",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var surveyShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show Survey",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mg := loadManager()
		i := args[0]
		id, err := strconv.Atoi(i)
		if err != nil {
			return fmt.Errorf("Unable to parse id : %s", err)
		}
		data, err := mg.GetSurveyData(uint(id), true)
		if err != nil {
			return fmt.Errorf("Unable to show survey : %s", err)
		}
		fmt.Println(string(data))
		return nil
	},
}
