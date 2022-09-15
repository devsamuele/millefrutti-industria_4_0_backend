package pasteurizer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/sys/opcuaconn"
	"github.com/devsamuele/service-kit/web"
	"github.com/devsamuele/service-kit/ws"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type Service struct {
	store    Store
	client   *opcua.Client
	io       *ws.EventEmitter
	log      *log.Logger
	shutdown chan os.Signal
}

func NewService(store Store, shutdown chan os.Signal, log *log.Logger, io *ws.EventEmitter) *Service {
	return &Service{
		store:    store,
		io:       io,
		log:      log,
		shutdown: shutdown,
	}
}

// Pasteurizer "opc.tcp://192.168.1.181:4840"

func (s *Service) OpcuaConnect(ctx context.Context) error {
	pasteurizerClient := opcua.NewClient("opc.tcp://192.168.1.181:4840", opcua.SecurityMode(ua.MessageSecurityModeNone), opcua.DialTimeout(time.Second*10))
	if err := pasteurizerClient.Connect(ctx); err != nil {
		return web.NewError("pasteurizer not connected", web.ErrReasonInternalError, "", "")
	}

	_ctx, cancel := context.WithCancel(context.Background())

	s.client = pasteurizerClient
	opcuaService := NewOpcuaService(_ctx, s.log, s.client, s.store, s.io)
	opcuaService.Run()

	go func() {
		defer func() {
			s.client.CloseWithContext(_ctx)
			cancel()
		}()
		for {
			if s.client.State() != opcua.Connected {
				s.client.CloseWithContext(_ctx)

				if err := s.io.Broadcast("pasteurizer-client-closed", nil); err != nil {
					log.Println(err)
				}
				return
			}
			time.Sleep(time.Second * 5)
		}
	}()
	return nil
}

func (s *Service) OpcuaDisconnect(ctx context.Context) error {
	if s.client != nil {
		if err := s.client.CloseWithContext(ctx); err != nil {
			log.Println("closing pasteurizer:", err)
		}
	}

	return nil
}

func (s Service) QueryWork(ctx context.Context) ([]Work, error) {
	works, err := s.store.QueryWork(ctx)
	if err != nil {
		return make([]Work, 0), err
	}
	return works, nil
}

func (s Service) InsertWork(ctx context.Context, nw NewWork, now time.Time) (Work, error) {

	if s.client == nil {
		return Work{}, web.NewError("pasteurizer is not connected", web.ErrReasonInternalError, "", "")
	}

	if err := nw.Validate(); err != nil {
		return Work{}, err
	}

	exist, err := s.store.ExistActiveWork(ctx)
	if err != nil {
		return Work{}, err
	}

	if exist {
		return Work{}, errors.New("active work already exist")
	}

	w := Work{
		CdLotto:         *nw.CdLotto,
		CdAr:            *nw.CdAr,
		DocumentCreated: false,
		Date:            now,
		Status:          PROCESSING_STATUS_SENT,
		Created:         now,
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return Work{}, err
	}

	defer tx.Rollback()

	found, err := s.store.CheckLottoAndAr(ctx, tx, w.CdLotto, w.CdAr)
	if err != nil {
		return Work{}, err
	}
	if !found {
		err = s.store.CreateLottoArca(ctx, tx, w.CdLotto, w.CdAr, now)
		if err != nil {
			return Work{}, err
		}
	}

	id, err := s.store.InsertWork(ctx, tx, w)
	if err != nil {
		return Work{}, err
	}
	w.ID = id

	_, err = opcuaconn.Write(ctx, s.client, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Receive.Numero_Lotto", w.CdLotto)
	if err != nil {
		return Work{}, err
	}

	var bit bool = true
	_, err = opcuaconn.Write(ctx, s.client, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Receive.Bit_Nuovo_Lotto", bit)
	if err != nil {
		return Work{}, err
	}

	if err := tx.Commit(); err != nil {
		return Work{}, err
	}

	return w, nil
}

func (s Service) GetOpcuaConnection(ctx context.Context) OpcuaConnection {

	if s.client != nil && s.client.State() == opcua.Connected {
		return OpcuaConnection{
			Connected: true,
		}
	}

	return OpcuaConnection{
		Connected: false,
	}
}

func (s Service) SetCreatedDocument(ctx context.Context, ids []ID) error {

	works := make([]Work, 0)
	for _, id := range ids {
		w, err := s.store.QueryWorkByID(ctx, id.ID)
		if err != nil {
			return err
		}
		works = append(works, w)
	}

	b, err := json.Marshal(&works)
	if err != nil {
		return err
	}

	if err := s.io.Broadcast("pasteurizer-created-documents", b); err != nil {
		return err
	}

	return nil
}

func (s Service) DeleteWork(ctx context.Context, id string) error {
	_id, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	tx, err := s.store.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	w, err := s.store.QueryWorkByID(ctx, _id)
	if err != nil {
		return err
	}

	// if w.Status != "send" {
	// 	return errors.New("unable to delete already sent work")
	// }

	err = s.store.DeleteWork(ctx, tx, _id)
	if err != nil {
		return err
	}

	found, err := s.store.CheckLottoAndArInDoc(ctx, tx, w.CdLotto, w.CdAr)
	if err != nil {
		return err
	}

	if !found {
		err = s.store.DeleteLottoArca(ctx, tx, w.CdLotto, w.CdAr)
		if err != nil {
			return err
		}
	}

	// var cdLotto string = " "
	// _, err = opcuaconn.Write(ctx, s.client, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Receive.Numero_Lotto", cdLotto)
	// if err != nil {
	// 	return err
	// }

	// var bit bool = false
	// _, err = opcuaconn.Write(ctx, s.client, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Receive.Bit_Nuovo_Lotto", bit)
	// if err != nil {
	// 	return err
	// }

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
