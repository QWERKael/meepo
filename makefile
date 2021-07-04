ifeq ($(OS),Windows_NT)
    env_os = windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        env_os = linux
    endif
    ifeq ($(UNAME_S),Darwin)
        env_os = mac
    endif
endif

all: main show gossip

main:
	go build -o target/$(env_os)/meepo main.go
show:
	go build -buildmode=plugin -o target/$(env_os)/plugin/show.so plugin/show/show.go
gossip:
	go build -buildmode=plugin -o target/$(env_os)/plugin/gossip.so plugin/gossip/gossip.go

.PHONY: clean
clean:
	rm -f meepo plugin/show.so plugin/gossip.so
