BIN := nginx-man

$(BIN):
	go build -o $(BIN) cmd/nginx_modsite/main.go

into_path: $(BIN)
	sudo mv ./$(BIN) /usr/local/bin/
