build_linux:
	go mod tidy
	
	go build -o ~/tools/mh/mh ./main.go
