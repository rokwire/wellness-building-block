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
	coreAdapter "wellness/driven/core"
	"wellness/driven/notifications"
	storage "wellness/driven/storage"
	driver "wellness/driver/web"

	"github.com/golang-jwt/jwt"
	"github.com/rokwire/core-auth-library-go/v2/authservice"
	"github.com/rokwire/core-auth-library-go/v2/sigauth"
	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logs"
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

	serviceID := "wellness"

	loggerOpts := logs.LoggerOpts{SuppressRequests: logs.NewStandardHealthCheckHTTPRequestProperties(serviceID + "/version")}
	loggerOpts.SuppressRequests = append(loggerOpts.SuppressRequests, logs.NewStandardHealthCheckHTTPRequestProperties("wellness/version")...)
	logger := logs.NewLogger(serviceID, &loggerOpts)

	port := getEnvKey("PORT", true)

	internalAPIKey := getEnvKey("INTERNAL_API_KEY", true)

	// web adapter
	host := getEnvKey("WELLNESS_HOST", true)
	coreBBHost := getEnvKey("WELLNESS_CORE_BB_HOST", true)
	serviceURL := getEnvKey("WELLNESS_SERVICE_URL", true)

	//mongoDB adapter
	mongoDBAuth := getEnvKey("WELLNESS_MONGO_AUTH", true)
	mongoDBName := getEnvKey("WELLNESS_MONGO_DATABASE", true)
	mongoTimeout := getEnvKey("WELLNESS_MONGO_TIMEOUT", false)
	storageAdapter := storage.NewStorageAdapter(mongoDBAuth, mongoDBName, mongoTimeout)
	err := storageAdapter.Start()
	if err != nil {
		log.Fatal("Cannot start the mongoDB adapter - " + err.Error())
	}

	mtAppID := getEnvKey("WELLNESS_MULTI_TENANCY_APP_ID", true)
	mtOrgID := getEnvKey("WELLNESS_MULTI_TENANCY_ORG_ID", true)
	//serviceAccountID := getEnvKey("WELLNESS_SERVICE_ACCOUNT_ID", false)

	authService := auth.Service{
		ServiceID:   serviceID,
		ServiceHost: serviceURL,
		FirstParty:  true,
		AuthBaseURL: coreBBHost,
	}

	serviceRegLoader, err := auth.NewRemoteServiceRegLoader(&authService, []string{"rewards"})
	if err != nil {
		log.Fatalf("Error initializing remote service registration loader: %v", err)
	}

	serviceRegManager, err := auth.NewServiceRegManager(&authService, serviceRegLoader)
	if err != nil {
		log.Fatalf("Error initializing service registration manager: %v", err)
	}

	serviceAccountID := getEnvKey("WELLNESS_SERVICE_ACCOUNT_ID", false)
	privKeyRaw := getEnvKey("WELLNESS_PRIV_KEY", true)
	privKeyRaw = strings.ReplaceAll(privKeyRaw, "\\n", "\n")
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privKeyRaw))
	if err != nil {
		log.Fatalf("Error parsing priv key: %v", err)
	}
	signatureAuth, err := sigauth.NewSignatureAuth(privKey, serviceRegManager, false)
	if err != nil {
		log.Fatalf("Error initializing signature auth: %v", err)
	}

	serviceAccountLoader, err := auth.NewRemoteServiceAccountLoader(&authService, serviceAccountID, signatureAuth)
	if err != nil {
		log.Fatalf("Error initializing remote service account loader: %v", err)
	}

	serviceAccountManager, err := auth.NewServiceAccountManager(&authService, serviceAccountLoader)
	if err != nil {
		log.Fatalf("Error initializing service account manager: %v", err)
	}

	// Core adapter
	coreAdapter := coreAdapter.NewCoreAdapter(coreBBHost, serviceAccountManager)

	// Notification adapter
	notificationsBaseURL := getEnvKey("NOTIFICATIONS_BASE_URL", true)
	notificationsAdapter := notifications.NewNotificationsAdapter(notificationsBaseURL, notificationsBaseURL, serviceAccountManager, mtAppID, mtOrgID)
	if err != nil {
		log.Fatalf("Error initializing notification adapter: %v", err)
	}

	// application
	application := core.NewApplication(Version, Build, logger, storageAdapter, coreAdapter, notificationsAdapter, mtAppID, mtOrgID)
	application.Start()

	config := model.Config{
		CoreBBHost:     coreBBHost,
		ServiceURL:     serviceURL,
		InternalAPIKey: internalAPIKey,
	}

	webAdapter := driver.NewWebAdapter(host, port, application, config, serviceRegManager)

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
