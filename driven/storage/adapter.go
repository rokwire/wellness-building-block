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

package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	"wellness/core/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

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

// PerformTransaction performs a transaction
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

		filter := bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "user_id", Value: userID},
			primitive.E{Key: "_id", Value: id},
		}

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
			primitive.E{Key: "category.id", Value: id},
		}
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

// GetTodoEntries gets user's todo entries
func (sa *Adapter) GetTodoEntriesForMigration() ([]model.TodoEntry, error) {
	filter := bson.D{}

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
func (sa *Adapter) CreateTodoEntry(appID string, orgID string, userID string, category *model.TodoEntry, messageIDs model.MessageIDs) (*model.TodoEntry, error) {
	category.ID = uuid.NewString()
	category.OrgID = orgID
	category.AppID = appID
	category.UserID = userID
	category.DateCreated = time.Now().UTC()
	category.MessageIDs = messageIDs

	_, err := sa.db.todoEntries.InsertOne(&category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateTodoEntry updates a todo entry
func (sa *Adapter) UpdateTodoEntry(appID string, orgID string, userID string, todo *model.TodoEntry, id string) (*model.TodoEntry, error) {

	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: id}}
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
			primitive.E{Key: "task_time", Value: todo.TaskTime},
			primitive.E{Key: "message_ids", Value: todo.MessageIDs},
			primitive.E{Key: "date_updated", Value: time.Now().UTC()},
		}},
	}
	_, err := sa.db.todoEntries.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("error updating user defined todo entry: %s", err)
		return nil, err
	}

	return sa.GetTodoEntry(appID, orgID, userID, id)
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
func (sa *Adapter) GetTodoEntriesWithCurrentReminderTime(context TransactionContext, reminderTime time.Time) ([]model.TodoEntry, error) {
	startDate := time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), reminderTime.Hour(), reminderTime.Minute(), 0, 0, reminderTime.Location())
	endDate := time.Date(reminderTime.Year(), reminderTime.Month(), reminderTime.Day(), reminderTime.Hour(), reminderTime.Minute(), 59, 999999, reminderTime.Location())
	filter := bson.D{
		primitive.E{Key: "completed", Value: false},
		primitive.E{Key: "reminder_date_time", Value: []primitive.E{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
		primitive.E{Key: "$or", Value: []primitive.M{
			{"task_time": primitive.M{"$exists": false}},
			{"task_time": nil},
			{"task_time": primitive.M{"$lt": startDate}},
		}},
	}
	update := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{Key: "task_time", Value: reminderTime},
		}},
	}

	var result []model.TodoEntry
	err := sa.db.todoEntries.FindOneAndUpdateWithContext(context, filter, update, &result, &options.FindOneAndUpdateOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTodoEntriesWithCurrentDueTime Gets all todo entries that are applied for the specified due datetime
func (sa *Adapter) GetTodoEntriesWithCurrentDueTime(context TransactionContext, dueTime time.Time) ([]model.TodoEntry, error) {
	startDate := time.Date(dueTime.Year(), dueTime.Month(), dueTime.Day(), dueTime.Hour(), dueTime.Minute(), 0, 0, dueTime.Location())
	endDate := time.Date(dueTime.Year(), dueTime.Month(), dueTime.Day(), dueTime.Hour(), dueTime.Minute(), 59, 999999, dueTime.Location())
	filter := bson.D{
		primitive.E{Key: "completed", Value: false},
		primitive.E{Key: "has_due_time", Value: true},
		primitive.E{Key: "due_date_time", Value: []primitive.E{
			{Key: "$gte", Value: startDate},
			{Key: "$lte", Value: endDate},
		}},
		primitive.E{Key: "$or", Value: []primitive.M{
			{"task_time": primitive.M{"$exists": false}},
			{"task_time": nil},
			{"task_time": primitive.M{"$lt": startDate}},
		}},
	}

	var result []model.TodoEntry
	err := sa.db.todoEntries.Find(filter, &result, &options.FindOptions{Sort: bson.D{{"name", 1}}})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateTodoEntriesTaskTime Updates task time field for the desired todo ids
func (sa *Adapter) UpdateTodoEntriesTaskTime(context TransactionContext, ids []string, taskTime time.Time) error {
	if len(ids) > 0 {
		filter := bson.D{
			bson.E{Key: "_id", Value: bson.M{"$in": ids}},
		}
		update := bson.D{
			bson.E{Key: "$set", Value: bson.D{
				bson.E{Key: "task_time", Value: taskTime},
			}},
		}
		_, err := sa.db.todoEntries.UpdateManyWithContext(context, filter, update, &options.UpdateOptions{})
		return err
	}
	return nil
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
func (sa *Adapter) CreateRing(appID string, orgID string, userID string, ring *model.Ring) (*model.Ring, error) {
	ring.ID = uuid.NewString()
	ring.OrgID = orgID
	ring.AppID = appID
	ring.UserID = userID
	ring.DateCreated = time.Now().UTC()

	for index := range ring.History {
		ring.History[index].RingID = ring.ID
	}

	_, err := sa.db.rings.InsertOne(&ring)
	if err != nil {
		return nil, err
	}
	return ring, nil
}

// DeleteRing deletes a user wellness ring
func (sa *Adapter) DeleteRing(appID string, orgID string, userID string, id string) error {

	err := sa.db.dbClient.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			sa.abortTransaction(sessionContext)
			log.Printf("error starting a transaction - %s", err)
			return err
		}

		filter := bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "user_id", Value: userID},
			primitive.E{Key: "_id", Value: id}}

		_, err = sa.db.rings.DeleteOneWithContext(sessionContext, filter, nil)
		if err != nil {
			sa.abortTransaction(sessionContext)
			log.Printf("error deleting user ring: %s", err)
			return err
		}

		filter = bson.D{
			primitive.E{Key: "app_id", Value: appID},
			primitive.E{Key: "org_id", Value: orgID},
			primitive.E{Key: "user_id", Value: userID},
			primitive.E{Key: "ring_id", Value: id}}

		_, err = sa.db.ringsRecords.DeleteOneWithContext(sessionContext, filter, nil)
		if err != nil {
			sa.abortTransaction(sessionContext)
			log.Printf("error deleting user ring records: %s", err)
			return err
		}

		//commit the transaction
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			sa.abortTransaction(sessionContext)
			fmt.Println(err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("error deleting user ring: %s", err)
		return fmt.Errorf("error deleting user ring: %s", err)
	}

	return nil
}

// CreateRingHistory create a new ring history item
func (sa *Adapter) CreateRingHistory(appID string, orgID string, userID string, ringID string, ringHistory *model.RingHistoryEntry) (*model.Ring, error) {

	ringHistory.ID = uuid.NewString()
	ringHistory.RingID = ringID
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
			primitive.E{Key: "history", Value: primitive.M{"id": ringHistoryID}},
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

// GetRingsRecords Get all ring records for the corresponding ring id
func (sa *Adapter) GetRingsRecords(appID string, orgID string, userID string, ringID *string, startDateEpoch *int64, endDateEpoch *int64, offset *int64, limit *int64, order *string) ([]model.RingRecord, error) {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
	}

	if ringID != nil {
		filter = append(filter, primitive.E{Key: "ring_id", Value: ringID})
	}

	if startDateEpoch != nil {
		seconds := *startDateEpoch / 1000
		timeValue := time.Unix(seconds, 0)
		filter = append(filter, primitive.E{Key: "date_created", Value: bson.D{primitive.E{Key: "$gte", Value: &timeValue}}})
	}
	if endDateEpoch != nil {
		seconds := *endDateEpoch / 1000
		timeValue := time.Unix(seconds, 0)
		filter = append(filter, primitive.E{Key: "date_created", Value: bson.D{primitive.E{Key: "$lte", Value: &timeValue}}})
	}

	findOptions := options.Find()
	if order != nil && *order == "asc" {
		findOptions.SetSort(bson.D{{"date_created", 1}})
	} else {
		findOptions.SetSort(bson.D{{"date_created", -1}})
	}
	if limit != nil {
		findOptions.SetLimit(*limit)
	}
	if offset != nil {
		findOptions.SetSkip(*offset)
	}

	var list []model.RingRecord
	err := sa.db.ringsRecords.Find(filter, &list, findOptions)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// GetRingsRecord gets a single ring record by id
func (sa *Adapter) GetRingsRecord(appID string, orgID string, userID string, id string) (*model.RingRecord, error) {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: id},
	}

	var list []model.RingRecord
	err := sa.db.ringsRecords.Find(filter, &list, nil)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}

// CreateRingsRecord creates a new ring record
func (sa *Adapter) CreateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error) {
	record.ID = uuid.NewString()
	record.OrgID = orgID
	record.AppID = appID
	record.UserID = userID
	record.DateCreated = time.Now().UTC()

	_, err := sa.db.ringsRecords.InsertOne(&record)
	if err != nil {
		return nil, err
	}

	return sa.GetRingsRecord(appID, orgID, userID, record.ID)
}

// UpdateRingsRecord updates a ring record
func (sa *Adapter) UpdateRingsRecord(appID string, orgID string, userID string, record *model.RingRecord) (*model.RingRecord, error) {
	now := time.Now().UTC()
	record.DateUpdated = &now
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
		primitive.E{Key: "_id", Value: record.ID}}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "value", Value: record.Value},
			primitive.E{Key: "date_updated", Value: record},
		}},
	}

	_, err := sa.db.ringsRecords.UpdateOne(filter, update, nil)
	if err != nil {
		log.Printf("error delete a ring record: %s", err)
		return nil, fmt.Errorf("error delete a ring record: %s", err)
	}

	return sa.GetRingsRecord(appID, orgID, userID, record.ID)
}

// DeleteRingsRecords deletes a ring record
func (sa *Adapter) DeleteRingsRecords(appID string, orgID string, userID string, ringID *string, recordID *string) error {
	filter := bson.D{
		primitive.E{Key: "app_id", Value: appID},
		primitive.E{Key: "org_id", Value: orgID},
		primitive.E{Key: "user_id", Value: userID},
	}
	if ringID != nil {
		filter = append(filter, primitive.E{Key: "ring_id", Value: *ringID})
	}
	if recordID != nil {
		filter = append(filter, primitive.E{Key: "_id", Value: *recordID})
	}

	_, err := sa.db.ringsRecords.DeleteMany(filter, nil)
	if err != nil {
		log.Printf("error delete a ring records: %s", err)
		return fmt.Errorf("error delete a ring records: %s", err)
	}

	return nil
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

// TransactionContext wraps mongo.SessionContext for use by external packages
type TransactionContext interface {
	mongo.SessionContext
}
