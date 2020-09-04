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

type RouteCallback func(r *Router)
type RouteCallbacks []RouteCallback

func (this *RouteCallbacks) Append(cb ...RouteCallback) *RouteCallbacks {
	*this = append(*this, cb...)
	return this
}

func (this RouteCallbacks) AppendCopy(cb ...RouteCallback) RouteCallbacks {
	return append(this, cb...)
}

type Router struct {
	Mux                *xroute.Mux
	PrioritaryHandlers httpu.FallbackHandlers
	PrefixHandlers     httpu.PrefixHandlers
	Handler            http.Handler
	Config             *httpu.Config
	preServeCallbacks  []func(srv *httpu.Server)
	Tasks              task.Tasks
	Cmd                *cobra.Command
	server             *httpu.Server
	RouteCallbacks     RouteCallbacks
}

func (r *Router) GetMux() *xroute.Mux {
	return r.Mux
}

func (r *Router) PreServe(cb ...func(srv *httpu.Server)) {
	r.preServeCallbacks = append(r.preServeCallbacks, cb...)
}

func (r *Router) CallPreServe(srv *httpu.Server) {
	for _, cb := range r.preServeCallbacks {
		cb(srv)
	}
}

func (r *Router) Server() *httpu.Server {
	if r.server == nil {
		srv := httpu.NewServer(r.Config, r)
		srv.PreSetup(func(srv *httpu.Server) error {
			r.CallPreServe(srv)
			return nil
		})
		r.server = srv
	}
	return r.server
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	httpu.Fallback(r.PrioritaryHandlers, r.PrefixHandlers, r.Handler).ServeHTTP(w, req)
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

func OnRoute(dis plug.EventDispatcherInterface, callbacks ...func(e *RouterEvent)) {
	for _, cb := range callbacks {
		dis.On(E_ROUTE, RouterEventCallback(cb))
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
