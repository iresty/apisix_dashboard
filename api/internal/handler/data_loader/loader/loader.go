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
package loader

import "github.com/apache/apisix-dashboard/api/internal/core/entity"

// DataSets are intermediate structures used to handle
// import and export data with APISIX entities.
// On import, raw data will be parsed as DataSets
// On export, DataSets will be encoded to raw data
type DataSets struct {
	Routes        []entity.Route
	Upstreams     []entity.Upstream
	Services      []entity.Service
	Consumers     []entity.Consumer
	SSLs          []entity.SSL
	StreamRoutes  []entity.StreamRoute
	GlobalPlugins []entity.GlobalPlugins
	PluginConfigs []entity.PluginConfig
	Protos        []entity.Proto
}

// Loader provide data loader abstraction
type Loader interface {
	// Import accepts data and converts it into entity data sets
	Import(input interface{}) (*DataSets, error)

	// Export accepts entity data sets and converts it into a specific format
	Export(data DataSets) (interface{}, error)
}
