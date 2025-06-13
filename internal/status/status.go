package status

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"

	"github.com/coreos/airlock/internal/config"
	"github.com/coreos/airlock/internal/lock"
)

const (
	StatusEndpoint = "/status"
)

func Status(settings *config.Settings) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if settings == nil {
			http.Error(w, "nil settings", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json_response := make(map[string]interface{})
		for groupName, group := range settings.LockGroups {
			lockManager, err := lock.NewManager(r.Context(), settings, groupName, group)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"group": groupName,
				}).WithError(err).Error("error creating lock manager")
				http.Error(w, fmt.Sprintf("error creating lock manager: %v", err), http.StatusInternalServerError)
				return
			}
			semaphore, err := lockManager.FetchSemaphore(r.Context())
			json_response[groupName] = semaphore
		}

		err := json.NewEncoder(w).Encode(json_response)
		if err != nil {
			logrus.WithError(err).Error("error encoding response")
			http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
