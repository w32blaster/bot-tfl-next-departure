package state

import "github.com/w32blaster/bot-tfl-next-departure/db"

// GetPreviouslySelectedStation returns the station ID that user selected before, or
// empty string if user haven't selected anything
func GetPreviouslySelectedStation(user int) string {
	state, _ := db.GetStateFor(user)
	if state == nil {
		return ""
	}
	return state.StationID
}

// SaveSelectedStationForUser save selected station for a given user
func SaveSelectedStationForUser(user int, stationID string) {
	db.SaveStateForStationID(user, stationID)
}

// ResetStateForUser removes any state for given user allowing him/her to start from
// the beginning
func ResetStateForUser(user int) {
	db.DeleteStateFor(user)
}
