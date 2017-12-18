package commands

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/w32blaster/bot-tfl-next-departure/db"
	"github.com/w32blaster/bot-tfl-next-departure/state"
	"github.com/w32blaster/bot-tfl-next-departure/structs"
	"gopkg.in/telegram-bot-api.v4"
)

const (
	buttonCommandReset = "startFromBeginning"
	// keep short to meet Telegram API limitation for the button date (64 bytes)
	commandUpdate = "u"
	commandSave   = "s"
)

// ProcessCommands acts when user sent to a bot some command, for example "/command arg1 arg2"
func ProcessCommands(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	chatID := message.Chat.ID
	command := extractCommand(message.Command())
	log.Println("This is command " + command)

	switch command {

	case "start":

		resp, _ := sendMsg(bot, chatID, "Yay! Welcome! It will be fun to work with me. Start with typing /help \n\n You can start by pressing this button:")

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				*renderKeyboardButtonActivateQuery(" 🚏 Enter the first station"),
			})

		keyboardMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, resp.MessageID, keyboard)
		bot.Send(keyboardMsg)

	case "mybookmarks":
		renderButtonsWithBookmarks(bot, chatID, message.From.ID)

	case "help":

		help := `This bot supports the following commands:
			 /help - Help
			 /mybookmarks - prints saved trips`
		sendMsg(bot, chatID, html.EscapeString(help))

	default:
		sendMsg(bot, chatID, "Sorry, I don't recognyze such command: "+command+", please call /help to get full list of commands I understand")
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

// ProcessButtonCallback fired when a user click to some button in screen
// We expect to have 4 buttons:
//    1. start from beginning
//    2. show only tube times
//    3. show only bus times
//    4. save bookmarks
func ProcessButtonCallback(bot *tgbotapi.BotAPI, callbackQuery *tgbotapi.CallbackQuery, opts *structs.Opts) {

	// in this switch we decide which exactly button was clicked
	if callbackQuery.Data == buttonCommandReset {

		// the button "start from the beginning" was clicked
		state.ResetStateForUser(callbackQuery.From.ID)

		// let's clean previous messages and button
		editConfig := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    callbackQuery.Message.Chat.ID,
				MessageID: callbackQuery.Message.MessageID,
			},
			Text: "Ok, let's start from the beginning. Start typing a station name again",
		}
		resp, _ := bot.Send(editConfig)

		// and send helping button
		renderButtonThatOpensInlineQuery(bot, callbackQuery.Message.Chat.ID, resp.MessageID)

	} else {

		// we assume the JSON is encoded to the button's data
		journeyRequest := fromJSON(callbackQuery.Data)

		if journeyRequest.Command == commandUpdate {

			// now update already printed timetables with new results
			markdownText, _ := GetTimesBetweenStations(journeyRequest.StationIDFrom, journeyRequest.StationIDTo, journeyRequest.Mode, opts)

			// let's update previous message with new timetables
			editConfig := tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    callbackQuery.Message.Chat.ID,
					MessageID: callbackQuery.Message.MessageID,
				},
				Text:      (markdownText + "\n_updated_"),
				ParseMode: "markdown",
			}
			resp, _ := bot.Send(editConfig)

			// and update keyboard
			keyboardMsg := renderKeyboard(journeyRequest.StationIDFrom, journeyRequest.StationIDTo, callbackQuery.Message.Chat.ID, resp.MessageID)
			bot.Send(keyboardMsg)

		} else {

			// save bookmark to the database
			textMarkdown := fmt.Sprintf("Ok, can you send me the name of your bookmark, please? For example, “way home” or “tain to work” (max " + string(db.MaxLengthBookmarkName) + " symbols)")
			sendMsg(bot, callbackQuery.Message.Chat.ID, textMarkdown)

			// save the state that next time user types some request we know what he/she wants to save
			db.SaveStateForBookmark(callbackQuery.From.ID, journeyRequest)
		}

	}

	// notify the telegram that we processed the button, it will turn "loading indicator" off
	bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: callbackQuery.ID,
	})
}

// ProcessSimpleText is called when a user typed simple text to the chat
func ProcessSimpleText(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {

	// firstly, check is there any state for this user
	state, err := db.GetStateFor(message.From.ID)
	if err != nil {
		log.Println(err)
		sendMsg(bot, message.Chat.ID, "Whoops, can't get a state for you, sorry")
	}

	if state != nil {
		if state.Command == db.StateBookmark {

			// ahha, user previousely wanted to save a bookmark and now he/she sends us the name of bookmark,
			// so now we need to create new bookmark with the given name
			err := db.SaveBookmark(message.From.ID, message.Text, &state.JourneyRequest)

			if err == nil {
				db.DeleteStateFor(message.From.ID)
				sendMsg(bot, message.Chat.ID, "The bookmark “"+message.Text+"” was saved")
			}

		} else if state.Command == db.StateStationID {

			// hm... user selectes previousely some station and sends us plain name of a station. We expect to get here ID but not a plain text
			resp, _ := sendMsg(bot, message.Chat.ID, "If you send me the name of a station, please don't forget to specify @nextTrainLondonBot before your search. Or use this button below")

			// and print helping button
			renderButtonThatOpensInlineQuery(bot, message.Chat.ID, resp.MessageID)
		}
	} else {

		// user enters some text and we don't have any state yet.
		// So we just return some helping message
		resp, _ := sendMsg(bot, message.Chat.ID, "If you want to begin searching a station, then type my name @nextTrainLondonBot and type station name. Or use this button below")

		// and print helping button
		renderButtonThatOpensInlineQuery(bot, message.Chat.ID, resp.MessageID)
	}

}

// OnStationSelected function processes the case when user selected some station.
// If this is the first selection (start station), then we save this value and suggest to
// select another one. If this is the second selection (destination), then find times for this journey
func OnStationSelected(bot *tgbotapi.BotAPI, chatID int64, userID int, command string, opts *structs.Opts) {

	previouslySelectedStation := state.GetPreviouslySelectedStation(userID)
	stationID := strings.Split(command, " ")[2]

	isStationTheSame := stationID == previouslySelectedStation // to prevent seatching jounrey from-to the same station
	if len(previouslySelectedStation) == 0 || isStationTheSame {

		if isStationTheSame {
			db.DeleteStateFor(userID)
		}

		// user selects the station the first time.
		// 1. save this station for the future reference
		state.SaveSelectedStationForUser(userID, stationID)

		// 2. send message that user selected station
		textMarkdown := fmt.Sprintf("You selected the station %s", stationID)
		resp, _ := sendMsg(bot, chatID, textMarkdown)

		// 3. send the keyboard layout with one button "start from the beginning"
		keyboard := tgbotapi.NewInlineKeyboardMarkup(

			// row 1
			[]tgbotapi.InlineKeyboardButton{
				*renderKeyboardButtonActivateQuery(" 🚏 Send me the destination station, please"),
			},

			// row 2
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData("↩  Start from the beginning ", buttonCommandReset),
			})

		keyboardMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, resp.MessageID, keyboard)
		bot.Send(keyboardMsg)

	} else {

		// here we appear after a user selects both stations and we ready to show timetables
		markdownText, err := GetTimesBetweenStations(previouslySelectedStation, stationID, "", opts)

		if err != nil {
			sendMsg(bot, chatID, "Ah, sorry, error occurred when I asked TFL for data journey")
		} else {

			// 1. Send result to client
			resp, _ := sendMsg(bot, chatID, markdownText)

			// 2. Print buttons to save the trip and to narrow to one type of transport
			keyboardMsg := renderKeyboard(previouslySelectedStation, stationID, chatID, resp.MessageID)
			bot.Send(keyboardMsg)
		}
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

func renderButtonsWithBookmarks(bot *tgbotapi.BotAPI, chatID int64, userID int) {
	bookmarks := db.GetBookmarksFor(userID)

	var r string
	for _, bookrmark := range *bookmarks {
		r = r + bookrmark.Name + ","
	}
	sendMsg(bot, chatID, "bookmakrs: "+r)
}

// shortcut method to encode object to JSON
func asJSON(stationFrom string, stationTo string, mode string, command string) string {
	journey := &structs.JourneyRequest{
		StationIDFrom: stationFrom,
		StationIDTo:   stationTo,
		Mode:          mode,
		Command:       command,
	}

	bytesJSON, _ := json.Marshal(journey)
	return string(bytesJSON)
}

func fromJSON(rawJSON string) *structs.JourneyRequest {
	var journeyRequest structs.JourneyRequest
	json.Unmarshal([]byte(rawJSON), &journeyRequest)
	return &journeyRequest
}

// simply render a keyboard with buttons that switch the modes
func renderKeyboard(stationIDFrom string, stationID string, chatID int64, messageID int) *tgbotapi.EditMessageReplyMarkupConfig {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(

		// row 1
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🚇  Show only tube ", asJSON(stationIDFrom, stationID, "tube", commandUpdate)),
			tgbotapi.NewInlineKeyboardButtonData("🚌  Show only buses ", asJSON(stationIDFrom, stationID, "bus", commandUpdate)),
		},

		// row 2
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("🔖 Bookmark this search", asJSON(stationIDFrom, stationID, "", commandSave)),
		})

	markup := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, keyboard)
	return &markup
}

func renderKeyboardButtonActivateQuery(message string) *tgbotapi.InlineKeyboardButton {
	emtpyString := ""
	button := tgbotapi.InlineKeyboardButton{
		Text: message,
		SwitchInlineQueryCurrentChat: &emtpyString,
	}
	return &button
}

// renders the button "enter the first station"
func renderButtonThatOpensInlineQuery(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			*renderKeyboardButtonActivateQuery(" 🚏 Enter the first station name"),
		})

	keyboardMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, keyboard)
	bot.Send(keyboardMsg)
}
