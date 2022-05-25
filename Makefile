.PHONY: linux
linux:
	mkdir -p bin/ &&env GOOS=linux GOARCH=amd64 go build  -o ./bin/server ./ChatServer/main.go
	mkdir -p bin/ &&env GOOS=linux GOARCH=amd64 go build  -o ./bin/client ./ChatClient/main.go

.PHONY: win
win:
	mkdir -p bin/ &&env GOOS=windows GOARCH=amd64 go build  -o ./bin/server.exe ./ChatServer/main.go
	mkdir -p bin/ &&env GOOS=windows GOARCH=amd64 go build  -o ./bin/client.exe ./ChatClient/main.go