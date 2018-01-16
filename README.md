![TFL Bot telegram](https://github.com/w32blaster/bot-tfl-next-departure/blob/master/img/tfl-bot-logo.png?raw=true)

# Unofficial TFL Bot for Telegram
Telegram bot that shows the nearest transport departure from your favourite station. You can find the bot itself by the nickname: **@nextTrainLondonBot** or by following the link:

https://t.me/nextTrainLondonBot

## Why this repository exists?
The bot is the demonstration of bot for my presentation "Don't talk to me. Talk to my bot", so it has many comments to help developers get familiar with Telegram bots. Please find the step-by-step tutorial how to create your own bot, here is the link:

📖 [Wiki: how to start with bot development](https://github.com/w32blaster/bot-tfl-next-departure/wiki)

## Save data
You can store your search as "bookmarks". We use BoltDB storage for that.

## TFL API

The bot uses TFL Api. Here are main requests that we perform on the request.

**Search station by name**. We call this endpoint when user searches a station and a popup window sith suggestions appears.

The documentation is [here](https://api.tfl.gov.uk/swagger/ui/index.html#!/StopPoint/StopPoint_Search). For example, find all stations matching, say, "Liverp" we call:

https://api.tfl.gov.uk/StopPoint/Search?query=Liverp

**Find a route between two stations**

For each station we use _icsId_ unique identifier. Here is the example of http query for the route between two stations (full documentation is [here](https://api.tfl.gov.uk/swagger/ui/index.html#!/Journey/Journey_JourneyResults)):

https://api.tfl.gov.uk/Journey/JourneyResults/1000147/to/1015495?date=20171129&time=1553&timeIs=Departing&mode=tube&accessibilityPreference=NoRequirements&app_id=<secret>&app_key=<secret>
