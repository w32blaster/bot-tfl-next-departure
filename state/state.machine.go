package state

// in the nearest future this will be replaced with proper storage
var previousStation string

// GetPreviouslySelectedStation returns the station ID that user selected before, or
// empty string if user haven't selected anything
func GetPreviouslySelectedStation(user int) string {
	return previousStation
}

// SaveSelectedStationForUser save selected station for a given user
func SaveSelectedStationForUser(user int, stationID string) {
	previousStation = stationID
}
