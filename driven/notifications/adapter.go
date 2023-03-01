package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wellness/core/model"

	"github.com/rokwire/core-auth-library-go/v2/authservice"
)

type MessageRef struct {
	OrgID string `json:"org_id" bson:"org_id"`
	AppID string `json:"app_id" bson:"app_id"`
	ID    string `json:"id" bson:"_id"`
}

// Adapter implements the Notifications interface
type Adapter struct {
	internalAPIKey    string
	baseURL           string
	accountManager    *authservice.ServiceAccountManager
	multiTenancyAppID string
	multiTenancyOrgID string
}

// NewNotificationsAdapter creates a new Notifications BB adapter instance
func NewNotificationsAdapter(internalAPIKey string, notificationsHost string, accountManager *authservice.ServiceAccountManager, mtAppID string, mtOrgID string) *Adapter {
	return &Adapter{internalAPIKey: internalAPIKey, baseURL: notificationsHost, accountManager: accountManager, multiTenancyAppID: mtAppID,
		multiTenancyOrgID: mtOrgID}
}

// SendNotification sends notification to a user
func (na *Adapter) SendNotification(recipients []model.NotificationRecipient, topic *string, title string, text string, appID string, orgID string, time *int64, data map[string]string) (*string, error) {
	return na.sendNotification(recipients, topic, title, text, appID, orgID, time, data)
}

func (na *Adapter) sendNotification(recipients []model.NotificationRecipient, topic *string, title string, text string, appID string, orgID string, time *int64, data map[string]string) (*string, error) {
	if len(recipients) > 0 {
		url := fmt.Sprintf("%s/api/bbs/message", na.baseURL)

		async := true
		message := map[string]interface{}{
			"org_id":     orgID,
			"app_id":     appID,
			"priority":   10,
			"recipients": recipients,
			"topic":      nil, //topic,
			"subject":    title,
			"body":       text,
			"data":       data,
			"time":       time,
		}
		bodyData := map[string]interface{}{
			"async":   async,
			"message": message,
		}
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			log.Printf("SendNotification::error creating notification request - %s", err)
			return nil, err
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
		if err != nil {
			log.Printf("SendNotification:error creating load user data request - %s", err)
			return nil, err
		}

		resp, err := na.accountManager.MakeRequest(req, appID, orgID)
		if err != nil {
			log.Printf("SendNotification: error sending request - %s", err)
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Printf("SendNotification: error with response code - %d", resp.StatusCode)
			return nil, fmt.Errorf("SendNotification:error with response code != 200")
		} else {
			var notificationResponse MessageRef
			err := json.NewDecoder(resp.Body).Decode(&notificationResponse)
			if err != nil {
				log.Printf("SendNotification: error with response code - %d", resp.StatusCode)
				return nil, fmt.Errorf("SendNotification: %s", err)
			}
			return &notificationResponse.ID, nil
		}
	}
	return nil, nil
}

// DeleteNotification deletes notification
func (na *Adapter) DeleteNotification(appID string, orgID string, id string) error {

	url := fmt.Sprintf("%s/api/bbs/message/%s", na.baseURL, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("DeleteNotification:error creating load user data request - %s", err)
		return err
	}

	resp, err := na.accountManager.MakeRequest(req, appID, orgID)
	if err != nil {
		log.Printf("DeleteNotification: error sending request - %s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("SendNotification: error with response code - %d", resp.StatusCode)
		return fmt.Errorf("DeleteNotification: error with response code != 200")
	}
	return nil
}
