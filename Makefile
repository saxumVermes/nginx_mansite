build:
	go build -ldflags "-X main.AvailablePath=${NGINX_A_PATH}" -ldflags "-X main.EnabledPath=${NGINX_E_PATH}" -ldflags "-X main.TemplatePath=$$(pwd)" -a "github.com/saxumVermes/nginx_mansite/cmd/nginx-man"

install: build
	sudo ln -s ${PWD}/nginx-man /usr/local/bin/nginx-man

clean:
	sudo rm /usr/local/bin/nginx-man
	rm ./nginx-man
