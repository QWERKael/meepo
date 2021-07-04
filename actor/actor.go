package actor

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

var System *actor.ActorSystem

func NewRootContext(host string, port int) *actor.RootContext {
	system := actor.NewActorSystem()
	System = system
	cfg := remote.Configure(host, port)
	remoter := remote.NewRemote(system, cfg)
	remoter.Start()
	rootContext := system.Root
	return rootContext
}

