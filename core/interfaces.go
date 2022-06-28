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
	"wellness/core/model"
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
	CreateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error)
	UpdateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error)
	DeleteTodoEntry(appID string, orgID string, userID string, id string) error
	DeleteCompletedTodoEntries(appID string, orgID string, userID string) error
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

func (s *servicesImpl) CreateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error) {
	return s.app.createTodoEntry(appID, orgID, userID, category)
}

func (s *servicesImpl) UpdateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error) {
	return s.app.updateTodoEntry(appID, orgID, userID, category)
}

func (s *servicesImpl) DeleteTodoEntry(appID string, orgID string, userID string, id string) error {
	return s.app.deleteTodoEntry(appID, orgID, userID, id)
}

func (s *servicesImpl) DeleteCompletedTodoEntries(appID string, orgID string, userID string) error {
	return s.app.deleteCompletedTodoEntries(appID, orgID, userID)
}

// Storage is used by core to storage data - DB storage adapter, file storage adapter etc
type Storage interface {
	GetTodoCategories(appID string, orgID string, userID string) ([]model.TodoCategory, error)
	GetTodoCategory(appID string, orgID string, userID string, id string) (*model.TodoCategory, error)
	CreateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error)
	UpdateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error)
	DeleteTodoCategory(appID string, orgID string, userID string, id string) error

	GetTodoEntries(appID string, orgID string, userID string) ([]model.TodoEntry, error)
	GetTodoEntry(appID string, orgID string, userID string, id string) (*model.TodoEntry, error)
	CreateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error)
	UpdateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error)
	DeleteTodoEntry(appID string, orgID string, userID string, id string) error
	DeleteCompletedTodoEntries(appID string, orgID string, userID string) error
}
