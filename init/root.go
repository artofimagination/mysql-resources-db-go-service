package init

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/proemergotech/log/v3"
	"github.com/spf13/cobra"

	"github.com/artofimagination/mysql-resources-db-go-service/config"
	"github.com/artofimagination/mysql-resources-db-go-service/di"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: config.AppName,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := &config.Config{}
		initConfig(cfg)

		container, err := di.NewContainer(cfg)
		if err != nil {
			log.Panic(context.Background(), "Couldn't load container", "error", err)
		}
		defer container.Close()

		runner := newRunner()
		defer runner.stop()

		//
		//_, err := restcontrollers.NewRESTController()
		//if err != nil {
		//	panic(err)
		//}
		//
		//// Start HTTP server that accepts requests from the offer process to exchange SDP and Candidates
		//panic(http.ListenAndServe(":8080", nil))
		runner.start("rest server", container.RestServer.Start, container.RestServer.Stop)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-sigs:
		case err := <-runner.errors():
			log.Panic(context.Background(), err.Error(), "error", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Panic(context.Background(), err.Error(), "error", err)
	}
}
