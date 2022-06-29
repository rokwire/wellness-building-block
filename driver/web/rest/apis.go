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

package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/tokenauth"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"wellness/core"
	"wellness/core/model"
)

const maxUploadSize = 15 * 1024 * 1024 // 15 mb

//ApisHandler handles the rest APIs implementation
type ApisHandler struct {
	app *core.Application
}

//Version gives the service version
// @Description Gives the service version.
// @Tags Client
// @ID Version
// @Produce plain
// @Success 200
// @Router /version [get]
func (h ApisHandler) Version(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.app.Services.GetVersion()))
}

// GetUserTodoCategories Retrieves all user todo categories
// @Description Retrieves all user todo categories
// @Tags Client
// @ID GetUserTodoCategories
// @Accept json
// @Success 200
// @Security UserAuth
// @Router  /api/user/todo_categories [get]
func (h ApisHandler) GetUserTodoCategories(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	resData, err := h.app.Services.GetTodoCategories(claims.AppID, claims.OrgID, claims.Subject)
	if err != nil {
		log.Printf("Error on getting user todo categories - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.TodoCategory{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal all user todo categories: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserTodoCategory Retrieves a user todo category by id
// @Description Retrieves a user todo category by id
// @Tags Client
// @ID GetUserTodoCategory
// @Accept json
// @Produce json
// @Success 200
// @Security UserAuth
// @Router  /api/user/todo_categories/{id} [get]
func (h ApisHandler) GetUserTodoCategory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetTodoCategory(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on getting user todo category by id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal user todo category: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateUserTodoCategory Updates a user todo category with the specified id
// @Description Updates a user todo category with the specified id
// @Tags Client
// @ID UpdateUserTodoCategory
// @Accept json
// @Produce json
// @Success 200
// @Security UserAuth
// @Router /api/user/todo_categories/{id} [put]
func (h ApisHandler) UpdateUserTodoCategory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user todo category - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.TodoCategory
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create user todo category request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id != item.ID {
		log.Printf("Inconsistent attempt to update query param  and json id are not equal")
		http.Error(w, "Inconsistent attempt to update query param  and json id are not equal", http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateTodoCategory(claims.AppID, claims.OrgID, claims.Subject, &item)
	if err != nil {
		log.Printf("Error on updating user todo category with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal the updated user todo category: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateUserTodoCategory Creates a user todo category
// @Description Creates a user todo category
// @Tags Client
// @ID CreateUserTodoCategory
// @Accept json
// @Success 200
// @Security UserAuth
// @Router /api/user/todo_categories [post]
func (h ApisHandler) CreateUserTodoCategory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user todo category - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.TodoCategory
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create user todo category request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateTodoCategory(claims.AppID, claims.OrgID, claims.Subject, &item)
	if err != nil {
		log.Printf("Error on creating user todo category: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on marshal the new user todo category: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteUserTodoCategory Deletes a user todo category with the specified id
// @Description Deletes a user todo category with the specified id
// @Tags Client
// @ID DeleteUserTodoCategory
// @Success 200
// @Security UserAuth
// @Router /api/user/todo_categories/{id} [delete]
func (h ApisHandler) DeleteUserTodoCategory(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteTodoCategory(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on deleting user todo category with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetUserTodoEntries Retrieves all user todo entries
// @Description Retrieves all user todo entries
// @Tags Client
// @ID GetUserTodoEntries
// @Accept json
// @Success 200
// @Security UserAuth
// @Router  /api/user/todo_entries [get]
func (h ApisHandler) GetUserTodoEntries(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	resData, err := h.app.Services.GetTodoEntries(claims.AppID, claims.OrgID, claims.Subject)
	if err != nil {
		log.Printf("Error on getting user todo entries - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.TodoEntry{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal all user todo entries: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserTodoEntry Retrieves a user todo entry by id
// @Description Retrieves a user todo entry by id
// @Tags Client
// @ID GetUserTodoEntry
// @Accept json
// @Produce json
// @Success 200
// @Security UserAuth
// @Router  /api/user/todo_entries/{id} [get]
func (h ApisHandler) GetUserTodoEntry(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetTodoEntry(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on getting user todo entry by id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		log.Printf("Error on getting user todo entry by id - %s (Not found)", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal user todo entry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateUserTodoEntry Updates a user todo entry with the specified id
// @Description Updates a user todo entry with the specified id
// @Tags Client
// @ID UpdateUserTodoEntry
// @Accept json
// @Produce json
// @Success 200
// @Security UserAuth
// @Router /api/user/todo_entries/{id} [put]
func (h ApisHandler) UpdateUserTodoEntry(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user todo entry - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.TodoEntry
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create user todo entry request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id != item.ID {
		log.Printf("Inconsistent attempt to update todo entry - query param  and json id are not equal")
		http.Error(w, "Inconsistent attempt to update todo entry - query param  and json id are not equal", http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateTodoEntry(claims.AppID, claims.OrgID, claims.Subject, &item)
	if err != nil {
		log.Printf("Error on updating user todo entry with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal the updated user todo entry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// CreateUserTodoEntry Creates a user todo entry
// @Description Creates a user todo entry
// @Tags Client
// @ID CreateUserTodoEntry
// @Accept json
// @Success 200
// @Security UserAuth
// @Router /api/user/todo_entries [post]
func (h ApisHandler) CreateUserTodoEntry(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user todo entry - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item model.TodoEntry
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create user todo entry request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateTodoEntry(claims.AppID, claims.OrgID, claims.Subject, &item)
	if err != nil {
		log.Printf("Error on creating user todo entry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on marshal the new user todo entry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteUserTodoEntry Deletes a user todo entry with the specified id
// @Description Deletes a user todo entry with the specified id
// @Tags Client
// @ID DeleteUserTodoEntry
// @Success 200
// @Security UserAuth
// @Router /api/user/todo_entries/{id} [delete]
func (h ApisHandler) DeleteUserTodoEntry(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteTodoEntry(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on deleting user todo entry with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// DeleteCompletedUserTodoEntry Deletes all completed user todo entries
// @Description Deletes all completed user todo entries
// @Tags Client
// @ID DeleteCompletedUserTodoEntry
// @Success 200
// @Security UserAuth
// @Router /user/todo_entries/clear_completed_entries [delete]
func (h ApisHandler) DeleteCompletedUserTodoEntry(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteCompletedTodoEntries(claims.AppID, claims.OrgID, claims.Subject)
	if err != nil {
		log.Printf("Error on deleting user todo entry with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func intPostValueFromString(stringValue string) int {
	var value int
	if len(stringValue) > 0 {
		val, err := strconv.Atoi(stringValue)
		if err == nil {
			value = val
		}
	}
	return value
}

// NewApisHandler creates new rest Handler instance
func NewApisHandler(app *core.Application) ApisHandler {
	return ApisHandler{app: app}
}

// NewAdminApisHandler creates new rest Handler instance
func NewAdminApisHandler(app *core.Application) AdminApisHandler {
	return AdminApisHandler{app: app}
}

// NewInternalApisHandler creates new rest Handler instance
func NewInternalApisHandler(app *core.Application) InternalApisHandler {
	return InternalApisHandler{app: app}
}
