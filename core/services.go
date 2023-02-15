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

func (app *Application) createTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry) (*model.TodoEntry, error) {
	return app.storage.CreateTodoEntry(appID, orgID, userID, todo)
}

func (app *Application) updateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry) (*model.TodoEntry, error) {
	return app.storage.UpdateTodoEntry(appID, orgID, userID, todo)
}

func (app *Application) deleteTodoEntry(appID string, orgID string, userID string, id string) error {
	return app.storage.DeleteTodoEntry(appID, orgID, userID, id)
}

func (app *Application) deleteCompletedTodoEntries(appID string, orgID string, userID string) error {
	return app.storage.DeleteCompletedTodoEntries(appID, orgID, userID)
}

func (app *Application) processReminders() error {

	log.Printf("Start reminder processing")
	now := time.Now()
	todos, err := app.storage.GetTodoEntriesWithCurrentDueTime(now)
	if err != nil {
		log.Printf("Error on retrieving reminders: %s", err)
	}
	//topic temporarly removed
	//	topic := "wellness.reminders"
	if len(todos) > 0 {
		for _, todo := range todos {
			err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: todo.UserID}}, nil, fmt.Sprintf("TODO Reminder: %s", todo.Title), todo.Description, todo.AppID, todo.OrgID, map[string]string{
				"type":        "wellness_todo_entry",
				"operation":   "todo_reminder",
				"entity_type": "wellness_todo_entry",
				"entity_id":   todo.ID,
				"entity_name": todo.Title,
				"app_id":      todo.AppID,
				"org_id":      todo.OrgID,
			})
			if err != nil {
				log.Printf("Error on sending reminder inbox message: %s", err)
			}
		}
	}
	log.Printf("Processed %d reminders", len(todos))

	todos, err = app.storage.GetTodoEntriesWithCurrentReminderTime(now)
	if err != nil {
		log.Printf("Error on retrieving due time reminders: %s", err)
	}

	if len(todos) > 0 {
		for _, todo := range todos {
			err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: todo.UserID}}, nil, fmt.Sprintf("TODO: %s", todo.Title), todo.Description, todo.AppID, todo.OrgID, map[string]string{
				"type":        "wellness_todo_entry",
				"operation":   "todo_reminder",
				"entity_type": "wellness_todo_entry",
				"entity_id":   todo.ID,
				"entity_name": todo.Title,
				"app_id":      todo.AppID,
				"org_id":      todo.OrgID,
			})
			if err != nil {
				log.Printf("Error on sending reminder inbox message: %s", err)
			}
		}
	}
	log.Printf("Processed %d due date reminders", len(todos))

	log.Printf("End reminder processing")

	return nil
}

func (app *Application) getRings(appID string, orgID string, userID string) ([]model.Ring, error) {
	return app.storage.GetRings(appID, orgID, userID)
}

func (app *Application) getRing(appID string, orgID string, userID string, id string) (*model.Ring, error) {
	return app.storage.GetRing(appID, orgID, userID, id)
}

func (app *Application) createRing(appID string, orgID string, userID string, category *model.Ring) (*model.Ring, error) {
	return app.storage.CreateRing(appID, orgID, userID, category)
}

func (app *Application) deleteRing(appID string, orgID string, userID string, id string) error {
	return app.storage.DeleteRing(appID, orgID, userID, id)
}

func (app *Application) createRingHistory(appID string, orgID string, userID string, ringID string, ringHistory *model.RingHistoryEntry) (*model.Ring, error) {
	return app.storage.CreateRingHistory(appID, orgID, userID, ringID, ringHistory)
}

func (app *Application) deleteRingHistory(appID string, orgID string, userID string, ringID string, ringHistoryID string) (*model.Ring, error) {
	return app.storage.DeleteRingHistory(appID, orgID, userID, ringID, ringHistoryID)
}

func (app *Application) getRingsRecords(appID string, orgID string, userID string, ringID *string, startDateEpoch *int64, endDateEpoch *int64, offset *int64, limit *int64, order *string) ([]model.RingRecord, error) {
	return app.storage.GetRingsRecords(appID, orgID, userID, ringID, startDateEpoch, endDateEpoch, offset, limit, order)
}

func (app *Application) getRingsRecord(appID string, orgID string, userID string, id string) (*model.RingRecord, error) {
	return app.storage.GetRingsRecord(appID, orgID, userID, id)
}

func (app *Application) createRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error) {
	return app.storage.CreateRingsRecord(appID, orgID, userID, record)
}

func (app *Application) updateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error) {
	return app.storage.UpdateRingsRecord(appID, orgID, userID, record)
}

func (app *Application) deleteRingsRecords(appID string, orgID string, userID string, ringID *string, recordID *string) error {
	return app.storage.DeleteRingsRecords(appID, orgID, userID, ringID, recordID)
}
