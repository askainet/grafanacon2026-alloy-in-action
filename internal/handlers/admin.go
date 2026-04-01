package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/grafana/alloy-mission-control/internal/missioncontrol"
)

type AdminHandler struct {
	controller *missioncontrol.Controller
}

func NewAdminHandler(controller *missioncontrol.Controller) *AdminHandler {
	return &AdminHandler{
		controller: controller,
	}
}

// Status returns the current mission status
func (h *AdminHandler) Status(w http.ResponseWriter, r *http.Request) {
	activeMissions := h.controller.GetActiveMissions()
	allMissions := h.controller.Registry().List()

	type missionInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Active      bool   `json:"active"`
	}

	response := struct {
		TotalMissions  int           `json:"total_missions"`
		ActiveMissions int           `json:"active_missions"`
		Missions       []missionInfo `json:"missions"`
	}{
		TotalMissions:  len(allMissions),
		ActiveMissions: len(activeMissions),
		Missions:       make([]missionInfo, 0),
	}

	// Build map of active mission IDs
	activeMissionIDs := make(map[string]bool)
	for _, m := range activeMissions {
		activeMissionIDs[m.ID()] = true
	}

	// Add all missions with active status
	for _, m := range allMissions {
		response.Missions = append(response.Missions, missionInfo{
			ID:          m.ID(),
			Name:        m.Name(),
			Description: m.Description(),
			Active:      activeMissionIDs[m.ID()],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// StartMission activates a mission by ID
func (h *AdminHandler) StartMission(w http.ResponseWriter, r *http.Request) {
	missionID := chi.URLParam(r, "id")

	if err := h.controller.Activate(r.Context(), missionID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "activated",
		"mission": missionID,
	})
}

// Reset deactivates all missions
func (h *AdminHandler) Reset(w http.ResponseWriter, r *http.Request) {
	h.controller.DeactivateAll(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "all missions reset",
	})
}
