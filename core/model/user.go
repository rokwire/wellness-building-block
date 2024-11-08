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

// RingResponse represents wellness ring wrapper
type RingResponse struct {
	ID     string `json:"id" bson:"_id"`
	UserID string `json:"user_id" bson:"user_id"`
} // @name RingResponse

// RingRecordResponse represents individual daily records for an individual ring
type RingRecordResponse struct {
	ID     string `json:"id" bson:"_id"`
	UserID string `json:"user_id" bson:"user_id"`
} //@name RingRecordResponse

// TodoCategoryResponse user defined todo category
type TodoCategoryResponse struct {
	ID     string `json:"id" bson:"_id"`
	UserID string `json:"user_id" bson:"user_id"`
} // @name TodoCategoryResponse

// TodoEntryResponse user todo entry
type TodoEntryResponse struct {
	ID     string `json:"id" bson:"_id"`
	UserID string `json:"user_id" bson:"user_id"`
} // @name TodoEntryResponse

// UserDataResponse user todo entry
type UserDataResponse struct {
	Rings          []RingResponse         `json:"rings"`
	RingsRecord    []RingRecordResponse   `json:"rings_records"`
	TodoEntries    []TodoEntryResponse    `json:"todo_entries"`
	TodoCategories []TodoCategoryResponse `json:"todo_categories"`
} // @name UserDataResponse
