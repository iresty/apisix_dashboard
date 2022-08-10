/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/apache/apisix-dashboard/api/internal/conf"
	"github.com/apache/apisix-dashboard/api/internal/core/server"
	"github.com/apache/apisix-dashboard/api/internal/log"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apisix-dashboard",
		Short: "Apache APISIX Dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf.InitConf()
			log.InitLogger()

			s, err := server.NewServer(&server.Options{})
			if err != nil {
				return err
			}

			// start Manager API server
			errSig := make(chan error, 5)
			s.Start(errSig)

			// Signal received to the process externally.
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case sig := <-quit:
				log.Infof("The Manager API server receive %s and start shutting down", sig.String())
				s.Stop()
				log.Infof("See you next time!")
			case err := <-errSig:
				log.Errorf("The Manager API server start failed: %s", err.Error())
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.ConfigFile, "config", "c", "", "config file")
	cmd.PersistentFlags().StringVarP(&conf.WorkDir, "work-dir", "p", ".", "current work directory")

	cmd.AddCommand(
		newVersionCommand(),
	)

	return cmd
}
