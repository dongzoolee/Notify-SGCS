module main

go 1.15

replace updateDB => ./updateDB

replace slackApi => ./slackApi

replace checkUpdate => ./checkUpdate

require (
	checkUpdate v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.3.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/slack-go/slack v0.9.4
	updateDB v0.0.0-00010101000000-000000000000
)
