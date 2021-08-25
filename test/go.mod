module test

replace updateDB => ../updateDB

replace checkUpdate => ../checkUpdate

replace slackApi => ../slackApi

go 1.16

require (
	checkUpdate v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/slack-go/slack v0.9.4
	updateDB v0.0.0-00010101000000-000000000000
)
