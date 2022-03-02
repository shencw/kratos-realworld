package pkg

import "github.com/spf13/cobra"

type CommandService interface {
	GetCommands() *cobra.Command
}
