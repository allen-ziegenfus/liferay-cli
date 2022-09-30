/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package ext

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"liferay.com/lcectl/cetypes"
	"liferay.com/lcectl/prereq"
)

var whitespace = regexp.MustCompile(`\s`)

// refreshCmd represents the refresh command
var createCmd = &cobra.Command{
	Use:   "create [OPTIONS] [FLAGS]",
	Short: "Creates new Client Extensions using a wizard-like interface",
	Run: func(cmd *cobra.Command, args []string) {
		prereq.Prereq(Verbose)

		validate := func(input string) error {
			if len(input) <= 0 {
				return errors.New("Project Name must not be empty")
			}
			if whitespace.MatchString(input) {
				return errors.New("Project Name must not contain spaces")
			}
			return nil
		}

		projectNamePrompt := promptui.Prompt{
			Label:    "Project Name",
			Validate: validate,
		}

		projectName, err := projectNamePrompt.Run()

		if err != nil {
			fmt.Printf("Project Name entry failed %v\n", err)
			return
		}

		dat, err := cetypes.ClientExtensionTypeKeys(Verbose)

		if err != nil {
			fmt.Printf("Error getting Client Extension Types %v\n", err)
			return
		}

		cetPrompt := promptui.Select{
			Label: "Type",
			Items: dat,
		}

		_, cetType, err := cetPrompt.Run()

		if err != nil {
			fmt.Printf("Error getting Client Extension Types %v\n", err)
			return
		}

		fmt.Printf("Project Name: %q\n", projectName)
		fmt.Printf("Type: %q\n", cetType)
	},
}

func init() {
	createCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "enable verbose output")
	extCmd.AddCommand(createCmd)
}
