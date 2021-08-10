
cd src

go get

## build for windows
# go build -o ./bin/win/mh.exe ./main.go

## build for linux
# env GOOS=linux go build -o ../bin/linux/mh ./main.go
env GOOS=linux go build -o ~/dev/tools/mh ./main.go
