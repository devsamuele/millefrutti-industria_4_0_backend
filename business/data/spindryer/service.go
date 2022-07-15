package spindryer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/sys/opcuaconn"
	"github.com/devsamuele/service-kit/ws"
	"github.com/gopcua/opcua"
)

type Service struct {
	store  Store
	client *opcua.Client
	io     *ws.EventEmitter
}

func NewService(store Store, client *opcua.Client, io *ws.EventEmitter) Service {
	return Service{
		store:  store,
		client: client,
		io:     io,
	}
}

func (s Service) QueryWork(ctx context.Context) ([]Work, error) {
	works, err := s.store.QueryWork(ctx)
	if err != nil {
		return make([]Work, 0), err
	}
	return works, nil
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

	if err := s.io.Broadcast("spindryer-created-documents", b); err != nil {
		return err
	}

	return nil
}

func (s Service) GetOpcuaConnection(ctx context.Context) OpcuaConnection {

	return OpcuaConnection{
		Connected: OpcuaConnected,
	}

}


func (s Service) InsertWork(ctx context.Context, nw NewWork, now time.Time) (Work, error) {

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
		Cycles:          0,
		TotalCycles:     0,
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

	_, err = opcuaconn.Write(ctx, s.client, "ns=2;s=DB_REPORT_4_0_LOTTO_DA_MES", w.CdLotto)
	if err != nil {
		return Work{}, err
	}

	var bit bool = true
	_, err = opcuaconn.Write(ctx, s.client, "ns=2;s=DB_REPORT_4_0_BIT_NUOVO_ORD_DA_MES", bit)
	if err != nil {
		return Work{}, err
	}

	// v, err := api.Read(ctx, c, "ns=2;s=DB_REPORT_4_0_LOTTO_DA_MES")
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Printf("%v", v)

	if err := tx.Commit(); err != nil {
		return Work{}, err
	}

	return w, nil
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

	err = s.store.DeleteLottoArca(ctx, tx, w.CdLotto, w.CdAr)
	if err != nil {
		return err
	}

	var cdLotto string
	_, err = opcuaconn.Write(ctx, s.client, "ns=2;s=DB_REPORT_4_0_LOTTO_DA_MES", cdLotto)
	if err != nil {
		log.Println(err)
	}

	var bit bool = false
	_, err = opcuaconn.Write(ctx, s.client, "ns=2;s=DB_REPORT_4_0_BIT_NUOVO_ORD_DA_MES", bit)
	if err != nil {
		log.Println(err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
