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

package web

import (
	"fmt"
	"github.com/rokwire/core-auth-library-go/authservice"
	"log"
	"net/http"
	"wellness/core"
	"wellness/core/model"
	"wellness/driver/web/rest"
	"wellness/utils"

	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	httpSwagger "github.com/swaggo/http-swagger"
)

//Adapter entity
type Adapter struct {
	host string
	port string
	auth *Auth

	apisHandler         rest.ApisHandler
	adminApisHandler    rest.AdminApisHandler
	internalApisHandler rest.InternalApisHandler

	app *core.Application
}

// @title Rokwire Wellness Building Block API
// @description Rokwire Content Building Block API Documentation.
// @version 1.0.1
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost
// @BasePath /wellness
// @schemes https

// @securityDefinitions.apikey UserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

// @securityDefinitions.apikey AdminUserAuth
// @in header (add Bearer prefix to the Authorization value)
// @name Authorization

// @securityDefinitions.apikey InternalAPIAuth
// @in header (add INTERNAL-API-KEY header with an appropriate value)
// @name Authorization

// @securityDefinitions.apikey AdminGroupAuth
// @in header
// @name GROUP

//Start starts the module
func (we Adapter) Start() {

	router := mux.NewRouter().StrictSlash(true)

	// handle apis
	subRouter := router.PathPrefix("/wellness").Subrouter()
	subRouter.PathPrefix("/doc/ui").Handler(we.serveDocUI())
	subRouter.HandleFunc("/doc", we.serveDoc)
	subRouter.HandleFunc("/version", we.wrapFunc(we.apisHandler.Version)).Methods("GET")

	subRouter = subRouter.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/int/process_reminders", we.internalAuthWrapFunc(we.internalApisHandler.ProcessReminders)).Methods("POST")

	// handle user todo categories apis
	subRouter.HandleFunc("/user/todo_categories", we.coreAuthWrapFunc(we.apisHandler.GetUserTodoCategories, we.auth.coreAuth.standardAuth)).Methods("GET")
	subRouter.HandleFunc("/user/todo_categories", we.coreAuthWrapFunc(we.apisHandler.CreateUserTodoCategory, we.auth.coreAuth.standardAuth)).Methods("POST")
	subRouter.HandleFunc("/user/todo_categories/{id}", we.coreAuthWrapFunc(we.apisHandler.GetUserTodoCategory, we.auth.coreAuth.standardAuth)).Methods("GET")
	subRouter.HandleFunc("/user/todo_categories/{id}", we.coreAuthWrapFunc(we.apisHandler.UpdateUserTodoCategory, we.auth.coreAuth.standardAuth)).Methods("PUT")
	subRouter.HandleFunc("/user/todo_categories/{id}", we.coreAuthWrapFunc(we.apisHandler.DeleteUserTodoCategory, we.auth.coreAuth.standardAuth)).Methods("DELETE")

	// handle user todo entries apis
	subRouter.HandleFunc("/user/todo_entries", we.coreAuthWrapFunc(we.apisHandler.GetUserTodoEntries, we.auth.coreAuth.standardAuth)).Methods("GET")
	subRouter.HandleFunc("/user/todo_entries", we.coreAuthWrapFunc(we.apisHandler.CreateUserTodoEntry, we.auth.coreAuth.standardAuth)).Methods("POST")
	subRouter.HandleFunc("/user/todo_entries/clear_completed_entries", we.coreAuthWrapFunc(we.apisHandler.DeleteCompletedUserTodoEntry, we.auth.coreAuth.standardAuth)).Methods("DELETE")
	subRouter.HandleFunc("/user/todo_entries/{id}", we.coreAuthWrapFunc(we.apisHandler.GetUserTodoEntry, we.auth.coreAuth.standardAuth)).Methods("GET")
	subRouter.HandleFunc("/user/todo_entries/{id}", we.coreAuthWrapFunc(we.apisHandler.UpdateUserTodoEntry, we.auth.coreAuth.standardAuth)).Methods("PUT")
	subRouter.HandleFunc("/user/todo_entries/{id}", we.coreAuthWrapFunc(we.apisHandler.DeleteUserTodoEntry, we.auth.coreAuth.standardAuth)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":"+we.port, router))
}

func (we Adapter) serveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("access-control-allow-origin", "*")
	http.ServeFile(w, r, "./docs/swagger.yaml")
}

func (we Adapter) serveDocUI() http.Handler {
	url := fmt.Sprintf("%s/wellness/doc", we.host)
	return httpSwagger.Handler(httpSwagger.URL(url))
}

func (we Adapter) wrapFunc(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		handler(w, req)
	}
}

type coreAuthFunc = func(*tokenauth.Claims, http.ResponseWriter, *http.Request)

func (we Adapter) coreAuthWrapFunc(handler coreAuthFunc, authorization Authorization) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		responseStatus, claims, err := authorization.check(req)
		if err != nil {
			log.Printf("error authorization check - %s", err)
			http.Error(w, http.StatusText(responseStatus), responseStatus)
			return
		}
		handler(claims, w, req)
	}
}

type internalAuthFunc = func(http.ResponseWriter, *http.Request)

func (we Adapter) internalAuthWrapFunc(handler internalAuthFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.LogRequest(req)

		status, err := we.auth.internalAuth.check(req)
		if err != nil {
			log.Printf("error authorization check - %s", err)
			http.Error(w, http.StatusText(status), status)
			return
		}

		handler(w, req)
	}
}

// NewWebAdapter creates new WebAdapter instance
func NewWebAdapter(host string, port string, app *core.Application, config model.Config, authService *authservice.AuthService) Adapter {
	auth := NewAuth(app, config, authService)

	apisHandler := rest.NewApisHandler(app)
	adminApisHandler := rest.NewAdminApisHandler(app)
	internalApisHandler := rest.NewInternalApisHandler(app)
	return Adapter{host: host, port: port, auth: auth, apisHandler: apisHandler, adminApisHandler: adminApisHandler,
		internalApisHandler: internalApisHandler, app: app}
}

// AppListener implements core.ApplicationListener interface
type AppListener struct {
	adapter *Adapter
}
