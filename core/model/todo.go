package model

import "time"

// TodoCategory user defined todo category
type TodoCategory struct {
	ID           string     `json:"id" bson:"_id"`
	OrgID        string     `json:"org_id" bson:"org_id"`
	AppID        string     `json:"app_id" bson:"app_id"`
	UserID       string     `json:"user_id" bson:"user_id"`
	Name         string     `json:"name" bson:"name"`
	Color        string     `json:"color" bson:"color"`
	ReminderType string     `json:"reminder_type" bson:"reminder_type"`
	DateCreated  time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated  *time.Time `json:"date_updated" bson:"date_updated"`
} // @name TodoCategory

// TodoEntry user todo entry
type TodoEntry struct {
	ID          string `json:"id" bson:"_id"`
	OrgID       string `json:"org_id" bson:"org_id"`
	AppID       string `json:"app_id" bson:"app_id"`
	UserID      string `json:"user_id" bson:"user_id"`
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Category    *struct {
		ID           string `json:"id" bson:"id"`
		Name         string `json:"name" bson:"name"`
		Color        string `json:"color" bson:"color"`
		ReminderType string `json:"reminder_type" bson:"reminder_type"`
	} `json:"category" bson:"category"`
	WorkDays []string `json:"work_days" bson:"work_days"`
	Location *struct {
		Latitude  float64 `json:"latitude" bson:"latitude"`
		Longitude float64 `json:"longitude" bson:"longitude"`
	} `json:"location" bson:"location"`
	Completed        bool       `json:"completed" bson:"completed"`
	HasDueTime       bool       `json:"has_due_time" bson:"has_due_time"`
	DueDateTime      string     `json:"due_date_time" bson:"due_date_time"`
	ReminderDateTime *time.Time `json:"reminder_date_time" bson:"reminder_date_time"`
	DateCreated      time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated      *time.Time `json:"date_updated" bson:"date_updated"`
} // @name TodoEntry
