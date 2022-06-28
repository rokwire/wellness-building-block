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

package main

import (
	"log"
	"os"
	"strings"
	"wellness/core"
	"wellness/core/model"
	cacheadapter "wellness/driven/cache"
	storage "wellness/driven/storage"
	driver "wellness/driver/web"
)

var (
	// Version : version of this executable
	Version string
	// Build : build date of this executable
	Build string
)

func main() {
	if len(Version) == 0 {
		Version = "dev"
	}

	port := getEnvKey("PORT", true)

	//mongoDB adapter
	mongoDBAuth := getEnvKey("WELLNESS_MONGO_AUTH", true)
	mongoDBName := getEnvKey("WELLNESS_MONGO_DATABASE", true)
	mongoTimeout := getEnvKey("WELLNESS_MONGO_TIMEOUT", false)
	storageAdapter := storage.NewStorageAdapter(mongoDBAuth, mongoDBName, mongoTimeout)
	err := storageAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the mongoDB adapter - " + err.Error())
	}

	defaultCacheExpirationSeconds := getEnvKey("WELLNESS_DEFAULT_CACHE_EXPIRATION_SECONDS", false)
	cacheAdapter := cacheadapter.NewCacheAdapter(defaultCacheExpirationSeconds)

	mtAppID := getEnvKey("WELLNESS_MULTI_TENANCY_APP_ID", true)
	mtOrgID := getEnvKey("WELLNESS_MULTI_TENANCY_ORG_ID", true)

	// application
	application := core.NewApplication(Version, Build, storageAdapter, cacheAdapter, mtAppID, mtOrgID)
	application.Start()

	// web adapter
	host := getEnvKey("WELLNESS_HOST", true)
	coreBBHost := getEnvKey("WELLNESS_CORE_BB_HOST", true)
	serviceURL := getEnvKey("WELLNESS_SERVICE_URL", true)

	config := model.Config{
		CoreBBHost: coreBBHost,
		ServiceURL: serviceURL,
	}

	webAdapter := driver.NewWebAdapter(host, port, application, config)

	webAdapter.Start()
}

func getEnvKeyAsList(key string, required bool) []string {
	stringValue := getEnvKey(key, required)

	// it is comma separated format
	stringListValue := strings.Split(stringValue, ",")
	if len(stringListValue) == 0 && required {
		log.Fatalf("missing or empty env var: %s", key)
	}

	return stringListValue
}

func getEnvKey(key string, required bool) string {
	// get from the environment
	value, exist := os.LookupEnv(key)
	if !exist {
		if required {
			log.Fatal("No provided environment variable for " + key)
		} else {
			log.Printf("No provided environment variable for " + key)
		}
	}
	return value
}
