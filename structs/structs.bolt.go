package structs

// Bookmark used to be saved into the BoltDB
type Bookmark struct {
	Name    string         `json:"n"`
	Journey JourneyRequest `json:"j"`
}

// State keeps the prevous step from a user. That's why we can't detect what this user selected on his previous step
type State struct {
	Command        string         `json:"c"`
	StationID      string         `json:"s"` // if user selected station, then which one?
	JourneyRequest JourneyRequest `json:"j"` // if user wants to save a bookmark, then which one?
}
