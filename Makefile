gen:
	sqlc generate

wipedb:
	if [[ -f ~/.yata/app.db ]]; then rm ~/.yata/app.db; fi

run:
	go run main.go

logs:
	tail -f ./yata.log
