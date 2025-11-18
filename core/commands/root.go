package commands

import (
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "registry",
		Short: "registry",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	//启动
	rootCmd.AddCommand(StartCmd())
	//版本
	rootCmd.AddCommand(VersionCmd())
	return rootCmd
}
