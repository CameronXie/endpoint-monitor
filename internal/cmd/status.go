package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/CameronXie/endpoint-monitor/internal/status"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	statusCommand  = "status"
	statusShort    = "check the endpoint status and store the result"
	stopMonitorMsg = "gracefully stop status monitoring"
)

func setupStatusCMD(
	svc status.MonitorService,
	l Logger,
) *cobra.Command {
	var cfgFile string
	s := make(chan os.Signal, 1)
	cmd := &cobra.Command{
		Use:   statusCommand,
		Short: statusShort,
		Run:   statusRun(&cfgFile, s, svc, l),
	}

	cmd.Flags().StringVarP(&cfgFile, "config", "f", "", "config file (required)")
	_ = cmd.MarkFlagRequired("config")

	return cmd
}

func statusRun(
	f *string,
	s chan os.Signal,
	svc status.MonitorService,
	l Logger,
) func(cmd *cobra.Command, args []string) {
	return func(_ *cobra.Command, _ []string) {
		signal.Notify(s, os.Interrupt, syscall.SIGTERM)

		viper.SetConfigFile(*f)
		exitOnError(viper.ReadInConfig(), l)

		eps := make([]status.Endpoint, 0)
		exitOnError(viper.UnmarshalKey("endpoints", &eps), l)

		ctx, cancel := context.WithCancel(context.Background())
		svc.Monitor(ctx, eps)

		<-s
		l.Infoln(stopMonitorMsg)
		cancel()
		svc.Stop()
	}
}
