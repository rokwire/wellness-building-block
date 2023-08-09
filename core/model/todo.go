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

package model

import "time"

// TodoCategory user defined todo category
type TodoCategory struct {
	ID          string     `json:"id" bson:"_id"`
	OrgID       string     `json:"org_id" bson:"org_id"`
	AppID       string     `json:"app_id" bson:"app_id"`
	UserID      string     `json:"user_id" bson:"user_id"`
	Name        string     `json:"name" bson:"name"`
	Color       string     `json:"color" bson:"color"`
	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
} // @name TodoCategory

// ToCategoryRef Converts to CategoryRef
func (c *TodoCategory) ToCategoryRef() CategoryRef {
	return CategoryRef{
		ID: c.ID, OrgID: c.OrgID, AppID: c.AppID, UserID: c.UserID,
		Name: c.Name, Color: c.Color,
	}
}

// TodoEntry user todo entry
type TodoEntry struct {
	ID               string       `json:"id" bson:"_id"`
	OrgID            string       `json:"org_id" bson:"org_id"`
	AppID            string       `json:"app_id" bson:"app_id"`
	UserID           string       `json:"user_id" bson:"user_id"`
	Title            string       `json:"title" bson:"title"`
	Description      string       `json:"description" bson:"description"`
	Category         *CategoryRef `json:"category" bson:"category"`
	WorkDays         []string     `json:"work_days" bson:"work_days"`
	Location         *string      `json:"location" bson:"location"`
	Completed        bool         `json:"completed" bson:"completed"`
	HasDueTime       bool         `json:"has_due_time" bson:"has_due_time"`
	DueDateTime      *time.Time   `json:"due_date_time" bson:"due_date_time"`
	ReminderType     string       `json:"reminder_type" bson:"reminder_type"`
	ReminderDateTime *time.Time   `json:"reminder_date_time" bson:"reminder_date_time"`
	MessageIDs       MessageIDs   `json:"message_ids" bson:"message_ids"`
	TaskTime         *time.Time   `json:"task_time" bson:"task_time"`
	DateCreated      time.Time    `json:"date_created" bson:"date_created"`
	DateUpdated      *time.Time   `json:"date_updated" bson:"date_updated"`
	RecurrenceType   *string      `json:"recurrence_type" bson:"recurrence_type"`
	RecurrenceID     *string      `json:"recurrence_id" bson:"recurrence_id"`
} // @name TodoEntry

// RequiresMessageIDsMigration Checks if the record requires db data migration
func (t *TodoEntry) RequiresMessageIDsMigration() bool {
	return (t.DueDateTime != nil && time.Now().Before(*t.DueDateTime) && t.MessageIDs.DueDateMessageID == nil) ||
		(t.ReminderDateTime != nil && time.Now().Before(*t.ReminderDateTime) && t.MessageIDs.ReminderDateMessageID == nil)
}

// MessageIDs is used to collect due and reminder time messages
type MessageIDs struct {
	ReminderDateMessageID *string `json:"reminder_date_message_id" bson:"reminder_date_message_id"`
	DueDateMessageID      *string `json:"due_date_message_id" bson:"due_date_message_id"`
}

// CategoryRef used as a reference within the TodoEntry
type CategoryRef struct {
	ID     string `json:"id" bson:"id"`
	OrgID  string `json:"org_id" bson:"org_id"`
	AppID  string `json:"app_id" bson:"app_id"`
	UserID string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Color  string `json:"color" bson:"color"`
} // @name CategoryRef
