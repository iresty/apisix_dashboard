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
package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/apache/apisix-dashboard/api/internal/conf"
	"github.com/apache/apisix-dashboard/api/internal/log"
	"github.com/apache/apisix-dashboard/api/internal/utils"
)

// server is the manager of Manager API, which is responsible for managing the life cycle of Manager API, including initialization, start, stop and so on
type server struct {
	server    *http.Server
	serverSSL *http.Server
	options   *Options
}

type Options struct{}

// NewServer Create a server manager
func NewServer(options *Options) (*server, error) {
	return &server{options: options}, nil
}

func (s *server) Start(errSig chan error) {
	// initialize server
	err := s.init()
	if err != nil {
		errSig <- err
		return
	}

	// print server info to stdout
	s.printInfo()

	// start HTTP server
	log.Infof("The Manager API is listening on %s", s.server.Addr)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Errorf("listen and serv fail: %s", err)
			errSig <- err
		}
	}()

	// start HTTPs server
	if conf.SSLCert != "" && conf.SSLKey != "" {
		go func() {
			err := s.serverSSL.ListenAndServeTLS(conf.SSLCert, conf.SSLKey)
			if err != nil && err != http.ErrServerClosed {
				log.Errorf("listen and serve for HTTPS failed: %s", err)
				errSig <- err
			}
		}()
	}
}

func (s *server) Stop() {
	utils.CloseAll()

	s.shutdownServer(s.server)
	s.shutdownServer(s.serverSSL)
}

func (s *server) init() error {
	log.Info("Initialize Manager API store")
	err := s.setupStore()
	if err != nil {
		return err
	}

	log.Info("Initialize Manager API server")
	s.setupAPI()

	return nil
}

func (s *server) shutdownServer(server *http.Server) {
	if server != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Errorf("Shutting down server error: %s", err)
		}
	}
}

func (s *server) printInfo() {
	fmt.Fprint(os.Stdout, "The manager-api is running successfully!\n\n")
	utils.PrintVersion()
	fmt.Fprintf(os.Stdout, "%-8s: %s\n", "Config File", viper.ConfigFileUsed())
	fmt.Fprintf(os.Stdout, "%-8s: %s:%d\n", "Listen", conf.ServerHost, conf.ServerPort)
	if conf.SSLCert != "" && conf.SSLKey != "" {
		fmt.Fprintf(os.Stdout, "%-8s: %s:%d\n", "HTTPS Listen", conf.SSLHost, conf.SSLPort)
	}
	fmt.Fprintf(os.Stdout, "%-8s: %s\n", "Loglevel", conf.ErrorLogLevel)
	fmt.Fprintf(os.Stdout, "%-8s: %s\n", "ErrorLogFile", conf.ErrorLogPath)
	fmt.Fprintf(os.Stdout, "%-8s: %s\n\n", "AccessLogFile", conf.AccessLogPath)
}
