package main

import (
	"net/http"
	"strconv"
	"strings"

	"log"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
	"github.com/w32blaster/bot-tfl-next-departure/commands"
	"github.com/w32blaster/bot-tfl-next-departure/db"
	"github.com/w32blaster/bot-tfl-next-departure/stats"
	"github.com/w32blaster/bot-tfl-next-departure/structs"
	"gopkg.in/telegram-bot-api.v4"
)

func main() {

	// get the command line arguments and parse them
	var opts = structs.Opts{}
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(opts.BotToken)
	if err != nil {
		panic("Bot doesn't work. Reason: " + err.Error())
	}

	bot.Debug = opts.IsDebug

	// initiate the database structure
	db.Init()

	log.Printf("Authorized on account %s", bot.Self.UserName)
	updates := bot.ListenForWebhook("/" + bot.Token)

	go http.ListenAndServe(":"+strconv.Itoa(opts.Port), nil)

	for update := range updates {

		if update.Message != nil {

			if update.Message.IsCommand() {

				// This is a command starting with slash
				commands.ProcessCommands(bot, update.Message)

			} else {

				if update.Message.ReplyToMessage == nil {

					if strings.HasPrefix(update.Message.Text, "Selected Station:") {

						// user selected some inline query here
						commands.OnStationSelected(bot, update.Message.Chat.ID, update.Message.From.ID, update.Message.Text, &opts)

					} else {

						// This is a simple text
						log.Println("This is plain text: " + update.Message.Text)
						commands.ProcessSimpleText(bot, update.Message)
					}
				}
			}

			stats.TrackMessage(update.Message, opts.BotanToken)

		} else if update.CallbackQuery != nil {

			// this is the callback after a button click
			commands.ProcessButtonCallback(bot, update.CallbackQuery, &opts)

		} else if update.InlineQuery != nil {

			// this is inline query
			commands.ProcessInlineQuery(bot, update.InlineQuery, &opts)
		}

	}

}
