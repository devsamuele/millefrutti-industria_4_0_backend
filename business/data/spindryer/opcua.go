package spindryer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/sys/opcuaconn"
	"github.com/devsamuele/service-kit/ws"
	"github.com/gopcua/opcua"
)

type OpcuaService struct {
	ctx   context.Context
	c     *opcua.Client
	log   *log.Logger
	store Store
	io    *ws.EventEmitter
}

func NewOpcuaService(ctx context.Context, log *log.Logger, c *opcua.Client, store Store, io *ws.EventEmitter) *OpcuaService {
	return &OpcuaService{
		ctx:   ctx,
		c:     c,
		log:   log,
		store: store,
		io:    io,
	}
}

func (o *OpcuaService) Run() {
	go o.WatchOrderConf("ns=2;s=DB_REPORT_4_0_BIT_NUOVO_ORD_CONF", 1)
	go o.WatchEndWork("ns=2;s=DB_REPORT_4_0_IMP_IN_CICLO_AUT", 1)
}

func (o *OpcuaService) WatchOrderConf(nodeID string, clientHandle uint32) {

	opcuaconn.Subscribe(o.ctx, o.c, nodeID, clientHandle, func(data interface{}) {

		log.Println("SPINDRYER SUBSCRIPTION START WORK:", data)
		bit, _ := data.(bool)
		if bit {
			found := true
			oldWork, err := o.store.QueryActiveWork(o.ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					found = false
				} else {
					o.log.Println(err)
					return
				}
			}

			if found && oldWork.Status == PROCESSING_STATUS_SENT {
				work := oldWork
				work.Status = PROCESSING_STATUS_WORK

				tx, err := o.store.BeginTx(o.ctx)
				if err != nil {
					o.log.Println(err)
					return
				}
				defer tx.Rollback()

				var totalCycles int32
				newTotalCycles, err := opcuaconn.Read(o.ctx, o.c, "ns=2;s=DB_REPORT_4_0_BATCH_TOTALIZZATORE")
				if err != nil {
					o.log.Println(err)
				}

				totalCycles, _ = newTotalCycles.(int32)
				log.Println("total cycles work:", totalCycles)
				work.TotalCycles = int(totalCycles)

				err = o.store.UpdateWork(o.ctx, tx, work)
				if err != nil {
					o.log.Println(err)
					return
				}

				b, err := json.Marshal(&work)
				if err != nil {
					o.log.Println(err)
					return
				}

				if err := o.io.Broadcast("spindryer-status-change", b); err != nil {
					o.log.Println(err)
					return
				}

				if err := tx.Commit(); err != nil {
					o.log.Println(err)
					return
				}

			}

		}
	})
}

func (o *OpcuaService) WatchEndWork(nodeID string, clientHandle uint32) {

	opcuaconn.Subscribe(o.ctx, o.c, nodeID, clientHandle, func(data interface{}) {
		log.Println("SPINDRYER SUBSCRIPTION END WORK:", data)
		bit, _ := data.(bool)
		if !bit {
			found := true
			oldWork, err := o.store.QueryActiveWork(o.ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					found = false
				} else {
					o.log.Println(err)
					return
				}
			}
			if found && oldWork.Status == PROCESSING_STATUS_WORK {
				work := oldWork
				work.Status = PROCESSING_STATUS_DONE

				tx, err := o.store.BeginTx(o.ctx)
				if err != nil {
					o.log.Println(err)
					return
				}
				defer tx.Rollback()

				var totalCycles int32
				newTotalCycles, err := opcuaconn.Read(o.ctx, o.c, "ns=2;s=DB_REPORT_4_0_BATCH_TOTALIZZATORE")
				if err != nil {
					o.log.Println(err)
				}

				totalCycles, _ = newTotalCycles.(int32)
				log.Println("total cycles done:", totalCycles)

				// TODO multiply by K
				work.Cycles = int(totalCycles) - oldWork.TotalCycles

				err = o.store.UpdateWork(o.ctx, tx, work)
				if err != nil {
					o.log.Println(err)
					return
				}

				b, err := json.Marshal(&work)
				if err != nil {
					o.log.Println(err)
					return
				}

				if err := o.io.Broadcast("spindryer-status-change", b); err != nil {
					o.log.Println(err)
					return
				}

				if err := tx.Commit(); err != nil {
					o.log.Println(err)
					return
				}

			}

		}
	})
}
