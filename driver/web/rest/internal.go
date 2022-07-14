package rest

import (
	"log"
	"net/http"
	"wellness/core"
)

//InternalApisHandler handles the rest Internal APIs implementation
type InternalApisHandler struct {
	app *core.Application
}

// ProcessReminders Process reminders. Invoked by AWS Scheduler API
// @Description Process reminders. Invoked by AWS Scheduler API
// @Tags Internal
// @ID ProcessReminders
// @Success 200
// @Security InternalAPIAuth
// @Router /int/process_reminders [post]
func (h InternalApisHandler) ProcessReminders(w http.ResponseWriter, r *http.Request) {
	err := h.app.Services.ProcessReminders()
	if err != nil {
		log.Printf("Error on processing reminders - %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
