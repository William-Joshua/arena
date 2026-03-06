package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"cc.io/arena/internal/server"
	log "cc.io/arena/pkg/logging"
)

func newApiserverCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "apiserver",
		Short: "Run the API server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			app := fx.New(
				log.Module,
				server.Module(server.Config{Port: port}),
				// Invoke the router so Fx registers the lifecycle hooks.
				fx.Invoke(func(*gin.Engine) {}),
			)
			app.Run()
			return app.Err()
		},
	}

	cmd.Flags().IntVar(&port, "port", 8080, "port to listen on")
	return cmd
}
