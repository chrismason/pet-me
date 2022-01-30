package cmd

import "github.com/spf13/cobra"

var (
	verbose bool
)

var RootCmd = &cobra.Command{
	Use:   "pet-me",
	Short: "Things you can do with pets",
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Log all the things")

	RootCmd.AddCommand(NewFactsCommand())
	RootCmd.AddCommand(NewPicsCommand())
}
