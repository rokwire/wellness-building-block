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
	"time"
	"wellness/core/model"
	"wellness/driven/storage"
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

func (app *Application) createTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry) error {
	return app.storage.PerformTransaction(func(context storage.TransactionContext) error {
		createTodoEntry, err := app.storage.CreateTodoEntry(appID, orgID, userID, todo)
		if err != nil {
			log.Printf("Error on retrieving reminders: %s", err)
		}
		topic := "create todo entry"
		dueDateTime := createTodoEntry.DueDateTime.Unix()
		_, err = app.notifications.SendNotification([]model.NotificationRecipient{{UserID: createTodoEntry.UserID}}, &topic, "TODO Reminder", createTodoEntry.Title, createTodoEntry.AppID, createTodoEntry.OrgID, &dueDateTime, map[string]string{
			"type":        "wellness_todo_entry",
			"operation":   "todo_reminder",
			"entity_type": "wellness_todo_entry",
			"entity_id":   createTodoEntry.ID,
			"entity_name": createTodoEntry.Title,
		})

		reminderDateTime := createTodoEntry.ReminderDateTime.Unix()
		_, err = app.notifications.SendNotification([]model.NotificationRecipient{{UserID: createTodoEntry.UserID}}, &topic, "TODO Reminder", createTodoEntry.Title, createTodoEntry.AppID, createTodoEntry.OrgID, &reminderDateTime, map[string]string{
			"type":        "wellness_todo_entry",
			"operation":   "todo_reminder",
			"entity_type": "wellness_todo_entry",
			"entity_id":   createTodoEntry.ID,
			"entity_name": createTodoEntry.Title,
		})
		if err != nil {
			log.Printf("Error on sending notification %s inbox message: %s", createTodoEntry.ID, err)
		} else {
			log.Printf("Sent notification %s successfully", createTodoEntry.ID)
		}
		return nil
	})
}

func (app *Application) updateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry, id string) error {
	return app.storage.PerformTransaction(func(context storage.TransactionContext) error {
		updateTodoEntry, err := app.storage.UpdateTodoEntry(appID, orgID, userID, todo, id)
		if err != nil {
			log.Printf("Error on retrieving reminders: %s", err)
		}
		topic := "update todo entry"
		_, err = app.notifications.SendNotification([]model.NotificationRecipient{{UserID: updateTodoEntry.UserID}}, &topic, "TODO Reminder", updateTodoEntry.Title, updateTodoEntry.AppID, updateTodoEntry.OrgID, nil, map[string]string{
			"type":        "wellness_todo_entry",
			"operation":   "todo_reminder",
			"entity_type": "wellness_todo_entry",
			"entity_id":   updateTodoEntry.ID,
			"entity_name": updateTodoEntry.Title,
		})
		if err != nil {
			log.Printf("Error on sending notification %s inbox message: %s", updateTodoEntry.ID, err)
		} else {
			log.Printf("Sent notification %s successfully", updateTodoEntry.ID)
		}
		return nil
	})
}

func (app *Application) deleteTodoEntry(appID string, orgID string, userID string, id string) error {

	return app.storage.PerformTransaction(func(context storage.TransactionContext) error {
		err := app.storage.DeleteTodoEntry(appID, orgID, userID, id)
		if err != nil {
			log.Printf("Error on retrieving reminders: %s", err)
		}

		topic := "delete todo entry"
		title := "delete todo entry"
		_, err = app.notifications.SendNotification([]model.NotificationRecipient{{UserID: userID}}, &topic, "TODO Reminder", title, appID, orgID, nil, map[string]string{
			"type":        "wellness_todo_entry",
			"operation":   "todo_reminder",
			"entity_type": "wellness_todo_entry",
			"entity_id":   id,
			"entity_name": title,
		})
		if err != nil {
			log.Printf("Error on sending notification %s inbox message: %s", id, err)
		} else {
			log.Printf("Sent notification %s successfully", id)
		}
		return nil
	})
}

func (app *Application) deleteCompletedTodoEntries(appID string, orgID string, userID string) error {
	return app.storage.DeleteCompletedTodoEntries(appID, orgID, userID)
}

func (app *Application) processReminders() error {

	return app.storage.PerformTransaction(func(context storage.TransactionContext) error {
		log.Printf("Start reminder processing")
		now := time.Now()
		todos, err := app.storage.GetTodoEntriesWithCurrentDueTime(context, now)
		if err != nil {
			log.Printf("Error on retrieving reminders: %s", err)
		}

		todoCount := len(todos)
		if todoCount > 0 {
			todoIDs := make([]string, todoCount)
			for index := range todos {
				todoIDs[index] = todos[index].ID
			}
			err := app.storage.UpdateTodoEntriesTaskTime(context, todoIDs, now)
			if err != nil {
				log.Printf("Error on updating reminders task time: %s", err)
			}
		}

		topic := "wellness.reminders"
		if len(todos) > 0 {
			for _, todo := range todos {
				_, err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: todo.UserID}}, &topic, "TODO Reminder", todo.Title, todo.AppID, todo.OrgID, nil, map[string]string{
					"type":        "wellness_todo_entry",
					"operation":   "todo_reminder",
					"entity_type": "wellness_todo_entry",
					"entity_id":   todo.ID,
					"entity_name": todo.Title,
				})
				if err != nil {
					log.Printf("Error on sending reminder %s inbox message: %s", todo.ID, err)
				} else {
					log.Printf("Sent notification for reminder %s successfully", todo.ID)
				}
			}
		}
		log.Printf("Processed %d reminders", len(todos))

		todos, err = app.storage.GetTodoEntriesWithCurrentReminderTime(context, now)
		if err != nil {
			log.Printf("Error on retrieving due time reminders: %s", err)
		}

		todoCount = len(todos)
		if todoCount > 0 {
			todoIDs := make([]string, todoCount)
			for index := range todos {
				todoIDs[index] = todos[index].ID
			}
			err := app.storage.UpdateTodoEntriesTaskTime(context, todoIDs, now)
			if err != nil {
				log.Printf("Error on updating reminders task time: %s", err)
			}
		}

		if len(todos) > 0 {
			for _, todo := range todos {
				_, err := app.notifications.SendNotification([]model.NotificationRecipient{{UserID: todo.UserID}}, &topic, "TODO Reminder", todo.Title, todo.AppID, todo.OrgID, nil, map[string]string{
					"type":        "wellness_todo_entry",
					"operation":   "todo_reminder",
					"entity_type": "wellness_todo_entry",
					"entity_id":   todo.ID,
					"entity_name": todo.Title,
				})
				if err != nil {
					log.Printf("Error on sending reminder %s inbox message: %s", todo.ID, err)
				} else {
					log.Printf("Sent notification for reminder %s successfully", todo.ID)
				}
			}
		}
		log.Printf("Processed %d due date reminders", len(todos))

		log.Printf("End reminder processing")

		return nil
	})

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
