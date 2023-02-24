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
	"time"
	"wellness/core/model"
	"wellness/driven/storage"
)

// Services exposes APIs for the driver adapters
type Services interface {
	GetVersion() string

	GetTodoCategories(appID string, orgID string, userID string) ([]model.TodoCategory, error)
	GetTodoCategory(appID string, orgID string, userID string, id string) (*model.TodoCategory, error)
	CreateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error)
	UpdateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error)
	DeleteTodoCategory(appID string, orgID string, userID string, id string) error

	GetTodoEntries(appID string, orgID string, userID string) ([]model.TodoEntry, error)
	GetTodoEntry(appID string, orgID string, userID string, id string) (*model.TodoEntry, error)
	CreateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry) error
	UpdateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry, id string) error
	DeleteTodoEntry(appID string, orgID string, userID string, id string) error
	DeleteCompletedTodoEntries(appID string, orgID string, userID string) error

	GetRings(appID string, orgID string, userID string) ([]model.Ring, error)
	GetRing(appID string, orgID string, userID string, id string) (*model.Ring, error)
	CreateRing(appID string, orgID string, userID string, category *model.Ring) (*model.Ring, error)
	DeleteRing(appID string, orgID string, userID string, id string) error
	CreateRingHistory(appID string, orgID string, userID string, ringID string, ringHistory *model.RingHistoryEntry) (*model.Ring, error)
	DeleteRingHistory(appID string, orgID string, userID string, ringID string, ringHistoryID string) (*model.Ring, error)

	GetRingsRecords(appID string, orgID string, userID string, ringID *string, startDateEpoch *int64, endDateEpoch *int64, offset *int64, limit *int64, order *string) ([]model.RingRecord, error)
	GetRingsRecord(appID string, orgID string, userID string, id string) (*model.RingRecord, error)
	CreateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error)
	UpdateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error)
	DeleteRingsRecords(appID string, orgID string, userID string, ringID *string, recordID *string) error
}

type servicesImpl struct {
	app *Application
}

func (s *servicesImpl) GetVersion() string {
	return s.app.getVersion()
}

func (s *servicesImpl) GetTodoCategories(appID string, orgID string, userID string) ([]model.TodoCategory, error) {
	return s.app.getTodoCategories(appID, orgID, userID)
}

func (s *servicesImpl) GetTodoCategory(appID string, orgID string, userID string, id string) (*model.TodoCategory, error) {
	return s.app.getTodoCategory(appID, orgID, userID, id)
}

func (s *servicesImpl) CreateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error) {
	return s.app.createTodoCategory(appID, orgID, userID, category)
}

func (s *servicesImpl) UpdateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error) {
	return s.app.updateTodoCategory(appID, orgID, userID, category)
}

func (s *servicesImpl) DeleteTodoCategory(appID string, orgID string, userID string, id string) error {
	return s.app.deleteTodoCategory(appID, orgID, userID, id)
}

func (s *servicesImpl) GetTodoEntries(appID string, orgID string, userID string) ([]model.TodoEntry, error) {
	return s.app.getTodoEntries(appID, orgID, userID)
}

func (s *servicesImpl) GetTodoEntry(appID string, orgID string, userID string, id string) (*model.TodoEntry, error) {
	return s.app.getTodoEntry(appID, orgID, userID, id)
}

func (s *servicesImpl) CreateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry) error /*(*model.TodoEntry, error)*/ {
	return s.app.createTodoEntry(appID, orgID, userID, todo)
}

func (s *servicesImpl) UpdateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry, id string) error {
	return s.app.updateTodoEntry(appID, orgID, userID, todo, id)
}

func (s *servicesImpl) DeleteTodoEntry(appID string, orgID string, userID string, id string) error {
	return s.app.deleteTodoEntry(appID, orgID, userID, id)
}

func (s *servicesImpl) DeleteCompletedTodoEntries(appID string, orgID string, userID string) error {
	return s.app.deleteCompletedTodoEntries(appID, orgID, userID)
}

func (s *servicesImpl) GetRings(appID string, orgID string, userID string) ([]model.Ring, error) {
	return s.app.getRings(appID, orgID, userID)
}

func (s *servicesImpl) GetRing(appID string, orgID string, userID string, id string) (*model.Ring, error) {
	return s.app.getRing(appID, orgID, userID, id)
}

func (s *servicesImpl) CreateRing(appID string, orgID string, userID string, category *model.Ring) (*model.Ring, error) {
	return s.app.createRing(appID, orgID, userID, category)
}

func (s *servicesImpl) DeleteRing(appID string, orgID string, userID string, id string) error {
	return s.app.deleteRing(appID, orgID, userID, id)
}

func (s *servicesImpl) CreateRingHistory(appID string, orgID string, userID string, ringID string, ringHistory *model.RingHistoryEntry) (*model.Ring, error) {
	return s.app.createRingHistory(appID, orgID, userID, ringID, ringHistory)
}

func (s *servicesImpl) DeleteRingHistory(appID string, orgID string, userID string, ringID string, ringHistoryID string) (*model.Ring, error) {
	return s.app.deleteRingHistory(appID, orgID, userID, ringID, ringHistoryID)
}

func (s *servicesImpl) GetRingsRecords(appID string, orgID string, userID string, ringID *string, startDateEpoch *int64, endDateEpoch *int64, offset *int64, limit *int64, order *string) ([]model.RingRecord, error) {
	return s.app.getRingsRecords(appID, orgID, userID, ringID, startDateEpoch, endDateEpoch, offset, limit, order)
}

func (s *servicesImpl) GetRingsRecord(appID string, orgID string, userID string, id string) (*model.RingRecord, error) {
	return s.app.getRingsRecord(appID, orgID, userID, id)
}

func (s *servicesImpl) CreateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error) {
	return s.app.createRingsRecord(appID, orgID, userID, record)
}

func (s *servicesImpl) UpdateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error) {
	return s.app.updateRingsRecord(appID, orgID, userID, record)
}

func (s *servicesImpl) DeleteRingsRecords(appID string, orgID string, userID string, ringID *string, recordID *string) error {
	return s.app.deleteRingsRecords(appID, orgID, userID, ringID, recordID)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	PerformTransaction(transaction func(context storage.TransactionContext) error) error

	GetTodoCategories(appID string, orgID string, userID string) ([]model.TodoCategory, error)
	GetTodoCategory(appID string, orgID string, userID string, id string) (*model.TodoCategory, error)
	CreateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error)
	UpdateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error)
	DeleteTodoCategory(appID string, orgID string, userID string, id string) error

	GetTodoEntriesWithCurrentReminderTime(context storage.TransactionContext, reminderTime time.Time) ([]model.TodoEntry, error)
	GetTodoEntriesWithCurrentDueTime(context storage.TransactionContext, dueTime time.Time) ([]model.TodoEntry, error)
	GetTodoEntries(appID string, orgID string, userID string) ([]model.TodoEntry, error)
	GetTodoEntry(appID string, orgID string, userID string, id string) (*model.TodoEntry, error)
	CreateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry) (*model.TodoEntry, error)
	UpdateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry, id string) (*model.TodoEntry, error)
	UpdateTodoEntriesTaskTime(context storage.TransactionContext, ids []string, taskTime time.Time) error
	DeleteTodoEntry(appID string, orgID string, userID string, id string) error
	DeleteCompletedTodoEntries(appID string, orgID string, userID string) error

	GetRings(appID string, orgID string, userID string) ([]model.Ring, error)
	GetRing(appID string, orgID string, userID string, id string) (*model.Ring, error)
	CreateRing(appID string, orgID string, userID string, category *model.Ring) (*model.Ring, error)
	DeleteRing(appID string, orgID string, userID string, id string) error
	CreateRingHistory(appID string, orgID string, userID string, ringID string, ringHistory *model.RingHistoryEntry) (*model.Ring, error)
	DeleteRingHistory(appID string, orgID string, userID string, ringID string, ringHistoryID string) (*model.Ring, error)

	GetRingsRecords(appID string, orgID string, userID string, ringID *string, startDateEpoch *int64, endDateEpoch *int64, offset *int64, limit *int64, order *string) ([]model.RingRecord, error)
	GetRingsRecord(appID string, orgID string, userID string, id string) (*model.RingRecord, error)
	CreateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error)
	UpdateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error)
	DeleteRingsRecords(appID string, orgID string, userID string, ringID *string, recordID *string) error
}

// Notifications wrapper
type Notifications interface {
	SendNotification(recipients []model.NotificationRecipient, topic *string, title string, text string, appID string, orgID string, time *int64, data map[string]string) (*string, error)
	DeleteNotification(appID string, orgID string, id string) error
}
