package cmd

import (
	"pwgen/internal/commands"
	"github.com/spf13/cobra"
)

var passCmd = &cobra.Command{
	Use:   "pass",
	Short: "Manage your passwords",
}

var newCmd = &cobra.Command{
	Use: "new [website] [label]",
	Short: "Create a password for [website]. Use [label] to retrieve the password once saved.",
	Args: cobra.MinimumNArgs(2),
	Run: func (cmd *cobra.Command, args []string)  {
		commands.NewPass(queries, cmd, args)
	},
}

var getCmd = &cobra.Command{
	Use: "get [label]",
	Short: "Get the password associated with [label]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commands.GetPass(queries, cmd, args)
	},
}

func init() {
	newCmd.Flags().Int("length", 25, "Password length")

	passCmd.AddCommand(newCmd)
	passCmd.AddCommand(getCmd)

	rootCmd.AddCommand(passCmd)
}