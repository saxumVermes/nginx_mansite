build:
	go build -a github.com/saxumVermes/nginx_mansite/cmd/nginx-man

install: build
	sudo ln -s ${PWD}/nginx-man /usr/local/bin/nginx-man

