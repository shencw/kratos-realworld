package commands

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2/log"
	cliFlag "github.com/shencw/kratos-realworld/pkg/cli/flag"
	"github.com/shencw/kratos-realworld/pkg/cli/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDefaultCommand() (*cobra.Command, func()) {
	rootCmd := &cobra.Command{
		Use:   "realworld",
		Short: "realworld controls the platform",
		Long: templates.LongDesc(`
		realworld controls the platform, is the client side tool for realworld platform.

		test`),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return flushProfiling()
		},
	}

	flags := rootCmd.PersistentFlags()
	flags.SetNormalizeFunc(cliFlag.WordSepNormalizeFunc)

	addProfilingFlags(flags)
	addLogFlags(flags)

	_ = viper.BindPFlags(rootCmd.PersistentFlags())
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	rootCmd.SetGlobalNormalizationFunc(cliFlag.WordSepNormalizeFunc)

	logger, cleanup1 := newLogger()
	cmdService, cleanup2, err := LoadCommandService(context.Background(), logger)
	if err != nil {
		cleanup2()
		cleanup1()
		log.NewHelper(logger).Fatalf("LoadCommandService Error:%s", err)
	}

	for _, cmdGroup := range cmdService {
		rootCmd.AddCommand(cmdGroup.GetCommands())
	}

	return rootCmd, func() {
		cleanup2()
		cleanup1()
	}
}
