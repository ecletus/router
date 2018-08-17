// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"net/http"

	"github.com/moisespsena/go-default-logger"
	"github.com/moisespsena/go-error-wrap"
	_ "github.com/aghape/session"
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
			log.Info("Listening on", addr)
			return http.ListenAndServe(router.ServerAddr, router)
		},
	}
	return serveCmd
}
