package router

import (
	"net/http"

	"github.com/moisespsena-go/task"
	"github.com/spf13/cobra"

	"github.com/ecletus/plug"
	"github.com/moisespsena-go/httpu"
	"github.com/moisespsena-go/xroute"
)

var E_ROUTE = PREFIX + ":route"

type RouterEvent struct {
	plug.PluginEventInterface
	Router *Router
}

type Router struct {
	rootMux           *xroute.Mux
	Mux               *xroute.Mux
	Handler           http.Handler
	Config            *httpu.Config
	preServeCallbacks []func(r *Router, ta task.Appender)
	Tasks             task.Tasks
	Cmd               *cobra.Command
}

func (r *Router) GetRootMux() *xroute.Mux {
	if r.rootMux != nil {
		return r.rootMux
	}
	return r.Mux
}

func (r *Router) RootMux(f func(mux *xroute.Mux)) {
	if r.rootMux == nil {
		r.rootMux = xroute.NewMux(PREFIX + ":root")
		r.rootMux.Mount("/", r.Mux)
	}
	f(r.rootMux)
}

func (r *Router) PreServe(cb ...func(r *Router, ta task.Appender)) {
	r.preServeCallbacks = append(r.preServeCallbacks, cb...)
}

func (r *Router) CallPreServe(ta task.Appender) {
	for _, cb := range r.preServeCallbacks {
		cb(r, ta)
	}
}

func (r *Router) CreateServer() *httpu.Server {
	srv := httpu.NewServer(r.Config, r)
	srv.AddSetup(func(ta task.Appender) error {
		r.CallPreServe(ta)
		return nil
	})
	return srv
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req = xroute.SetOriginalURLIfNotSetted(req)
	r.Handler.ServeHTTP(w, req)
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

func Trigger(dis plug.PluginEventDispatcherInterface, r *Router) error {
	return dis.TriggerPlugins(&RouterEvent{plug.NewPluginEvent(E_ROUTE), r})
}
