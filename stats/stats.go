package stats

import (
	"strings"

	"github.com/botanio/sdk/go"
	"github.com/w32blaster/bot-tfl-next-departure/structs"
	"gopkg.in/telegram-bot-api.v4"
)

type BotanCommand struct {
	Params []string
}

// TrackMessage sends statistics about a message to the Botan.io
func TrackMessage(message *tgbotapi.Message, botanToken string) {

	tokens := strings.Fields(message.Text)
	command := tokens[0]

	// track it
	_track(botanToken, message.From.ID, BotanCommand{tokens}, command)
}

// TrackJourney track a bookmark showing
func TrackJourney(journey *structs.JourneyRequest, userID int, botanToken string) {

	// track it
	_track(botanToken, userID, journey, "ShowJourney")
}

// TrackError track an error
func TrackError(err error, userFrom int, botanToken string) {

	// track it
	_track(botanToken, userFrom, err.Error(), "Error")
}

// common method to send track asynchronousely
func _track(botanToken string, userID int, message interface{}, eventName string) {

	ch := make(chan bool)

	bot := botan.New(botanToken)

	bot.TrackAsync(userID, message, eventName, func(ans botan.Answer, err []error) {
		ch <- true
	})

	<-ch // Synchronization receive

}
