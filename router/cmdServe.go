package router

import (
	"github.com/ecletus/ecletus"
	"github.com/moisespsena-go/httpu"

	"github.com/ecletus/router"

	_ "github.com/ecletus/session"
	defaultlogger "github.com/moisespsena-go/default-logger"
	errwrap "github.com/moisespsena-go/error-wrap"
	"github.com/spf13/cobra"
)

func serveHttpCmd(r *router.Router, agp *ecletus.Ecletus, setupRoutes func(r *router.Router) error) *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve [ADDR]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Init server.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var log = defaultlogger.GetOrCreateLogger(router.PREFIX + ":server")
			log.Debug("Setup routes")
			if err = setupRoutes(r); err != nil {
				return errwrap.Wrap(err, "Setup Routes")
			}

			var listeners []httpu.ListenerConfig
			for _, arg := range args {
				listeners = append(listeners, httpu.ListenerConfig{Addr: httpu.Addr(arg)})
			}

			if len(listeners) > 0 {
				r.Config.Listeners = listeners
			}

			r.Cmd = cmd
			agp.AddTask(r.Server())
			return nil
		},
	}
	return serveCmd
}
