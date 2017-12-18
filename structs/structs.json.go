package structs

// Opts command line arguments
type Opts struct {
	Port    int    `short:"p" long:"port" description:"The port for the bot. The default is 8444" default:"8444"`
	Host    string `short:"h" long:"host" description:"The hostname for the bot. Default is localhost" default:"localhost"`
	IsDebug bool   `short:"d" long:"debug" description:"Is it debug? Default is true. Disable it for production."`

	BotToken string `short:"b" long:"bot-token" description:"The Bot-Token. As long as it is the sencitive data, we can't keep it in Github" required:"true"`
	AppID    string `short:"a" long:"appid" description:"AppID that you can find from the TFL website"`
	APIKEY   string `short:"k" long:"apikey" description:"ApiKEY that you can find fom TFL website"`
}

// TFLInlineQueryResult is JSON wrapper for the result returning by TFL
type TFLInlineQueryResult struct {
	Total   int       `json:"total"`
	Query   string    `json:"query"`
	Matches []Station `json:"matches"`
}

// Station is the JSON wrapper for one station, used in TFL response
type Station struct {
	IcsID string   `json:"icsId"`
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Modes []string `json:"modes"`
}

// TFLJourneyQueryResult root response with query results
type TFLJourneyQueryResult struct {
	Journeys []Journey `json:"journeys"`
}

// Journey json wrapper
type Journey struct {
	StartDateTime string `json:"startDateTime"`
	Diration      int    `json:"duration"`
	Legs          []Leg  `json:"legs"`
}

// Leg
type Leg struct {
	DepartureTime  string         `json:"departureTime"`
	IsDisrupted    bool           `json:"isDisrupted"`
	Mode           Mode           `json:"mode"`
	DeparturePoint DeparturePoint `json:"departurePoint"`
	Instruction    Instruction    `json:"instruction"`
}

// Instruction
type Instruction struct {
	Summary  string `json:"Summary"`
	Detailed string `json:"Detailed"`
}

// DeparturePoint
type DeparturePoint struct {
	IcsID        string `json:"icsId"`
	PlatformName string `json:"platformName"`
	StopLetter   string `json:"stopLetter"`
	CommonName   string `json:"commonName"`
	PlaceType    string `json:"placeType"`
}

// Mode
type Mode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// JourneyRequest object that could be encoded to JSON for buttons and bookmarks
// There is limitation for the data sent to buttons - 64 bytes, so we try to keep JSON as short as possible
type JourneyRequest struct {
	StationIDFrom string `json:"f"`
	StationIDTo   string `json:"t"`
	Mode          string `json:"m"`
	Command       string `json:"c"` // update view or save bookmarks
}