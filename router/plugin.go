package router

import (
	"github.com/ecletus/cli"
	"github.com/ecletus/ecletus"
	"github.com/ecletus/plug"
	"github.com/ecletus/router"
	errwrap "github.com/moisespsena-go/error-wrap"
	"github.com/moisespsena-go/httpu"
	"github.com/moisespsena-go/pluggable"
	"github.com/moisespsena-go/xroute"
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

func (p *Plugin) ProvidesOptions(options *plug.Options) error {
	var cfg httpu.Config
	configDir := options.GetInterface(p.ConfigDirKey).(*ecletus.ConfigDir)
	if err := configDir.Load(&cfg, "router.yaml"); err != nil {
		return errwrap.Wrap(err, "Load config file router.yaml")
	}

	if len(cfg.Listeners) == 0 {
		cfg.Listeners = append(cfg.Listeners, httpu.ListenerConfig{Addr: ":5000"})
	}

	r := &router.Router{
		Mux:    xroute.NewMux(router.PREFIX).InterseptErrors(),
		Config: &cfg,
	}

	options.Set(p.RouterKey, r)
	return nil
}

type ServerPlugin struct {
	pluggable.EventDispatcher
	RouterKey string
	PreServe  []func(srv *httpu.Server)
}

func (p *ServerPlugin) OnRegister(options *plug.Options) {
	p.On(cli.E_REGISTER, func(e pluggable.PluginEventInterface) {
		r := e.Options().GetInterface(p.RouterKey).(*router.Router)
		if len(p.PreServe) > 0 {
			r.PreServe(func(srv *httpu.Server) {
				for _, f := range p.PreServe {
					f(srv)
				}
			})
		}
		rootCmd := e.(*cli.RegisterEvent).RootCmd
		ect := e.Options().GetInterface(ecletus.ECLETUS).(*ecletus.Ecletus)
		rootCmd.AddCommand(serveHttpCmd(r, ect, func(r *router.Router) error {
			for _, f := range r.RouteCallbacks {
				f(r)
			}
			dis := pluggable.Dis(options)
			return router.Trigger(dis, r)
		}))
	})
}

func (p *ServerPlugin) RequireOptions() []string {
	return []string{p.RouterKey}
}
