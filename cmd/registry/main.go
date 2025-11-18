package main

import (
	"agent-network-protocol/registry/core/commands"
)

func main() {
	rootCmd := commands.RootCmd()
	rootCmd.Execute()
}
