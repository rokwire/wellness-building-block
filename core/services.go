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
	"wellness/core/model"
	"wellness/driven/storage"

	"github.com/google/uuid"
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
	var createTodoEntry *model.TodoEntry
	entityID := uuid.NewString()
	err := app.storage.PerformTransaction(func(context storage.TransactionContext) error {
		topic := "create todo entry"
		var dueMsgID *string
		var reminderMsgID *string
		var err error
		if todo.DueDateTime != nil {
			dueDateTime := todo.DueDateTime.Unix()
			dueMsgID, err = app.notifications.SendNotification([]model.NotificationRecipient{{UserID: userID}}, &topic, "TODO Reminder", todo.Title, appID, orgID, &dueDateTime, map[string]string{
				"type":        "wellness_todo_entry",
				"operation":   "todo_reminder",
				"entity_type": "wellness_todo_entry",
				"entity_id":   entityID,
				"entity_name": todo.Title,
			})

			if err != nil {
				log.Printf("Error on sending DueDateTime notification %s inbox message: %s", todo.ID, err)
				//return err // Don't propagate the error. Just create the reminder.
			}
			log.Printf("Sent DueDateTime notification %s successfully", entityID)
		}
		if todo.ReminderDateTime != nil {
			reminderDateTime := todo.ReminderDateTime.Unix()
			reminderMsgID, err = app.notifications.SendNotification([]model.NotificationRecipient{{UserID: userID}}, &topic, "TODO Reminder", todo.Title, appID, orgID, &reminderDateTime, map[string]string{
				"type":        "wellness_todo_entry",
				"operation":   "todo_reminder",
				"entity_type": "wellness_todo_entry",
				"entity_id":   entityID,
				"entity_name": todo.Title,
			})
			if err != nil {
				log.Printf("Error on sending ReminderDateTime notification %s inbox message: %s", todo.ID, err)
				//return err // Don't propagate the error. Just create the reminder.
			}
			log.Printf("Sent ReminderDateTime notification %s successfully", entityID)
		}

		createTodoEntry, err = app.storage.CreateTodoEntry(appID, orgID, userID, todo, model.MessageIDs{ReminderDateMessageID: reminderMsgID, DueDateMessageID: dueMsgID}, entityID)
		if err != nil {
			log.Printf("Error on creating todo entry: %s", err)
		}
		return nil
	})
	return createTodoEntry, err
}

func (app *Application) updateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry, id string) (*model.TodoEntry, error) {
	var updateTodoEntry *model.TodoEntry
	err := app.storage.PerformTransaction(func(context storage.TransactionContext) error {
		todoEntry, err := app.storage.GetTodoEntry(appID, orgID, userID, id)
		if err != nil {
			log.Printf("Error on getting todo entry: %s", err)
		}
		if todoEntry.MessageIDs.DueDateMessageID != nil {
			err = app.notifications.DeleteNotification(appID, orgID, *todoEntry.MessageIDs.DueDateMessageID)
			if err != nil {
				log.Printf("Error on delete notification with DueDateMessageID %s", todoEntry.MessageIDs.DueDateMessageID)
				//return err // Don't propagate the error. Just create the reminder.
			}
		}

		if todo.DueDateTime != nil {
			if todoEntry.DueDateTime != todo.DueDateTime {
				topic := "update due date time"
				dueDateTime := todo.DueDateTime.Unix()
				duoMsg, err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: userID}}, &topic, "TODO Reminder", todo.Title, appID, orgID, &dueDateTime, map[string]string{
					"type":        "wellness_todo_entry",
					"operation":   "todo_reminder",
					"entity_type": "wellness_todo_entry",
					"entity_id":   todo.ID,
					"entity_name": todo.Title,
				})

				if err != nil {
					log.Printf("Error on sending DueDateTime notification %s inbox message: %s", id, err)
					//return err // Don't propagate the error. Just create the reminder.
				}
				todo.MessageIDs.DueDateMessageID = duoMsg
				log.Printf("Sent DueDateTime notification %s successfully", id)
			}
		}

		if todoEntry.MessageIDs.ReminderDateMessageID != nil {
			err = app.notifications.DeleteNotification(appID, orgID, *todoEntry.MessageIDs.ReminderDateMessageID)
			if err != nil {
				log.Printf("Error on delete notification with ReminderDateMessageID %s", todoEntry.MessageIDs.ReminderDateMessageID)
				//return err // Don't propagate the error. Just create the reminder.
			}
		}
		if todo.ReminderDateTime != nil {
			if todoEntry.ReminderDateTime != todo.ReminderDateTime {
				topic := "update due date time"
				reminderDateTime := todo.ReminderDateTime.Unix()
				reminderMsg, err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: userID}}, &topic, "TODO Reminder", todo.Title, appID, orgID, &reminderDateTime, map[string]string{
					"type":        "wellness_todo_entry",
					"operation":   "todo_reminder",
					"entity_type": "wellness_todo_entry",
					"entity_id":   todo.ID,
					"entity_name": todo.Title,
				})

				if err != nil {
					log.Printf("Error on sending ReminderDateTime notification %s inbox message: %s", id, err)
					//return err // Don't propagate the error. Just create the reminder.
				}

				todo.MessageIDs.ReminderDateMessageID = reminderMsg
				log.Printf("Sent ReminderDateTime notification %s successfully", id)
			}
		}

		updateTodoEntry, err = app.storage.UpdateTodoEntry(appID, orgID, userID, todo, id)
		if err != nil {
			log.Printf("Error on updating todo entry: %s", err)
		}

		return nil
	})
	return updateTodoEntry, err
}

func (app *Application) deleteTodoEntry(appID string, orgID string, userID string, id string) error {

	return app.storage.PerformTransaction(func(context storage.TransactionContext) error {
		todoEntry, err := app.storage.GetTodoEntry(appID, orgID, userID, id)
		if err != nil {
			log.Printf("Error on getting todo entry: %s", err)
		}
		if todoEntry.MessageIDs.DueDateMessageID != nil {
			err = app.notifications.DeleteNotification(appID, orgID, *todoEntry.MessageIDs.DueDateMessageID)
			if err != nil {
				log.Printf("Error on delete DueDateMessageID notificarion with id %s", *todoEntry.MessageIDs.DueDateMessageID)
				//return err // Don't propagate the error. Just create the reminder.
			}
		}
		if todoEntry.MessageIDs.ReminderDateMessageID != nil {
			err = app.notifications.DeleteNotification(appID, orgID, *todoEntry.MessageIDs.ReminderDateMessageID)
			if err != nil {
				log.Printf("Error on delete ReminderDateMessageID notificarion with id %s", *todoEntry.MessageIDs.ReminderDateMessageID)
				//return err // Don't propagate the error. Just create the reminder.
			}
		}
		err = app.storage.DeleteTodoEntry(appID, orgID, userID, id)
		if err != nil {
			log.Printf("Error on delete todo entry: %s", err)
		}

		return nil
	})
}

func (app *Application) deleteCompletedTodoEntries(appID string, orgID string, userID string) error {
	return app.storage.DeleteCompletedTodoEntries(appID, orgID, userID)
}

// MigrateMessageIDs migrate message ids
func (app *Application) MigrateMessageIDs() error {
	transaction := func(context storage.TransactionContext) error {
		todoEntries, err := app.storage.GetTodoEntriesForMigration()
		if err != nil {
			log.Printf("error on getting todo entries - %s", err)
		}

		for _, todo := range todoEntries {
			if todo.RequiresMessageIDsMigration() {
				_, err := app.updateTodoEntry(todo.AppID, todo.OrgID, todo.UserID, &todo, todo.ID)
				if err != nil {
					log.Printf("error on updating todo entries - %s", err)
				}

			}
		}
		return nil
	}
	return app.storage.PerformTransaction(transaction)
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
