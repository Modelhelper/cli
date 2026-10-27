build:
	go mod tidy
	
	go build -o ~/dev/tools/bin/mh ./cmd/mh

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/win/mh.exe ./cmd/mh