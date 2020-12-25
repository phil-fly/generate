# Go parameters
# build with version infos

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

BINARY_NAME=generate.exe

RemoteAddr = 10.10.27.11
RemotePort = 8080

ldflags="-s -w -H=windowsgui -w -X generate.RemoteAddr=10.10.27.11 -X generate.RemotePort=8080 -X generate.GenerateMod=autotrace"
#go build -ldflags "-s -w -H=windowsgui -w -X main.RemoteAddr=10.10.27.11 -X main.RemotePort=8080 -X main.GenerateMod=autotrace"

all: clean build
build :
	GOOS=windows GOARCH=386 $(GOBUILD) -ldflags ${ldflags} -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)