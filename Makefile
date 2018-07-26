BIN := nginx-man

$(BIN):
	go build -o $(BIN) cmd/nginx_modsite/main.go
