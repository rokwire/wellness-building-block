package model

import "time"

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
	RingID      string     `json:"ring_id" bson:"ring_id"`
	Color       string     `json:"color_hex" bson:"color_hex"`
	Name        string     `json:"name" bson:"name"`
	Unit        string     `json:"unit" bson:"unit"`
	Value       float64    `json:"value" bson:"value"`
	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
} // @name RingHistoryEntry

// RingRecord represents individual daily records for an individual ring
type RingRecord struct {
	ID          string     `json:"id" bson:"_id"`
	AppID       string     `json:"app_id" bson:"app_id"`
	OrgID       string     `json:"org_id" bson:"org_id"`
	UserID      string     `json:"user_id" bson:"user_id"`
	RingID      string     `json:"ring_id" bson:"ring_id"`
	Value       float64    `json:"value" bson:"value"`
	DateCreated time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated *time.Time `json:"date_updated" bson:"date_updated"`
} //@name RingRecord
