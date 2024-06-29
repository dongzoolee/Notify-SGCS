module slackApi

go 1.15

replace updateDB => ../updateDB

require (
	github.com/joho/godotenv v1.3.0
	github.com/slack-go/slack v0.9.4
	updateDB v0.0.0-00010101000000-000000000000
)
