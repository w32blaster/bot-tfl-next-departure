# bot-tfl-nextdeparture
Telegram bot that shows the nearest transport departure from your favourite station

Search station by name:

https://api.tfl.gov.uk/swagger/ui/index.html#!/StopPoint/StopPoint_Search

https://api.tfl.gov.uk/StopPoint/Search?query=Mill%20Hill%20East

Examples:
MIll Hill East: 940GZZLUMHL / 1000147 (icsId)
Finchley central: 940GZZLUFYC / 1015495 (icsId) / 

----------


Find all the routes from this station:

https://api.tfl.gov.uk/StopPoint/ServiceTypes?id=940GZZLUMHL
https://api.tfl.gov.uk/StopPoint/940GZZLUMHL/Route

https://api.tfl.gov.uk/Line/Route?ids=northern&serviceTypes=Regular

-----


find times by stations:

https://api.tfl.gov.uk/Journey/JourneyResults/1000147/to/1015495?date=20171129&time=1553&timeIs=Departing&mode=tube&accessibilityPreference=NoRequirements&app_id=4c754c2a&app_key=9eec9fd4bb56bf3732b2627b391d05b9


# Development
mkdir storage