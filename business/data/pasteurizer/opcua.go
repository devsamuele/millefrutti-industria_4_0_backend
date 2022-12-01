package pasteurizer

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
	// go o.WatchOrderConf("ns=8;s=Siemens S7-1200/S7-1500.Tags.Send.Conferma_Nuovo_Lotto", 1)
	// go o.WatchEndWork("ns=8;s=Siemens S7-1200/S7-1500.Tags.Send.Fine_Produzione", 2)
	// go o.WatchBasilAmount("ns=8;s=Siemens S7-1200/S7-1500.Tags.Send.Quantità_Basilico_Lavorato", 3)
	// go o.WatchBasilPackages("ns=8;s=Siemens S7-1200/S7-1500.Tags.Send.Numero_Di_Imballi", 4)
	go o.WatchOrderConf("ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Conferma_Nuovo_Lotto", 1)
	go o.WatchEndWork("ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Fine_Produzione", 1)
	go o.WatchBasilAmount("ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Quantità_Basilico_Lavorato", 1)
	go o.WatchBasilPackages("ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Numero_Di_Imballi", 1)
}

func (o *OpcuaService) WatchOrderConf(nodeID string, clientHandle uint32) {

	opcuaconn.Subscribe(o.ctx, o.c, nodeID, clientHandle, func(data interface{}) {
		bit, _ := data.(bool)
		if bit {
			found := true
			oldWork, err := o.store.QueryActiveWork(o.ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					found = false
				} else {
					return
				}
			}

			if found && oldWork.Status == PROCESSING_STATUS_SENT {
				log.Println("PASTEURIZER SUBSCRIPTION - START WORK")
				work := oldWork
				work.Status = PROCESSING_STATUS_WORK

				tx, err := o.store.BeginTx(o.ctx)
				if err != nil {
					o.log.Println(err)
					return
				}
				defer tx.Rollback()

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

				if err := o.io.Broadcast("pasteurizer-status-change", b); err != nil {
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

func (o *OpcuaService) WatchBasilAmount(nodeID string, clientHandle uint32) {

	opcuaconn.Subscribe(o.ctx, o.c, nodeID, clientHandle, func(data interface{}) {

		found := true
		oldWork, err := o.store.QueryActiveWork(o.ctx)
		if err != nil {
			log.Println(err)
			if errors.Is(err, sql.ErrNoRows) {
				found = false
			} else {
				return
			}
		}

		if found && oldWork.Status == PROCESSING_STATUS_WORK {
			log.Println("SPINDRYER SUBSCRIPTION - START UPDATE BASIL AMOUNT")

			currentBasilAmount, _ := data.(int64)

			work := oldWork

			tx, err := o.store.BeginTx(o.ctx)
			if err != nil {
				o.log.Println(err)
				return
			}
			defer tx.Rollback()

			work.BasilAmount = int(currentBasilAmount)
			log.Println("current basil amount:", currentBasilAmount)
			// TODO: TEMPORARY
			if currentBasilAmount == 0 {
				work.BasilAmount = 400
			}

			// var packages uint16
			// newPackages, err := opcuaconn.Read(o.ctx, o.c, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Numero_Di_Imballi")
			// if err != nil {
			// 	o.log.Println(err)
			// }

			// packages, _ = newPackages.(uint16)
			// log.Println("basil packages:", packages)
			// work.Packages = int(packages)
			// // TODO: TEMPORARY
			// work.Packages = 2

			err = o.store.UpdateWork(o.ctx, tx, work)
			if err != nil {
				o.log.Println(err)
				return
			}

			if err := tx.Commit(); err != nil {
				o.log.Println(err)
				return
			}
			log.Println("SPINDRYER SUBSCRIPTION - END UPDATE BASIL AMOUNT")
		}
	})
}

func (o *OpcuaService) WatchBasilPackages(nodeID string, clientHandle uint32) {

	opcuaconn.Subscribe(o.ctx, o.c, nodeID, clientHandle, func(data interface{}) {

		found := true
		oldWork, err := o.store.QueryActiveWork(o.ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				found = false
			} else {
				return
			}
		}

		if found && oldWork.Status == PROCESSING_STATUS_WORK {
			log.Println("SPINDRYER SUBSCRIPTION - START UPDATE BASIL PACKAGE")

			currentBasilPackages, _ := data.(uint16)

			work := oldWork

			tx, err := o.store.BeginTx(o.ctx)
			if err != nil {
				o.log.Println(err)
				return
			}
			defer tx.Rollback()

			// work.BasilAmount = int(currentBasilAmount)
			// log.Println("current basil amount:", currentBasilAmount)
			// // TODO: TEMPORARY
			// if currentBasilAmount == 0 {
			// 	work.BasilAmount = 400
			// }

			// var packages uint16
			// newPackages, err := opcuaconn.Read(o.ctx, o.c, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Numero_Di_Imballi")
			// if err != nil {
			// 	o.log.Println(err)
			// }

			// packages, _ = newPackages.(uint16)
			work.Packages = int(currentBasilPackages)
			log.Println("basil packages:", currentBasilPackages)
			if work.Packages == 0 {
				work.Packages = 2
			}

			err = o.store.UpdateWork(o.ctx, tx, work)
			if err != nil {
				o.log.Println(err)
				return
			}

			if err := tx.Commit(); err != nil {
				o.log.Println(err)
				return
			}
			log.Println("SPINDRYER SUBSCRIPTION - END UPDATE BASIL PACKAGE")
		}
	})
}

func (o *OpcuaService) WatchEndWork(nodeID string, clientHandle uint32) {

	opcuaconn.Subscribe(o.ctx, o.c, nodeID, clientHandle, func(data interface{}) {
		bit, _ := data.(bool)
		if bit {
			found := true
			oldWork, err := o.store.QueryActiveWork(o.ctx)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					found = false
				} else {
					return
				}
				// o.log.Println(err)
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

				// var basilAmount int64
				// newBasilAmount, err := opcuaconn.Read(o.ctx, o.c, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Quantità_Basilico_Lavorato")
				// if err != nil {
				// 	o.log.Println(err)
				// }

				// // TODO um
				// basilAmount, _ = newBasilAmount.(int64)
				// log.Println("basil amount:", basilAmount)
				// work.BasilAmount = int(basilAmount)
				// // TODO: TEMPORARY
				// work.BasilAmount = 400

				// var packages uint16
				// newPackages, err := opcuaconn.Read(o.ctx, o.c, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Numero_Di_Imballi")
				// if err != nil {
				// 	o.log.Println(err)
				// }

				// packages, _ = newPackages.(uint16)
				// log.Println("basil packages:", packages)
				// work.Packages = int(packages)
				// // TODO: TEMPORARY
				// work.Packages = 2

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

				if err := o.io.Broadcast("pasteurizer-status-change", b); err != nil {
					o.log.Println(err)
					return
				}

				if err := tx.Commit(); err != nil {
					o.log.Println(err)
					return
				}
				log.Println("PASTEURIZER SUBSCRIPTION - END WORK")

			}

		}
	})
}
