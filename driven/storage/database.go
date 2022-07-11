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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client

	todoCategories *collectionWrapper
	todoEntries    *collectionWrapper
	rings          *collectionWrapper
	ringsRecords   *collectionWrapper
}

func (m *database) start() error {

	log.Println("database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(m.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(m.mongoDBName)

	todoCategories := &collectionWrapper{database: m, coll: db.Collection("todo_categories")}
	err = m.applyTodoCategoriesChecks(todoCategories)
	if err != nil {
		return err
	}

	todoEntries := &collectionWrapper{database: m, coll: db.Collection("todo_entries")}
	err = m.applyTodoEntriesChecks(todoEntries)
	if err != nil {
		return err
	}

	rings := &collectionWrapper{database: m, coll: db.Collection("rings")}
	err = m.applyRingsChecks(rings)
	if err != nil {
		return err
	}

	ringsRecords := &collectionWrapper{database: m, coll: db.Collection("rings_records")}
	err = m.applyRingsRecordsChecks(ringsRecords)
	if err != nil {
		return err
	}

	m.todoCategories = todoCategories
	m.todoEntries = todoEntries
	m.rings = rings
	m.ringsRecords = ringsRecords

	//asign the db, db client and the collections
	m.db = db
	m.dbClient = client

	return nil
}

// Event

func (m *database) onDataChanged(changeDoc map[string]interface{}) {
	if changeDoc == nil {
		return
	}
	log.Printf("onDataChanged: %+v\n", changeDoc)
	ns := changeDoc["ns"]
	if ns == nil {
		return
	}
	nsMap := ns.(map[string]interface{})
	coll := nsMap["coll"]

	if "configs" == coll {
		log.Println("configs collection changed")
	} else {
		log.Println("other collection changed")
	}
}

func (m *database) applyTodoCategoriesChecks(categories *collectionWrapper) error {
	log.Println("apply todo_categories checks.....")

	//Add org_id + app_id index
	err := categories.AddIndex(
		bson.D{
			primitive.E{Key: "org_id", Value: 1},
			primitive.E{Key: "app_id", Value: 1},
			primitive.E{Key: "user_id", Value: 1},
		},
		false)
	if err != nil {
		return err
	}

	//Add user_id index
	err = categories.AddIndex(
		bson.D{primitive.E{Key: "user_id", Value: 1}},
		false)
	if err != nil {
		return err
	}

	log.Println("todo_categories passed")
	return nil
}

func (m *database) applyTodoEntriesChecks(entries *collectionWrapper) error {
	log.Println("apply todo_entries checks.....")

	err := entries.AddIndex(
		bson.D{
			primitive.E{Key: "org_id", Value: 1},
			primitive.E{Key: "app_id", Value: 1},
			primitive.E{Key: "user_id", Value: 1},
			primitive.E{Key: "category.id", Value: 1},
		},

		false)
	if err != nil {
		return err
	}

	err = entries.AddIndex(
		bson.D{
			primitive.E{Key: "due_date_time", Value: -1},
			primitive.E{Key: "org_id", Value: 1},
			primitive.E{Key: "app_id", Value: 1},
			primitive.E{Key: "user_id", Value: 1},
			primitive.E{Key: "category.id", Value: 1},
		},

		false)
	if err != nil {
		return err
	}

	err = entries.AddIndex(
		bson.D{
			primitive.E{Key: "reminder_date_time", Value: -1},
			primitive.E{Key: "org_id", Value: 1},
			primitive.E{Key: "app_id", Value: 1},
			primitive.E{Key: "user_id", Value: 1},
			primitive.E{Key: "category.id", Value: 1},
		},

		false)
	if err != nil {
		return err
	}

	//Add category.id index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "category.id", Value: 1}},
		false)
	if err != nil {
		return err
	}

	//Add user_id index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "user_id", Value: 1}},
		false)
	if err != nil {
		return err
	}

	//Add due_date_time index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "due_date_time", Value: 1}},
		false)
	if err != nil {
		return err
	}

	//Add reminder_date_time index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "reminder_date_time", Value: 1}},
		false)
	if err != nil {
		return err
	}

	log.Println("todo_entries passed")
	return nil
}

func (m *database) applyRingsChecks(entries *collectionWrapper) error {
	log.Println("apply rings checks.....")

	//Add org_id + app_id index
	err := entries.AddIndex(
		bson.D{
			primitive.E{Key: "org_id", Value: 1},
			primitive.E{Key: "app_id", Value: 1},
			primitive.E{Key: "user_id", Value: 1},
		},
		false)
	if err != nil {
		return err
	}

	//Add user_id index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "user_id", Value: 1}},
		false)
	if err != nil {
		return err
	}

	//Add history.id index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "history.id", Value: 1}},
		false)
	if err != nil {
		return err
	}

	log.Println("rings passed")
	return nil
}

func (m *database) applyRingsRecordsChecks(entries *collectionWrapper) error {
	log.Println("apply rings_records checks.....")

	//Add org_id + app_id index
	err := entries.AddIndex(
		bson.D{
			primitive.E{Key: "org_id", Value: 1},
			primitive.E{Key: "app_id", Value: 1},
			primitive.E{Key: "user_id", Value: 1},
			primitive.E{Key: "ring_id", Value: 1},
		},
		false)
	if err != nil {
		return err
	}

	err = entries.AddIndex(
		bson.D{
			primitive.E{Key: "date_created", Value: 1},
			primitive.E{Key: "org_id", Value: 1},
			primitive.E{Key: "app_id", Value: 1},
			primitive.E{Key: "user_id", Value: 1},
			primitive.E{Key: "ring_id", Value: 1},
		},
		false)
	if err != nil {
		return err
	}

	//Add user_id index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "user_id", Value: 1}},
		false)
	if err != nil {
		return err
	}

	//Add ring_id index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "ring_id", Value: 1}},
		false)
	if err != nil {
		return err
	}

	//Add date_created index
	err = entries.AddIndex(
		bson.D{primitive.E{Key: "date_created", Value: -1}},
		false)
	if err != nil {
		return err
	}

	log.Println("rings_records passed")
	return nil
}
