/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package core

import (
	"sync"
	cacheadapter "wellness/driven/cache"
)

//Application represents the core application code based on hexagonal architecture
type Application struct {
	version string
	build   string

	cacheLock *sync.Mutex

	Services Services //expose to the drivers adapters

	cache   *cacheadapter.CacheAdapter
	storage Storage

	//TODO - remove this when applied to all environemnts
	multiTenancyAppID string
	multiTenancyOrgID string
}

// Start starts the core part of the application
func (app *Application) Start() {

}

// NewApplication creates new Application
func NewApplication(version string, build string, storage Storage,
	cacheadapter *cacheadapter.CacheAdapter, mtAppID string, mtOrgID string) *Application {
	cacheLock := &sync.Mutex{}
	application := Application{version: version, build: build, cacheLock: cacheLock, storage: storage,
		cache: cacheadapter, multiTenancyAppID: mtAppID, multiTenancyOrgID: mtOrgID}

	// add the drivers ports/interfaces
	application.Services = &servicesImpl{app: &application}

	return &application
}