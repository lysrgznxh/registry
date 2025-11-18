package commands

import (
	"agent-network-protocol/registry/core"
	"fmt"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "lookup version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("***************************************")
			fmt.Println("************ version info *************")
			fmt.Println("***************************************")
			fmt.Println("client version:", core.VERSION)
		},
	}
	return cmd
}
