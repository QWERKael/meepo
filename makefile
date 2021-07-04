all: main show gossip

main:
	go build -o target/mac/meepo main.go
show:
	go build -buildmode=plugin -o target/mac/plugin/show.so plugin/show/show.go
gossip:
	go build -buildmode=plugin -o target/mac/plugin/gossip.so plugin/gossip/gossip.go

clean:
	rm -f meepo plugin/show.so plugin/gossip.so
