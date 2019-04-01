package router

import (
	"github.com/ecletus/ecletus"
	"github.com/ecletus/cli"
	"github.com/ecletus/plug"
	"github.com/ecletus/router"
	"github.com/moisespsena-go/httpu"
	"github.com/moisespsena-go/task"
	"github.com/moisespsena-go/xroute"
	"github.com/moisespsena/go-error-wrap"
	"github.com/moisespsena/go-pluggable"
)

type Plugin struct {
	plug.EventDispatcher
	RouterKey    string
	ConfigDirKey string
	SingleSite   bool
}

func (p *Plugin) ProvideOptions() []string {
	return []string{p.RouterKey}
}

func (p *Plugin) RequireOptions() []string {
	return []string{p.ConfigDirKey}
}

func (p *Plugin) Init(options *plug.Options) error {
	var cfg httpu.Config
	configDir := options.GetInterface(p.ConfigDirKey).(*ecletus.ConfigDir)
	if err := configDir.Load(&cfg, "router.yaml"); err != nil {
		return errwrap.Wrap(err, "Load config file router.yaml")
	}

	if len(cfg.Listeners) == 0 {
		cfg.Listeners = append(cfg.Listeners, httpu.ListenerConfig{Addr: ":5000"})
	}

	r := &router.Router{
		Mux:    xroute.NewMux(router.PREFIX).LogRequests().InterseptErrors(),
		Config: &cfg,
	}

	options.Set(p.RouterKey, r)
	return nil
}

type ServerPlugin struct {
	pluggable.EventDispatcher
	RouterKey string
	PreServe  []func()
}

func (p *ServerPlugin) OnRegister(dis plug.PluginEventDispatcherInterface) {
	p.On(cli.E_REGISTER, func(e pluggable.PluginEventInterface) {
		r := e.Options().GetInterface(p.RouterKey).(*router.Router)
		if len(p.PreServe) > 0 {
			r.PreServe(func(r *router.Router, ta task.Appender) {
				for _, f := range p.PreServe {
					f()
				}
			})
		}
		rootCmd := e.(*cli.RegisterEvent).RootCmd
		agp := e.Options().GetInterface(ecletus.AGHAPE).(*ecletus.Ecletus)
		rootCmd.AddCommand(serveHttpCmd(r, agp, func(r *router.Router) error {
			return router.Trigger(dis, r)
		}))
	})
}
func (p *ServerPlugin) RequireOptions() []string {
	return []string{p.RouterKey}
}
