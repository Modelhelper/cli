build:
	go mod tidy
	
	go build -o ~/dev/tools/bin/mh ./main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/win/mh.exe ./main.go