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
	"fmt"
	"log"
	"time"
	"wellness/core/model"
)

func (app *Application) getVersion() string {
	return app.version
}

func (app *Application) getTodoCategories(appID string, orgID string, userID string) ([]model.TodoCategory, error) {
	return app.storage.GetTodoCategories(appID, orgID, userID)
}

func (app *Application) getTodoCategory(appID string, orgID string, userID string, id string) (*model.TodoCategory, error) {
	return app.storage.GetTodoCategory(appID, orgID, userID, id)
}

func (app *Application) createTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error) {
	return app.storage.CreateTodoCategory(appID, orgID, userID, category)
}

func (app *Application) updateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error) {
	return app.storage.UpdateTodoCategory(appID, orgID, userID, category)
}

func (app *Application) deleteTodoCategory(appID string, orgID string, userID string, id string) error {
	return app.storage.DeleteTodoCategory(appID, orgID, userID, id)
}

func (app *Application) getTodoEntries(appID string, orgID string, userID string) ([]model.TodoEntry, error) {
	return app.storage.GetTodoEntries(appID, orgID, userID)
}

func (app *Application) getTodoEntry(appID string, orgID string, userID string, id string) (*model.TodoEntry, error) {
	return app.storage.GetTodoEntry(appID, orgID, userID, id)
}

func (app *Application) createTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error) {
	return app.storage.CreateTodoEntry(appID, orgID, userID, category)
}

func (app *Application) updateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error) {
	return app.storage.UpdateTodoEntry(appID, orgID, userID, category)
}

func (app *Application) deleteTodoEntry(appID string, orgID string, userID string, id string) error {
	return app.storage.DeleteTodoEntry(appID, orgID, userID, id)
}

func (app *Application) deleteCompletedTodoEntries(appID string, orgID string, userID string) error {
	return app.storage.DeleteCompletedTodoEntries(appID, orgID, userID)
}

func (app *Application) processReminders() error {
	now := time.Now()
	todos, err := app.storage.GetTodoEntriesWithCurrentDueTime(now)
	if err != nil {
		log.Printf("Error on retrieving reminders: %s", err)
	}

	topic := "wellness.reminders"
	if len(todos) > 0 {
		for _, todo := range todos {
			err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: todo.UserID}}, &topic, fmt.Sprintf("TODO: %s", todo.Title), todo.Description, map[string]string{
				"type":        "wellness_todo_entry",
				"operation":   "todo_reminder",
				"entity_type": "wellness_todo_entry",
				"entity_id":   todo.ID,
				"entity_name": todo.Title,
			})
			if err != nil {
				log.Printf("Error on sending reminder inbox message: %s", err)
			}
		}
	}

	todos, err = app.storage.GetTodoEntriesWithCurrentDueTime(now)
	if err != nil {
		log.Printf("Error on retrieving due time reminders: %s", err)
	}

	if len(todos) > 0 {
		for _, todo := range todos {
			err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: todo.UserID}}, &topic, fmt.Sprintf("TODO: %s", todo.Title), todo.Description, map[string]string{
				"type":        "wellness_todo_entry",
				"operation":   "todo_reminder",
				"entity_type": "wellness_todo_entry",
				"entity_id":   todo.ID,
				"entity_name": todo.Title,
			})
			if err != nil {
				log.Printf("Error on sending reminder inbox message: %s", err)
			}
		}
	}

	return nil
}
