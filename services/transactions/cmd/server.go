// Package cmd provides the command line interface for the application.
package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/model"
	"github.com/fardinabir/digital-wallet-demo/services/transactions/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(_ *cobra.Command, _ []string) {
		if err := runServe(cfg); err != nil {
			log.Fatal(err)
		}
	},
}

func runServe(cfg model.Config) error {
	var servers []server.Server

	apiOpts := server.TxnAPIServerOpts{
		ListenPort: cfg.APIServer.Port,
		Config:     cfg,
	}
	apiServer, err := server.NewAPI(apiOpts)
	if err != nil {
		return err
	}
	servers = append(servers, apiServer)

	if cfg.SwaggerServer.Enable {
		SwaggerOpts := server.SwaggerServerOpts{
			ListenPort: cfg.SwaggerServer.Port,
		}
		swagServer := server.NewSwagger(SwaggerOpts)
		servers = append(servers, swagServer)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for _, s := range servers {
		srvr := s
		go func() {
			if err := srvr.Run(); err != nil && err != http.ErrServerClosed {
				log.Fatal("shutting down ", srvr.Name(), " err: ", err)
			}
		}()
	}

	log.Info("server started")
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	log.Info("server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, s := range servers {
		if err := s.Shutdown(ctx); err != nil {
			return err
		}
	}
	log.Info("server shutdown gracefully")
	return nil
}
