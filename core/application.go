// Copyright 2022 Board of Trustees of the University of Illinois.
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

package core

import (
	"log"
	"sync"

	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logs"
)

// Application represents the core application code based on hexagonal architecture
type Application struct {
	version string
	build   string

	logger *logs.Logger

	cacheLock *sync.Mutex

	Services Services //expose to the drivers adapters

	storage       Storage
	core          Core
	notifications Notifications

	//TODO - remove this when applied to all environemnts
	multiTenancyAppID string
	multiTenancyOrgID string

	deleteDataLogic deleteDataLogic
}

// Start starts the core part of the application
func (app *Application) Start() {
	err := app.deleteDataLogic.start()
	if err != nil {
		log.Fatalf("error on starting the delete data logic - %s", err)
	}

	err = app.MigrateMessageIDs()
	if err != nil {
		log.Printf("error on migrate message ids - %s", err)
	}
}

// NewApplication creates new Application
func NewApplication(version string, build string,
	logger *logs.Logger, storage Storage,
	core Core, notifications Notifications, mtAppID string, mtOrgID string) *Application {
	cacheLock := &sync.Mutex{}

	deleteDataLogic := deleteDataLogic{logger: logger, coreAdapter: core, storage: storage}

	application := Application{version: version, build: build, logger: logger, cacheLock: cacheLock, storage: storage,
		core: core, notifications: notifications, multiTenancyAppID: mtAppID, multiTenancyOrgID: mtOrgID,
		deleteDataLogic: deleteDataLogic}

	// add the drivers ports/interfaces
	application.Services = &servicesImpl{app: &application}

	return &application
}
