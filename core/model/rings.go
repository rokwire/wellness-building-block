package model

import "time"

/*

{
	"id":"ldlkjdslkcdslknc",
	"history":[
		{
			"id":"1dddd",
			"color":"",
			"name":"ddds",
			"goal":4.0,
			"date_created":"2022.2..."
		}
	]
}

*/

// Ring represents wellness ring wrapper
type Ring struct {
	ID          string             `json:"id" bson:"_id"`
	AppID       string             `json:"app_id" bson:"app_id"`
	OrgID       string             `json:"org_id" bson:"org_id"`
	UserID      string             `json:"user_id" bson:"user_id"`
	History     []RingHistoryEntry `json:"history" bson:"history"`
	DateCreated time.Time          `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time         `json:"date_updated" bson:"date_updated"`
} // @name Ring

// RingHistoryEntry represents single history entry
type RingHistoryEntry struct {
	ID          string     `json:"id" bson:"id"`
	Color       string     `json:"color_hex" bson:"color_hex"`
	Name        string     `json:"name" bson:"name"`
	Unit        string     `json:"unit" bson:"unit"`
	Value       float64    `json:"value" bson:"value"`
	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
} // @name RingHistoryEntry

// WellnessRingData represent ring data
type WellnessRingData struct {
	ID             string     `json:"id" bson:"_id"`
	AppID          string     `json:"app_id" bson:"app_id"`
	OrgID          string     `json:"org_id" bson:"org_id"`
	UserID         string     `json:"user_id" bson:"user_id"`
	WellnessRingID string     `json:"wellness_ring_id" bson:"wellness_ring_id"`
	Color          string     `json:"color_hex" bson:"color_hex"`
	Name           string     `json:"name" bson:"name"`
	Unit           string     `json:"unit" bson:"unit"`
	Value          float64    `json:"value" bson:"value"`
	DateCreated    time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated    *time.Time `json:"date_updated" bson:"date_updated"`
} //@name WellnessRingData
