all:	
	go build src/webserver.go src/admin.go src/generate_rss.go src/setup.go

windows: src/admin.go src/webserver.go src/generate_rss.go
	go build -o pogoapp.exe src/webserver.go src/admin.go src/generate_rss.go src/setup.go

linux: src/admin.go src/webserver.go src/generate_rss.go
	go build -o pogoapp src/webserver.go src/admin.go src/generate_rss.go src/setup.go

install:
	go get github.com/gmemstr/feeds
	go get github.com/fsnotify/fsnotify
	go get github.com/spf13/viper
	go get github.com/gorilla/mux

docker:
	docker build .

and run:
	go build src/webserver.go src/admin.go src/generate_rss.go src/setup.go
	./pogoapp.exe