package appse

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-sesam/internal/se"
	"github.com/dumacp/smartcard"
)

// type reader struct {
// 	rootctx *actor.RootContext
// 	dev     *actor.PID
// }

func NewActor(dev smartcard.ICard) actor.Actor {
	// r := &reader{rootctx: rootctx}
	app := se.ActorSAM(dev)
	// r.dev = pid
	return app
}

// func (r *reader) PID() *actor.PID {
// 	return r.dev
// }

// func (r *reader) RootContext() *actor.RootContext {
// 	return r.rootctx
// }
