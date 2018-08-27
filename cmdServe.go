package router

import (
	"net/http"

	_ "github.com/aghape/session"
	"github.com/moisespsena/go-default-logger"
	"github.com/moisespsena/go-error-wrap"
	"github.com/spf13/cobra"
)

func serveHttpCmd(router *Router, setupRoutes func(r *Router) error) *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve [ADDR]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Init server.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var log = defaultlogger.NewLogger(PREFIX + ":server")
			log.Debug("Setup routes")
			if err = setupRoutes(router); err != nil {
				return errwrap.Wrap(err, "Setup Routes")
			}
			addr := router.ServerAddr
			if len(args) > 0 && args[0] != "" {
				addr = args[0]
			}
			router.preServe()
			log.Info("Listening on", addr)
			return http.ListenAndServe(router.ServerAddr, router)
		},
	}
	return serveCmd
}
