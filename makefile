build_linux:
	go get


	env GOOS=linux go build -o ./bin/linux/mh ./main.go
## build for linux

build_win:
	## build for windows
	go get
	
	go build -o ./bin/win/mh.exe ./main.go
