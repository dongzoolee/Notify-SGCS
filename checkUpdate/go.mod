module checkUpdate

go 1.15

replace slackApi => ../slackApi
replace updateDB => ../updateDB

require (
	github.com/anaskhan96/soup v1.2.4
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	slackApi v0.0.0-00010101000000-000000000000
	updateDB v0.0.0-00010101000000-000000000000 // indirect
)
