package main

import (
	"net/http"
	"strconv"

	"log"

	flags "github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/telegram-bot-api.v4"
)

// command line arguments
var opts struct {
	Port    int    `short:"p" long:"port" description:"The port for the bot. The default is 8444" default:"8444"`
	Host    string `short:"h" long:"host" description:"The hostname for the bot. Default is localhost" default:"localhost"`
	IsDebug bool   `short:"d" long:"debug" description:"Is it debug? Default is true. Disable it for production." default:"True"`

	BotToken string `short:"b" long:"bot-token" description:"The Bot-Token. As long as it is the sencitive data, we can't keep it in Github" required:"true"`
}

func main() {

	// get the command line arguments and parse them
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(opts.BotToken)
	if err != nil {
		panic("Bot doesn't work")
	}

	bot.Debug = true
	updates := bot.ListenForWebhook("/bot/" + bot.Token)

	go http.ListenAndServe(":"+strconv.Itoa(opts.Port), nil)

	for update := range updates {

		if update.Message != nil {

			if update.Message.IsCommand() {

				// This is a command
				log.Println("This is a command starting with '/' ")

			} else {

				if update.Message.ReplyToMessage == nil {

					// This is a simple text
					log.Println("This is plain text")
				}
			}

		} else if update.CallbackQuery != nil {

			// this is the callback after a button click
			log.Println("this is the callback that will be fired when a user clicks some button generated by bot")

		} else if update.InlineQuery != nil {

			// this is inline query
			log.Println("This is inline query, that shows some ")
		}

	}

}
