
cd src

go get

## build for windows
go build -o /c/dev/tools/mh/mm.exe ./main.go

## build for linux
# env GOOS=linux go build -o ../bin/linux/mh ./main.go
