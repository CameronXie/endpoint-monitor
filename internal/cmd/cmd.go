package cmd

import (
	"os"

	"github.com/CameronXie/endpoint-monitor/internal/status"
	"github.com/spf13/cobra"
)

type Logger interface {
	// Infoln logs any given args as information.
	Infoln(args ...interface{})

	// Errorln logs any given args as error.
	Errorln(args ...interface{})
}

func Execute(svc status.MonitorService, l Logger) error {
	rootCmd := &cobra.Command{
		Use: "EndpointMonitor",
	}

	rootCmd.AddCommand(setupStatusCMD(svc, l))
	return rootCmd.Execute()
}

func exitOnError(err error, l Logger) {
	if err != nil {
		l.Errorln(err)
		os.Exit(1)
	}
}
