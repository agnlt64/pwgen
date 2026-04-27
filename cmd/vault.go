package cmd

import (
	"pwgen/internal/commands"
	"github.com/spf13/cobra"
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage your vaults",
}

var addCmd = &cobra.Command{
	Use: "new [name]",
	Short: "Create a new vault",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commands.NewVault(queries, cmd, args)
	},
}

var useCmd = &cobra.Command{
	Use: "use [name]",
	Short: "Use a vault",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commands.UseVault(queries, cmd, args)
	},
}

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List all vaults",
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		commands.ListVaults(queries, cmd, args)
	},
}

func init() {
	vaultCmd.AddCommand(addCmd)
	vaultCmd.AddCommand(useCmd)
	vaultCmd.AddCommand(listCmd)

	rootCmd.AddCommand(vaultCmd)
}
