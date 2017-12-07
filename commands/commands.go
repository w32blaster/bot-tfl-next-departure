package commands

import (
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/w32blaster/bot-tfl-next-departure/state"
	"github.com/w32blaster/bot-tfl-next-departure/structs"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	buttonCommandReset = "startFromBeginning"
)

// ProcessCommands acts when user sent to a bot some command, for example "/command arg1 arg2"
func ProcessCommands(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	chatID := message.Chat.ID
	command := extractCommand(message.Command())
	log.Println("This is command " + command)

	switch command {

	case "start":
		sendMsg(bot, chatID, "Yay! Welcome! It will be fun to work with me. Start with typing /help")

	case "help":

		help := `This bot supports the following commands:
			 /help - Help
			 /mybookmarks - prints saved trips`
		sendMsg(bot, chatID, html.EscapeString(help))

	default:
		if strings.HasPrefix(command, "/station") {
			OnStationSelected(bot, chatID, message.From.ID, command)
		} else {
			sendMsg(bot, chatID, "Sorry, I don't recognyze such command: "+command+", please call /help to get full list of commands I understand")
		}
	}

}

// ProcessInlineQuery returns the array of suggested stations for the Inline Query
func ProcessInlineQuery(bot *tgbotapi.BotAPI, inlineQuery *tgbotapi.InlineQuery, opts *structs.Opts) error {

	// firstly, make query to TFL
	searchQuery := inlineQuery.Query
	foundStations := GetStationListByPattern(searchQuery, opts)

	answer := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		CacheTime:     3,
		Results:       foundStations,
	}

	if resp, err := bot.AnswerInlineQuery(answer); err != nil {
		log.Fatal("ERROR! bot.answerInlineQuery:", err, resp)
		return err
	}

	return nil
}

// OnStationSelected function processes the case when user selected some station.
// If this is the first selection (start station), then we save this value and suggest to
// select another one. If this is the second selection (destination), then find times for this journey
func OnStationSelected(bot *tgbotapi.BotAPI, chatID int64, userID int, command string) {

	previouslySelectedStation := state.GetPreviouslySelectedStation(userID)
	stationID := strings.Split(command, " ")[1]

	if len(previouslySelectedStation) == 0 {

		// user selects the station the first time.
		// 1. save this station for the future reference
		state.SaveSelectedStationForUser(userID, stationID)

		// 2. send message that user selected station
		textMarkdown := fmt.Sprintf("You selected the station %s. Now send me the destination station name please", stationID)
		resp, _ := sendMsg(bot, chatID, textMarkdown)

		// 3. send the keyboard layout with one button "start from the beginning"
		keyboard := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("❌  Start from the beginning ", buttonCommandReset),
		})

		keyboardMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, resp.MessageID, keyboard)
		bot.Send(keyboardMsg)

	} else {

		// user selected two stations, print the journey data
		sendMsg(bot, chatID, "here will be journey")
	}

}

// properly extracts command from the input string, removing all unnecessary parts
// please refer to unit tests for details
func extractCommand(rawCommand string) string {

	command := rawCommand

	// remove slash if necessary
	if rawCommand[0] == '/' {
		command = command[1:]
	}

	// if command contains the name of our bot, remote it
	command = strings.Split(command, "@")[0]
	command = strings.Split(command, " ")[0]

	return command
}

// simply send a message to bot in Markdown format
func sendMsg(bot *tgbotapi.BotAPI, chatID int64, textMarkdown string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, textMarkdown)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	// send the message
	resp, err := bot.Send(msg)
	if err != nil {
		log.Println("bot.Send:", err, resp)
		return resp, err
	}

	return resp, err
}
