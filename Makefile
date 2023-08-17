tidy:
	go mod tidy

client: tidy
	go build client/

server: tidy
	go build server/



