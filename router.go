package router

import (
	"net/http"

	"github.com/moisespsena/go-pluggable"
	"github.com/moisespsena/go-route"
	"github.com/aghape/cli"
	"github.com/aghape/plug"
)

var E_ROUTE = PREFIX + ":route"

type RouterEvent struct {
	plug.PluginEventInterface
	Router *Router
}

type Plugin struct {
	plug.EventDispatcher
	RouterKey  string
	ServerAddr string
	SingleSite bool
}

func (p *Plugin) ProvideOptions() []string {
	return []string{p.RouterKey}
}

func (p *Plugin) Init(options *plug.Options) error {
	if p.ServerAddr == "" {
		p.ServerAddr = ":5000"
	}
	router := &Router{
		Mux:        route.NewMux(PREFIX).LogRequests().InterseptErrors(),
		ServerAddr: p.ServerAddr,
	}

	options.Set(p.RouterKey, router)
	return nil
}

type Router struct {
	Mux        *route.Mux
	Handler    http.Handler
	ServerAddr string
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req = route.SetOriginalUrlIfNotSetted(req)
	r.Handler.ServeHTTP(w, req)
}

type ServerPlugin struct {
	pluggable.EventDispatcher
	RouterKey string
}

func (p *ServerPlugin) OnRegister(dis plug.PluginEventDispatcherInterface) {
	p.On(cli.E_REGISTER, func(e pluggable.PluginEventInterface) {
		router := e.Options().GetInterface(p.RouterKey).(*Router)
		rootCmd := e.(*cli.RegisterEvent).RootCmd
		rootCmd.AddCommand(serveHttpCmd(router, func(r *Router) error {
			return Trigger(dis, r)
		}))
	})
}

func RouterEventCallback(cb func(e *RouterEvent)) func(e plug.EventInterface) {
	return func(e plug.EventInterface) {
		cb(e.(*RouterEvent))
	}
}

func RouterEventCallbackE(cb func(e *RouterEvent) error) func(e plug.EventInterface) error {
	return func(e plug.EventInterface) error {
		return cb(e.(*RouterEvent))
	}
}

func OnRoute(dis plug.EventDispatcherInterface, callbacks ...func(e *RouterEvent)) (err error) {
	for _, cb := range callbacks {
		err = dis.OnE(E_ROUTE, RouterEventCallback(cb))
		if err != nil {
			return
		}
	}
	return
}

func OnRouteE(dis plug.EventDispatcherInterface, callbacks ...func(e *RouterEvent) error) (err error) {
	for _, cb := range callbacks {
		err = dis.OnE(E_ROUTE, RouterEventCallbackE(cb))
		if err != nil {
			return
		}
	}
	return
}

func (p *ServerPlugin) RequireOptions() []string {
	return []string{p.RouterKey}
}

func Trigger(dis plug.PluginEventDispatcherInterface, r *Router) error {
	return dis.TriggerPlugins(&RouterEvent{plug.NewPluginEvent(E_ROUTE), r})
}
