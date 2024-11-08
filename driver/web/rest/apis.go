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

package rest

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
	"wellness/core"
	"wellness/core/model"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rokwire/core-auth-library-go/v2/tokenauth"
)

const maxUploadSize = 15 * 1024 * 1024 // 15 mb

// ApisHandler handles the rest APIs implementation
type ApisHandler struct {
	app *core.Application
}

// Version gives the service version
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
// @Tags Client-TodoCategories
// @ID GetUserTodoCategories
// @Accept json
// @Success 200 {array} model.TodoCategory
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
// @Tags Client-TodoCategories
// @ID GetUserTodoCategory
// @Accept json
// @Produce json
// @Success 200 {object} model.TodoCategory
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
// @Tags Client-TodoCategories
// @ID UpdateUserTodoCategory
// @Accept json
// @Produce json
// @Param data body model.TodoCategory true "body json"
// @Success 200 {object} model.TodoCategory
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
// @Tags Client-TodoCategories
// @ID CreateUserTodoCategory
// @Param data body model.TodoCategory true "body json"
// @Success 200 {object} model.TodoCategory
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
// @Tags Client-TodoCategories
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
// @Tags Client-TodoEntries
// @ID GetUserTodoEntries
// @Accept json
// @Success 200 {array} model.TodoEntry
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
// @Tags Client-TodoEntries
// @ID GetUserTodoEntry
// @Produce json
// @Success 200 {object} model.TodoEntry
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
// @Tags Client-TodoEntries
// @ID UpdateUserTodoEntry
// @Accept json
// @Produce json
// @Param data body model.TodoEntry true "body json"
// @Success 200 {object} model.TodoEntry
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

	resData, err := h.app.Services.UpdateTodoEntry(claims.AppID, claims.OrgID, claims.Subject, &item, id)
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
// @Tags Client-TodoEntries
// @ID CreateUserTodoEntry
// @Accept json
// @Produce json
// @Param data body model.TodoEntry true "body json"
// @Success 200 {object} model.TodoEntry
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
// @Tags Client-TodoEntries
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
// @Tags Client-TodoEntries
// @ID DeleteCompletedUserTodoEntry
// @Success 200
// @Security UserAuth
// @Router /api/user/todo_entries/clear_completed_entries [delete]
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

// GetUserRings Retrieves all user wellness ring entries
// @Description Retrieves all user wellness ring entries
// @Tags Client-Rings
// @ID GetUserRings
// @Accept json
// @Success 200 {array} model.Ring
// @Security UserAuth
// @Router  /api/user/rings [get]
func (h ApisHandler) GetUserRings(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	resData, err := h.app.Services.GetRings(claims.AppID, claims.OrgID, claims.Subject)
	if err != nil {
		log.Printf("Error on getting user wellness ring entries - %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.Ring{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal all wellness ring todo entries: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserRing Retrieves a user wellness ring entry by id
// @Description Retrieves a user wellness ring entry by id
// @Tags Client-Rings
// @ID GetUserRing
// @Accept json
// @Produce json
// @Param data body model.Ring true "body json"
// @Success 200 {object} model.Ring
// @Security UserAuth
// @Router  /api/user/rings/{id} [get]
func (h ApisHandler) GetUserRing(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRing(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on getting user wellness ring entry by id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		log.Printf("Error on getting user wellness ring entry by id - %s (Not found)", id)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal user wellness ring entry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// createUserRingRequestBody represents request body data which is required for the initial state
type createUserRingRequestBody struct {
	Color string  `json:"color_hex" bson:"color_hex"`
	Name  string  `json:"name" bson:"name"`
	Unit  string  `json:"unit" bson:"unit"`
	Value float64 `json:"value" bson:"value"`
} // @name createUserRingRequestBody

// CreateUserRing Creates a user wellness ring entry
// @Description Creates a user wellness ring entry
// @Tags Client-Rings
// @ID CreateUserRing
// @Accept json
// @Param data body model.Ring true "body json"
// @Success 200 {object} model.Ring
// @Security UserAuth
// @Router /api/user/rings [post]
func (h ApisHandler) CreateUserRing(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user wellness ring entry - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var historyEntry createUserRingRequestBody
	err = json.Unmarshal(data, &historyEntry)
	if err != nil {
		log.Printf("Error on unmarshal the create user wellness ring entry request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRing(claims.AppID, claims.OrgID, claims.Subject, &model.Ring{
		History: []model.RingHistoryEntry{{
			ID:          uuid.NewString(),
			Color:       historyEntry.Color,
			Name:        historyEntry.Name,
			Unit:        historyEntry.Unit,
			Value:       historyEntry.Value,
			DateCreated: time.Now().UTC(),
		}},
	})
	if err != nil {
		log.Printf("Error on creating user wellness ring entry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on marshal the new user wellness ring entry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteUserRing Deletes a user wellness ring entry with the specified id
// @Description Deletes a user wellness ring entry with the specified id
// @Tags Client-Rings
// @ID DeleteUserRing
// @Security UserAuth
// @Router /api/user/rings/{id} [delete]
func (h ApisHandler) DeleteUserRing(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.app.Services.DeleteRing(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on deleting user wellness ring entry with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// CreateUserRingHistoryEntry Creates a user wellness ring history entry
// @Description Creates a user wellness ring history entry
// @Tags Client-Rings
// @ID CreateUserRingHistoryEntry
// @Accept json
// @Param data body createUserRingRequestBody true "body json"
// @Success 200
// @Security UserAuth
// @Router /api/user/rings/{id}/history [post]
func (h ApisHandler) CreateUserRingHistoryEntry(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user wellness ring history entry - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var historyEntry createUserRingRequestBody
	err = json.Unmarshal(data, &historyEntry)
	if err != nil {
		log.Printf("Error on unmarshal the create user wellness ring history entry request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.GetRing(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on creating user wellness ring history entry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		log.Printf("Error on creating user wellness ring history entry: %s\n Not found", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	createdItem, err := h.app.Services.CreateRingHistory(claims.AppID, claims.OrgID, claims.Subject, id, &model.RingHistoryEntry{
		ID:          uuid.NewString(),
		Color:       historyEntry.Color,
		Name:        historyEntry.Name,
		Unit:        historyEntry.Unit,
		Value:       historyEntry.Value,
		DateCreated: time.Now().UTC(),
	})
	if err != nil {
		log.Printf("Error on creating user wellness ring history entry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on marshal the new user wellness ring history entry: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteUserRingHistoryEntry Deletes a user wellness ring history entry with the specified id & history id
// @Description Deletes a user wellness ring history entry with the specified id & history id
// @Tags Client-Rings
// @ID DeleteUserRingHistoryEntry
// @Success 200
// @Security UserAuth
// @Router /api/user/rings/{id}/history/{history-id} [delete]
func (h ApisHandler) DeleteUserRingHistoryEntry(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	historyID := vars["history-id"]

	resData, err := h.app.Services.GetRing(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil {
		log.Printf("Error on deleting user wellness ring history entry: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		log.Printf("Error on deleting user wellness ring history entry: %s\n Not found", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = h.app.Services.DeleteRingHistory(claims.AppID, claims.OrgID, claims.Subject, id, historyID)
	if err != nil {
		log.Printf("Error on deleting user wellness ring history entry with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// GetUserAllRingRecords Retrieves all user ring record
// @Description Retrieves all user ring record
// @Tags Client-RingsRecords
// @ID GetUserAllRingRecords
// @Param offset query string false "offset"
// @Param limit query string false "limit - limit the result"
// @Param order query string false "order - Possible values: asc, desc. Default: desc"
// @Param start_date query string false "start_date - Start date filter in milliseconds as an integer epoch value"
// @Param end_date query string false "end_date - End date filter in milliseconds as an integer epoch value"
// @Success 200 {array} model.RingRecord
// @Security UserAuth
// @Router  /api/user/all_rings_records [get]
func (h ApisHandler) GetUserAllRingRecords(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	offsetFilter := getInt64QueryParam(r, "offset")
	limitFilter := getInt64QueryParam(r, "limit")
	orderFilter := getStringQueryParam(r, "order")
	startDateFilter := getInt64QueryParam(r, "start_date")
	endDateFilter := getInt64QueryParam(r, "end_date")

	resData, err := h.app.Services.GetRingsRecords(claims.AppID, claims.OrgID, claims.Subject, nil, startDateFilter, endDateFilter, offsetFilter, limitFilter, orderFilter)
	if err != nil {
		log.Printf("Error on getting user ring records- %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.RingRecord{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal all user ring records: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserRingRecords Retrieves all user ring record for a ring id
// @Description Retrieves all user ring record for a ring id
// @Tags Client-RingsRecords
// @ID GetUserRingRecords
// @Param offset query string false "offset"
// @Param limit query string false "limit - limit the result"
// @Param order query string false "order - Possible values: asc, desc. Default: desc"
// @Param start_date query string false "start_date - Start date filter in milliseconds as an integer epoch value"
// @Param end_date query string false "end_date - End date filter in milliseconds as an integer epoch value"
// @Success 200 {array} model.RingRecord
// @Security UserAuth
// @Router  /api/user/rings/{id}/records [get]
func (h ApisHandler) GetUserRingRecords(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	offsetFilter := getInt64QueryParam(r, "offset")
	limitFilter := getInt64QueryParam(r, "limit")
	orderFilter := getStringQueryParam(r, "order")
	startDateFilter := getInt64QueryParam(r, "start_date")
	endDateFilter := getInt64QueryParam(r, "end_date")
	vars := mux.Vars(r)
	id := vars["id"]

	resData, err := h.app.Services.GetRingsRecords(claims.AppID, claims.OrgID, claims.Subject, &id, startDateFilter, endDateFilter, offsetFilter, limitFilter, orderFilter)
	if err != nil {
		log.Printf("Error on getting user ring records- %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resData == nil {
		resData = []model.RingRecord{}
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal all user ring records: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// GetUserGetUserRingRecord Retrieves a user ring record by id
// @Description Retrieves a user ring record by id
// @Tags Client-RingsRecords
// @ID GetUserGetUserRingRecord
// @Produce json
// @Success 200 {array} model.RingRecord
// @Security UserAuth
// @Router  /api/user/rings/{id}/records/{record-id} [get]
func (h ApisHandler) GetUserGetUserRingRecord(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	recordID := vars["record-id"]

	resData, err := h.app.Services.GetRingsRecord(claims.AppID, claims.OrgID, claims.Subject, recordID)
	if err != nil {
		log.Printf("Error on getting user ring record by id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal user ring record: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// UpdateUserRingRecord Updates a user ring record with the specified id
// @Description Updates a user ring record with the specified id
// @Tags Client-RingsRecords
// @ID UpdateUserRingRecord
// @Accept json
// @Produce json
// @Param data body model.RingRecord true "body json"
// @Success 200 {array} model.RingRecord
// @Security UserAuth
// @Router /api/user/rings/{id}/records/{record-id} [put]
func (h ApisHandler) UpdateUserRingRecord(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	recordID := vars["record-id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user ring record - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ring, err := h.app.Services.GetRingsRecord(claims.AppID, claims.OrgID, claims.Subject, recordID)
	if err != nil || ring == nil {
		if err != nil {
			log.Printf("Error on marshal create a user ring record - ring not found %s\n", err)
		} else {
			log.Printf("Error on marshal create a user ring record - ring not found\n")

		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusBadRequest)
		return
	}

	var item model.RingRecord
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create user ring record request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if recordID != item.ID || item.RingID != id {
		log.Printf("Inconsistent attempt to update query param  and json id are not equal")
		http.Error(w, "Inconsistent attempt to update query param  and json id are not equal", http.StatusBadRequest)
		return
	}

	resData, err := h.app.Services.UpdateRingsRecord(claims.AppID, claims.OrgID, claims.Subject, &item)
	if err != nil {
		log.Printf("Error on updating user ring record with id - %s\n %s", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(resData)
	if err != nil {
		log.Printf("Error on marshal the updated user ring record: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// createUserRingRecordRequestBody represents individual daily record for an individual ring as a request body
type createUserRingRecordRequestBody struct {
	RingID string  `json:"ring_id" bson:"ring_id"`
	Value  float64 `json:"value" bson:"value"`
} //@name createUserRingRecordRequestBody

// CreateUserRingRecord Creates a user ring record
// @Description Creates a user ring record
// @Tags Client-RingsRecords
// @ID CreateUserRingRecord
// @Accept json
// @Param data body createUserRingRecordRequestBody true "body json"
// @Success 200 {array} model.RingRecord
// @Security UserAuth
// @Router /api/user/rings/{id}/records [post]
func (h ApisHandler) CreateUserRingRecord(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on marshal create a user ring record - %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ring, err := h.app.Services.GetRing(claims.AppID, claims.OrgID, claims.Subject, id)
	if err != nil || ring == nil {
		log.Printf("Error on marshal create a user ring record - ring not found %s\n", err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var item createUserRingRecordRequestBody
	err = json.Unmarshal(data, &item)
	if err != nil {
		log.Printf("Error on unmarshal the create user ring record request data - %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if item.RingID == "" {
		item.RingID = id
	}

	if item.RingID != id {
		log.Printf("api.CreateUserRingRecord() - ring id is different")
		http.Error(w, "api.CreateUserRingRecord() - ring id is different", http.StatusBadRequest)
		return
	}

	createdItem, err := h.app.Services.CreateRingsRecord(claims.AppID, claims.OrgID, claims.Subject, &model.RingRecord{
		RingID: item.RingID,
		Value:  item.Value,
	})
	if err != nil {
		log.Printf("Error on creating user ring record: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(createdItem)
	if err != nil {
		log.Printf("Error on marshal the new user ring record: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// DeleteAllUserRingRecords Deletes all user ring records (no matter of ring_id)
// @Description Deletes all user ring records (no matter of ring_id)
// @Tags Client-RingsRecords
// @ID DeleteAllUserRingRecords
// @Success 200
// @Security UserAuth
// @Router /api/user/all_rings_records [delete]
func (h ApisHandler) DeleteAllUserRingRecords(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	err := h.app.Services.DeleteRingsRecords(claims.AppID, claims.OrgID, claims.Subject, nil, nil)
	if err != nil {
		log.Printf("Error on deleting all user ring records - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// DeleteUserRingRecords Deletes all user ring record for a ring id
// @Description Deletes all user ring record for a ring id
// @Tags Client-RingsRecords
// @ID DeleteUserRingRecords
// @Success 200
// @Security UserAuth
// @Router /api/user/rings/{id}/records [delete]
func (h ApisHandler) DeleteUserRingRecords(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ringID := vars["id"]

	err := h.app.Services.DeleteRingsRecords(claims.AppID, claims.OrgID, claims.Subject, &ringID, nil)
	if err != nil {
		log.Printf("Error on deleting user ring records with ring_id - %s\n %s", ringID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// DeleteUserRingRecord Deletes a user ring record with the specified id
// @Description Deletes a user ring record with the specified id
// @Tags Client-RingsRecords
// @ID DeleteUserRingRecord
// @Success 200
// @Security UserAuth
// @Router /api/user/rings/{id}/records/{record-id} [delete]
func (h ApisHandler) DeleteUserRingRecord(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ringID := vars["id"]
	recordID := vars["record-id"]

	err := h.app.Services.DeleteRingsRecords(claims.AppID, claims.OrgID, claims.Subject, &ringID, &recordID)
	if err != nil {
		log.Printf("Error on deleting user ring record with id - %s\n %s", recordID, err)
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

// GetUserData Gets all related user data
// @Description  Gets all related user data
// @ID GetUserData
// @Tags Client
// @Success 200 {object} model.UserDataResponse
// @Security UserAuth
// @Router /api/user-data [get]
func (h ApisHandler) GetUserData(claims *tokenauth.Claims, w http.ResponseWriter, r *http.Request) {
	userData, err := h.app.Services.GetUserData(claims.AppID, claims.OrgID, claims.Subject)
	if err != nil {
		log.Printf("Error on creating user ring record: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Printf("Error on marshal the new user data: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
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
