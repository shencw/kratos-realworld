package main

import (
	"fmt"
	"github.com/shencw/kratos-realworld/internal/commands"
	"os"
)

func main() {
	cmd, cleanup := commands.NewDefaultCommand()
	defer cleanup()
	if err := cmd.Execute(); err != nil {
		fmt.Println("Execute Error: ", err)
		os.Exit(1)
	}
}
