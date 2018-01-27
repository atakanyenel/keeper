
STRIP= -ldflags="-s -w"
mac:
	go build $(STRIP) -o dist/keeper_mac

linux:
	env GOOS=linux GOARCH=amd64 go build $(STRIP) -o dist/keeper_linux
raspberry:
	env GOOS=linux GOARCH=arm GOARM=7 go build $(STRIP) -o dist/keeper_rasp
clean:
	-rm -f dist/keeper
	-rm -f dist/keeper_linux
	-rm -f dist/keeper_mac
	-rm -f dist/keeper_rasp
	-rm -f dist/.kept
	-rm -f dist/.env

all: clean mac linux raspberry

distribute:
	scp keeper_linux root@141.40.254.141:/root
	scp keeper_linux root@141.40.254.2:/root
test:
	go test -cover