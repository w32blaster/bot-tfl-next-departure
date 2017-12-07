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