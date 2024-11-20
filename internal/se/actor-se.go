package se

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dumacp/go-logs/pkg/logs"
	"github.com/dumacp/go-sesam/internal/sam"
	"github.com/dumacp/go-sesam/pkg/messages"
	"github.com/dumacp/go-sesam/pkg/se"
	"github.com/dumacp/smartcard"
	"github.com/looplab/fsm"
)

type samActor struct {
	// pid      *actor.PID
	// rootctx  *actor.RootContext
	fm       *fsm.FSM
	behavior actor.Behavior
	card     smartcard.ICard
	uid      string
	sam      se.SE
	ctx      actor.Context
	contxt   context.Context
	cancel   func()
}

// type Actor interface {
// 	PID() *actor.PID
// 	RootContext() *actor.RootContext
// }

func ActorSAM(r smartcard.ICard) actor.Actor {
	app := &samActor{}
	app.card = r
	// app.initFSM()
	// props := actor.PropsFromFunc(app.Receive)

	// if ctx == nil {
	// 	ctx = actor.NewActorSystem().Root
	// }
	// app.rootctx = ctx
	// uid, err := uuid.NewRandom()
	// if err != nil {
	// 	return nil, err
	// }
	// pid, err := ctx.SpawnNamed(props, fmt.Sprintf("sam-actor-%s", uid.String()))
	// if err != nil {
	// 	return nil, err
	// }
	// app.pid = pid
	app.behavior = make(actor.Behavior, 0)
	app.behavior.Become(app.CloseState)
	app.behavior = actor.NewBehavior()
	app.behavior.Become(app.CloseState)

	return app
}

// func (app *samActor) PID() *actor.PID {
// 	return app.pid
// }
// func (app *samActor) RootContext() *actor.RootContext {
// 	return app.rootctx
// }

func (a *samActor) Receive(ctx actor.Context) {
	// logs.LogBuild.Printf("Message arrived in appActor: %s, %T", ctx.Message(), ctx.Message())
	a.ctx = ctx
	a.behavior.Receive(ctx)
}

func (a *samActor) CloseState(ctx actor.Context) {
	logs.LogBuild.Printf("Message arrived in samActor, behavior (CloseState): %s, %T. %s", ctx.Message(), ctx.Message(), ctx.Sender())
	switch ctx.Message().(type) {
	case *actor.Started:
		contxt, cancel := context.WithCancel(context.Background())
		a.contxt = contxt
		a.cancel = cancel
		a.initFSM()
		ctx.Send(ctx.Self(), &messages.MsgOpen{})
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
	case *messages.MsgOpen:
		if err := func() error {
			c, err := sam.NewSamAV2(a.card)
			if err != nil {
				return err
			}
			for range []int{0, 1, 2} {
				if err = c.Connect(); err != nil {
					continue
				}
				break
			}
			if err != nil {
				return err
			}
			a.sam = c
			a.fm.Event(a.contxt, eOpenCmd)
			logs.LogInfo.Printf("sam UID: [% X]", c.Serial())
			a.uid = hex.EncodeToString(c.Serial())
			return nil
		}(); err != nil {
			logs.LogError.Printf("sam err: %s", err)
			a.fm.Event(a.contxt, eError, err)
		}
	}
}

func (a *samActor) WaitState(ctx actor.Context) {
	logs.LogBuild.Printf("Message arrived in samActor (%s), behavior (WaitState): %+v, %T, %s",
		ctx.Self().GetId(), ctx.Message(), ctx.Message(), ctx.Sender())
	switch msg := ctx.Message().(type) {
	case *actor.Stopping:
		if a.cancel != nil {
			a.cancel()
		}
	case *messages.MsgClose:
		if err := func() error {
			if err := a.sam.Disconnect(); err != nil {
				return err
			}
			a.fm.Event(a.contxt, eClosed)
			return nil
		}(); err != nil {
			logs.LogError.Println(err)
		}
	case *messages.MsgAuth:
		if err := func() error {
			if len(msg.Key) <= 0 {
				return fmt.Errorf("auth sam error: len key is %d", len(msg.Key))
			}
			keyAuth, err := hex.DecodeString(msg.Key)
			if err != nil {
				return fmt.Errorf("auth sam error: %w", err)
			}
			if err := a.sam.Auth(keyAuth, msg.Slot, msg.Version); err != nil {
				return fmt.Errorf("auth sam error: %w", err)
			}
			return nil
		}(); err != nil {
			fmt.Println(err)
			if ctx.Sender != nil {
				ctx.Respond(&messages.MsgAck{
					Error: err.Error(),
				})
			}
		}
	case *messages.MsgApdu:
		if err := func() error {
			resp, err := a.sam.Apdu(msg.Data)
			if err != nil {
				return err
			}
			if ctx.Sender() != nil {
				ctx.Respond(&messages.MsgApduResponse{
					Data: resp,
				})
			}
			return nil
		}(); err != nil {
			logs.LogError.Println(err)
			if ctx.Sender() != nil {
				// time.Sleep(3 * time.Second)
				ctx.Respond(&messages.MsgAck{Error: err.Error()})
			}
			a.fm.Event(a.contxt, eError, err)
		}
	case *messages.MsgEncryptRequest:
		if err := func() error {
			cipher, err := a.sam.Encrypt(msg.Data, msg.IV, msg.DevInput, msg.KeySlot)
			if err != nil {
				return err
			}
			if ctx.Sender() != nil {
				ctx.Respond(&messages.MsgEncryptResponse{
					Cipher: cipher,
					SamUID: a.uid,
					MsgID:  msg.MsgID,
				})
			}
			return nil
		}(); err != nil {
			logs.LogError.Println(err)
			if ctx.Sender() != nil {
				// time.Sleep(3 * time.Second)
				ctx.Respond(&messages.MsgAck{Error: err.Error()})
			}
			a.fm.Event(a.contxt, eError, err)
		}
	case *messages.MsgDecryptRequest:
		if err := func() error {
			cipher, err := a.sam.Decrypt(msg.Data, msg.IV, msg.DevInput, msg.KeySlot)
			if err != nil {
				return err
			}
			if ctx.Sender() != nil {
				ctx.Respond(&messages.MsgDecryptResponse{
					MsgID:  msg.MsgID,
					Plain:  cipher,
					SamUID: a.uid,
				})
			}
			return nil
		}(); err != nil {
			logs.LogError.Println(err)
			if ctx.Sender() != nil {
				// time.Sleep(3 * time.Second)
				ctx.Respond(&messages.MsgAck{Error: err.Error()})
			}
			a.fm.Event(a.contxt, eError, err)
		}
	case *messages.MsgDumpSecretKeyRequest:
		if err := func() error {
			key, err := a.sam.DumpSecretKey(msg.KeySlot)
			if err != nil {
				return err
			}
			if ctx.Sender() != nil {
				ctx.Respond(&messages.MsgDumpSecretKeyResponse{
					Data: key,
				})
			}
			return nil
		}(); err != nil {
			logs.LogError.Println(err)
			if ctx.Sender() != nil {
				time.Sleep(3 * time.Second)
				ctx.Respond(&messages.MsgAck{Error: err.Error()})
			}
			a.fm.Event(a.contxt, eError, err)
		}
	case *messages.MsgCreateKeyRequest:
		if err := func() error {
			if err := a.sam.GenerateKey(msg.KeySlot, msg.Alg); err != nil {
				return err
			}
			logs.LogBuild.Printf("Key Create, %d", msg.KeySlot)
			if ctx.Sender() != nil {
				ctx.Respond(&messages.MsgAck{})
			}
			return nil
		}(); err != nil {
			logs.LogError.Println(err)
			// if ctx.Sender() != nil {
			// 	time.Sleep(3 * time.Second)
			// 	ctx.Respond(&MsgAck{Error: err.Error()})
			// }
			a.fm.Event(a.contxt, eError, err)
		}
	case *messages.MsgImportKeyRequest:
		if err := func() error {
			// if a.sam.EnableKeys() msg.KeySlot
			if err := a.sam.ImportKey(msg.Data, msg.KeySlot, msg.Alg); err != nil {
				return err
			}
			logs.LogBuild.Printf("Key import, %d", msg.KeySlot)
			if ctx.Sender() != nil {
				ctx.Respond(&messages.MsgAck{})
			}
			return nil
		}(); err != nil {
			logs.LogError.Println(err)
			// if ctx.Sender() != nil {
			// 	time.Sleep(3 * time.Second)
			// 	ctx.Respond(&MsgAck{Error: err.Error()})
			// }
			a.fm.Event(a.contxt, eError, err)
		}
		// case *messages.MsgEnableKeysRequest:
		// 	if err := func() error {
		// 		keys, err := a.sam.EnableKeys()
		// 		if err != nil {
		// 			return err
		// 		}
		// 		logs.LogBuild.Printf("enable Keys, %+v", keys)
		// 		if ctx.Sender() != nil {
		// 			ctx.Request(ctx.Sender(), &messages.MsgEnableKeysResponse{Data: keys})
		// 		}
		// 		return nil
		// 	}(); err != nil {
		// 		logs.LogError.Println(err)
		// 		// if ctx.Sender() != nil {
		// 		// 	time.Sleep(3 * time.Second)
		// 		// 	ctx.Respond(&MsgAck{Error: err.Error()})
		// 		// }
		// 		a.fm.Event(a.contxt, eError, err)
		// 	}
	}
}
