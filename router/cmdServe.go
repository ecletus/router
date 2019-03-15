package router

import (
	"github.com/aghape/aghape"
	"github.com/moisespsena-go/httpu"

	"github.com/aghape/router"

	_ "github.com/aghape/session"
	"github.com/moisespsena/go-default-logger"
	"github.com/moisespsena/go-error-wrap"
	"github.com/spf13/cobra"
)

func serveHttpCmd(r *router.Router, agp *aghape.Aghape, setupRoutes func(r *router.Router) error) *cobra.Command {
	var serveCmd = &cobra.Command{
		Use:   "serve [ADDR]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Init server.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var log = defaultlogger.NewLogger(router.PREFIX + ":server")
			log.Debug("Setup routes")
			if err = setupRoutes(r); err != nil {
				return errwrap.Wrap(err, "Setup Routes")
			}

			var serverConfigs []httpu.ServerConfig
			for _, arg := range args {
				serverConfigs = append(serverConfigs, httpu.ServerConfig{Addr: httpu.Addr(arg)})
			}

			if len(serverConfigs) > 0 {
				r.Config.Servers = serverConfigs
			}

			r.Cmd = cmd
			agp.AddTask(r.CreateServer())
			return nil
		},
	}
	return serveCmd
}
