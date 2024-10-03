build:
	go mod tidy
	
	go build -o ~/tools/mh/mh ./main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/win/mh.exe ./main.go