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

package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"time"
	"wellness/core/model"

	"github.com/rokwire/logging-library-go/errors"
	"github.com/rokwire/logging-library-go/logutils"
	"go.mongodb.org/mongo-driver/mongo"
)

// Adapter implements the Storage interface
type Adapter struct {
	db *database
}

// Start starts the storage
func (sa *Adapter) Start() error {
	err := sa.db.start()
	return err
}

//PerformTransaction performs a transaction
func (sa *Adapter) PerformTransaction(transaction func(context TransactionContext) error) error {
	// transaction
	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			sa.abortTransaction(sessionContext)
			return errors.WrapErrorAction(logutils.ActionStart, logutils.TypeTransaction, nil, err)
		}

		err = transaction(sessionContext)
		if err != nil {
			sa.abortTransaction(sessionContext)
			return errors.WrapErrorAction("performing", logutils.TypeTransaction, nil, err)
		}

		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			sa.abortTransaction(sessionContext)
			return errors.WrapErrorAction(logutils.ActionCommit, logutils.TypeTransaction, nil, err)
		}
		return nil
	})

	return err
}

// GetTodoCategories gets all user defined todo categories
func (sa *Adapter) GetTodoCategories(appID string, orgID string, userID string) ([]model.TodoCategory, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "user_id", Value: userID},
	}

	var result []model.TodoCategory
	err := sa.db.todoCategories.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTodoCategory gets a single user defined todo category by id
func (sa *Adapter) GetTodoCategory(appID string, orgID string, userID string, id string) (*model.TodoCategory, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: id},
	}

	var result []model.TodoCategory
	err := sa.db.todoCategories.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return &result[0], nil
	}

	return nil, nil
}

// CreateTodoCategory create a new user defined todo category
func (sa *Adapter) CreateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error) {
	category.ID = uuid.NewString()
	category.OrgID = orgID
	category.AppID = appID
	category.UserID = userID
	category.DateCreated = time.Now().UTC()

	_, err := sa.db.todoCategories.InsertOne(&category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateTodoCategory updates a user defined todo category
func (sa *Adapter) UpdateTodoCategory(appID string, orgID string, userID string, category *model.TodoCategory) (*model.TodoCategory, error) {

	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		filter := bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "user_id", Value: userID},
			primitive.E{Key: "_id", Value: category.ID}}
		update := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "name", Value: category.Name},
				primitive.E{Key: "color", Value: category.Color},
				primitive.E{Key: "date_updated", Value: time.Now().UTC()},
			}},
		}
		_, err = sa.db.todoCategories.UpdateOneWithContext(sessionContext, filter, update, nil)
		if err != nil {
			sa.abortTransaction(sessionContext)
			log.Printf("error updating user defined todo category: %s", err)
			return err
		}

		filter = bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "user_id", Value: userID},
			primitive.E{Key: "category.id", Value: category.ID}}
		update = bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "category", Value: category.ToCategoryRef()},
			}},
		}
		_, err = sa.db.todoEntries.UpdateManyWithContext(sessionContext, filter, update, nil)
		if err != nil {
			sa.abortTransaction(sessionContext)
			log.Printf("error updating user defined todo category: %s", err)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("error on update category: %s", err)
		return nil, fmt.Errorf("error on update category: %s", err)
	}

	return sa.GetTodoCategory(appID, orgID, userID, category.ID)
}

// DeleteTodoCategory deletes a user defined todo category
func (sa *Adapter) DeleteTodoCategory(appID string, orgID string, userID string, id string) error {

	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		filter := bson.D{primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "user_id", Value: userID},
			primitive.E{Key: "_id", Value: id}}

		_, err = sa.db.todoCategories.DeleteOneWithContext(sessionContext, filter, nil)
		if err != nil {
			sa.abortTransaction(sessionContext)
			log.Printf("error deleting todo category: %s", err)
			return err
		}

		filter = bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "user_id", Value: userID},
			primitive.E{Key: "category.id", Value: id}}
		update := bson.D{
			primitive.E{Key: "$set", Value: bson.D{
				primitive.E{Key: "category", Value: bson.TypeNull},
			}},
		}
		_, err = sa.db.todoEntries.UpdateManyWithContext(sessionContext, filter, update, nil)
		if err != nil {
			sa.abortTransaction(sessionContext)
			log.Printf("error deleting todo category: %s", err)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("error on delete todo category: %s", err)
		return fmt.Errorf("error on delete todo category: %s", err)
	}

	return nil
}

// GetTodoEntries gets user's todo entries
func (sa *Adapter) GetTodoEntries(appID string, orgID string, userID string) ([]model.TodoEntry, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "user_id", Value: userID},
	}

	var result []model.TodoEntry
	err := sa.db.todoEntries.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTodoEntry get a single todo entry
func (sa *Adapter) GetTodoEntry(appID string, orgID string, userID string, id string) (*model.TodoEntry, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: id},
	}

	var result []model.TodoEntry
	err := sa.db.todoEntries.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return &result[0], nil
	}

	return nil, nil
}

// CreateTodoEntry create a todo entry
func (sa *Adapter) CreateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry) (*model.TodoEntry, error) {
	category.ID = uuid.NewString()
	category.OrgID = orgID
	category.AppID = appID
	category.UserID = userID
	category.DateCreated = time.Now().UTC()

	_, err := sa.db.todoEntries.InsertOne(&category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateTodoEntry updates a todo entry
func (sa *Adapter) UpdateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry) (*model.TodoEntry, error) {

	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: todo.ID}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "title", Value: todo.Title},
			primitive.E{Key: "description", Value: todo.Description},
			primitive.E{Key: "category", Value: todo.Category},
			primitive.E{Key: "completed", Value: todo.Completed},
			primitive.E{Key: "has_due_time", Value: todo.HasDueTime},
			primitive.E{Key: "due_date_time", Value: todo.DueDateTime},
			primitive.E{Key: "reminder_type", Value: todo.ReminderType},
			primitive.E{Key: "reminder_date_time", Value: todo.ReminderDateTime},
			primitive.E{Key: "work_days", Value: todo.WorkDays},
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err := sa.db.todoEntries.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("error updating user defined todo entry: %s", err)
		return nil, err
	}

	return sa.GetTodoEntry(appID, orgID, userID, todo.ID)
}

// DeleteTodoEntry deletes a todo entry
func (sa *Adapter) DeleteTodoEntry(appID string, orgID string, userID string, id string) error {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: id}}

	_, err := sa.db.todoEntries.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("error deleting todo entry: %s", err)
		return err
	}

	return nil
}

// DeleteCompletedTodoEntries deletes a completed todo entries
func (sa *Adapter) DeleteCompletedTodoEntries(appID string, orgID string, userID string) error {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "completed", Value: true}}

	_, err := sa.db.todoEntries.DeleteMany(filter, nil)
	if err != nil {
		log.Printf("error deleting comleted todo entries: %s", err)
		return err
	}

	return nil
}

// GetTodoEntriesWithCurrentReminderTime Gets all todo entries that are applied for the specified reminder datetime
func (sa *Adapter) GetTodoEntriesWithCurrentReminderTime(reminderTime time.Time) ([]model.TodoEntry, error) {
	startDate := time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, reminderTime.Location())
	endDate := time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), reminderTime.Hour(), reminderTime.Minute(), 59, 999999, reminderTime.Location())
	filter := bson.D{
		primitive.E{Key: "completed", Value: false},
		primitive.E{Key: "reminder_date_time", Value: []primitive.E{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}

	var result []model.TodoEntry
	err := sa.db.todoEntries.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTodoEntriesWithCurrentDueTime Gets all todo entries that are applied for the specified due datetime
func (sa *Adapter) GetTodoEntriesWithCurrentDueTime(dueTime time.Time) ([]model.TodoEntry, error) {
	startDate := time.Date(dueTime.Year(), dueTime.Month(), dueTime.Day(), dueTime.Hour(), dueTime.Minute(), 0, 0, dueTime.Location())
	endDate := time.Date(dueTime.Year(), dueTime.Month(), dueTime.Day(), dueTime.Hour(), dueTime.Minute(), 59, 999999, dueTime.Location())
	filter := bson.D{
		primitive.E{Key: "completed", Value: false},
		primitive.E{Key: "has_due_time", Value: true},
		primitive.E{Key: "due_date_time", Value: []primitive.E{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
	}

	var result []model.TodoEntry
	err := sa.db.todoEntries.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Wellness Rings

// GetRings gets user's wellness rings
func (sa *Adapter) GetRings(appID string, orgID string, userID string) ([]model.Ring, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "user_id", Value: userID},
	}

	var result []model.Ring
	err := sa.db.rings.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetRing get a single user wellness ring
func (sa *Adapter) GetRing(appID string, orgID string, userID string, id string) (*model.Ring, error) {
	filter := bson.D{
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: id},
	}

	var result []model.Ring
	err := sa.db.rings.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}

	if len(result) > 0 {
		return &result[0], nil
	}

	return nil, nil
}

// CreateRing create a user wellness ring
func (sa *Adapter) CreateRing(appID string, orgID string, userID string, category *model.Ring) (*model.Ring, error) {
	category.ID = uuid.NewString()
	category.OrgID = orgID
	category.AppID = appID
	category.UserID = userID
	category.DateCreated = time.Now().UTC()

	_, err := sa.db.rings.InsertOne(&category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// DeleteRing deletes a user wellness ring
func (sa *Adapter) DeleteRing(appID string, orgID string, userID string, id string) error {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: id}}

	_, err := sa.db.rings.DeleteOne(filter, nil)
	if err != nil {
		log.Printf("error deleting todo entry: %s", err)
		return err
	}

	return nil
}

// CreateRingHistory create a new ring history item
func (sa *Adapter) CreateRingHistory(appID string, orgID string, userID string, ringID string, ringHistory *model.RingHistoryEntry) (*model.Ring, error) {

	ringHistory.ID = uuid.NewString()
	ringHistory.DateCreated = time.Now().UTC()

	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: ringID}}
	update := bson.D{
		primitive.E{Key: "$push", Value: bson.D{
			primitive.E{Key: "history", Value: ringHistory},
		}},
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err := sa.db.rings.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("error updating ring history entry: %s", err)
		return nil, fmt.Errorf("error updating ring history entry: %s", err)
	}

	return sa.GetRing(appID, orgID, userID, ringID)
}

// DeleteRingHistory deletes a ring history item
func (sa *Adapter) DeleteRingHistory(appID string, orgID string, userID string, ringID string, ringHistoryID string) (*model.Ring, error) {
	ring, err := sa.GetRing(appID, orgID, userID, ringID)
	if err != nil {
		log.Printf("error on deleting ring history entry: %s", err)
		return nil, fmt.Errorf("error on deleting ring history entry: %s", err)
	}

	if ring == nil || len(ring.History) < 2 {
		log.Printf("ring contains one or less history items: %s", err)
		return nil, fmt.Errorf("ring contains one or less history items: %s", err)
	}

	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: ringID}}
	update := bson.D{
		primitive.E{Key: "$pull", Value: bson.D{
			primitive.E{Key: "history", Value: primitive.E{Key: "id", Value: ringHistoryID}},
		}},
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err = sa.db.rings.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("error updating ring history entry: %s", err)
		return nil, fmt.Errorf("error updating ring history entry: %s", err)
	}

	return sa.GetRing(appID, orgID, userID, ringID)
}

func (sa *Adapter) abortTransaction(sessionContext mongo.SessionContext) {
	err := sessionContext.AbortTransaction(sessionContext)
	if err != nil {
		log.Printf("error aborting a transaction - %s", err)
	}
}

// NewStorageAdapter creates a new storage adapter instance
func NewStorageAdapter(mongoDBAuth string, mongoDBName string, mongoTimeout string) *Adapter {
	timeout, err := strconv.Atoi(mongoTimeout)
	if err != nil {
		log.Println("Set default timeout - 500")
		timeout = 500
	}
	timeoutMS := time.Millisecond * time.Duration(timeout)

	db := &database{mongoDBAuth: mongoDBAuth, mongoDBName: mongoDBName, mongoTimeout: timeoutMS}
	return &Adapter{db: db}
}

//TransactionContext wraps mongo.SessionContext for use by external packages
type TransactionContext interface {
	mongo.SessionContext
}
