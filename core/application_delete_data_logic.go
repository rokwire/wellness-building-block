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

package core

import (
	"time"
	"wellness/core/model"

	"github.com/rokwire/rokwire-building-block-sdk-go/utils/logging/logs"
)

type deleteDataLogic struct {
	logger *logs.Logger

	storage     Storage
	coreAdapter Core

	//delete data timer
	dailyDeleteTimer *time.Timer
	timerDone        chan bool
}

func (d deleteDataLogic) start() error {

	// set up web tools timer
	go d.setupTimerForDelete()

	return nil
}

func (d deleteDataLogic) setupTimerForDelete() {
	d.logger.Info("Delete data timer")

	//cancel if active
	if d.dailyDeleteTimer != nil {
		d.logger.Info("setupTimerForDelete -> there is active timer, so cancel it")

		d.timerDone <- true
		d.dailyDeleteTimer.Stop()
	}

	//wait until it is the correct moment from the day
	location, err := time.LoadLocation("America/Chicago")
	if err != nil {
		d.logger.Errorf("Error getting location:%s\n", err.Error())
	}
	now := time.Now().In(location)
	d.logger.Infof("setupTimerForDelete -> now - hours:%d minutes:%d seconds:%d\n", now.Hour(), now.Minute(), now.Second())

	nowSecondsInDay := 60*60*now.Hour() + 60*now.Minute() + now.Second()
	desiredMoment := 14400 //4 AM

	var durationInSeconds int
	d.logger.Infof("setupTimerForDelete -> nowSecondsInDay:%d desiredMoment:%d\n", nowSecondsInDay, desiredMoment)
	if nowSecondsInDay <= desiredMoment {
		d.logger.Infof("setupTimerForDelete -> not delete process today, so the first process will be today")
		durationInSeconds = desiredMoment - nowSecondsInDay
	} else {
		d.logger.Infof("setupTimerForDelete -> the delete process has already been processed today, so the first process will be tomorrow")
		leftToday := 86400 - nowSecondsInDay
		durationInSeconds = leftToday + desiredMoment // the time which left today + desired moment from tomorrow
	}
	//log.Println(durationInSeconds)
	//duration := time.Second * time.Duration(3)
	duration := time.Second * time.Duration(durationInSeconds)
	d.logger.Infof("setupTimerForDelete -> first call after %s", duration)

	d.dailyDeleteTimer = time.NewTimer(duration)
	select {
	case <-d.dailyDeleteTimer.C:
		d.logger.Info("setupTimerForDelete -> delete timer expired")
		d.dailyDeleteTimer = nil

		d.process()
	case <-d.timerDone:
		// timer aborted
		d.logger.Info("setupTimerForDelete -> delete timer aborted")
		d.dailyDeleteTimer = nil
	}
}

func (d deleteDataLogic) process() {
	d.logger.Info("Deleting data process")

	//process work
	d.processDelete()

	//generate new processing after 24 hours
	duration := time.Hour * 24
	d.logger.Infof("Deleting data process -> next call after %s", duration)
	d.dailyDeleteTimer = time.NewTimer(duration)
	select {
	case <-d.dailyDeleteTimer.C:
		d.logger.Info("Deleting data process -> timer expired")
		d.dailyDeleteTimer = nil

		d.process()
	case <-d.timerDone:
		// timer aborted
		d.logger.Info("Deleting data process -> timer aborted")
		d.dailyDeleteTimer = nil
	}
}

func (d deleteDataLogic) processDelete() {

	//load deleted accounts
	deletedMemberships, err := d.coreAdapter.LoadDeletedMemberships()

	if err != nil {
		d.logger.Errorf("error on loading deleted accounts - %s", err)
		return
	}

	//process by app org

	for _, appOrgSection := range deletedMemberships {
		d.logger.Infof("delete - [app-id:%s org-id:%s]", appOrgSection.AppID, appOrgSection.OrgID)

		accountsIDs := d.getAccountsIDs(appOrgSection.Memberships)
		if len(accountsIDs) == 0 {
			d.logger.Info("no accounts for deletion")
			continue
		}

		d.logger.Infof("accounts for deletion - %s", accountsIDs)

		//delete the data
		d.deleteAppOrgUsersData(appOrgSection.AppID, appOrgSection.OrgID, accountsIDs)
	}
}

func (d deleteDataLogic) deleteAppOrgUsersData(appID string, orgID string, accountsIDs []string) {

	// delete the todo categories
	err := d.storage.DeleteTodoCategoriesForUsers(appID, orgID, accountsIDs)
	if err != nil {
		d.logger.Errorf("error deleting todo categories for users - %s", err)
		return
	}

	// delete the todo entries
	err = d.storage.DeleteTodoEntriesForUsers(appID, orgID, accountsIDs)
	if err != nil {
		d.logger.Errorf("error deleting todo entries for users - %s", err)
		return
	}

	// delete the rings
	err = d.storage.DeleteRingsForUsers(appID, orgID, accountsIDs)
	if err != nil {
		d.logger.Errorf("error deleting rings for users - %s", err)
		return
	}

	// delete the rings records
	err = d.storage.DeleteRingsForUsers(appID, orgID, accountsIDs)
	if err != nil {
		d.logger.Errorf("error deleting rings records for users - %s", err)
		return
	}
}

func (d deleteDataLogic) getAccountsIDs(memberships []model.DeletedMembership) []string {
	res := make([]string, len(memberships))
	for i, item := range memberships {
		res[i] = item.AccountID
	}
	return res
}

// deleteLogic creates new deleteLogic
func deleteLogic( /*coreAdapter Core,*/ logger *logs.Logger) deleteDataLogic {
	timerDone := make(chan bool)
	return deleteDataLogic{ /*coreAdapter: coreAdapter,*/ timerDone: timerDone, logger: logger}
}
