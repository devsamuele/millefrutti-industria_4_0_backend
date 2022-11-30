package main

import (
	"context"
	"log"
	"time"

	"github.com/devsamuele/millefrutti-industria_4_0_backend/business/sys/opcuaconn"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

func main() {
	ctx := context.Background()
	pasteurizerClient := opcua.NewClient("opc.tcp://192.168.1.181:4840", opcua.SecurityMode(ua.MessageSecurityModeNone), opcua.DialTimeout(time.Second*10))
	if err := pasteurizerClient.Connect(ctx); err != nil {
		log.Println(err)
	}

	var basilAmount int64
	newBasilAmount, err := opcuaconn.Read(ctx, pasteurizerClient, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Quantità_Basilico_Lavorato")
	if err != nil {
		log.Println(err)
	}

	// TODO um
	basilAmount, _ = newBasilAmount.(int64)
	log.Println("basil amount:", basilAmount)
	log.Println(int(basilAmount))

	var packages uint16
	newPackages, err := opcuaconn.Read(ctx, pasteurizerClient, "ns=2;s=Siemens S7-1200/S7-1500.Tags.Send.Numero_Di_Imballi")
	if err != nil {
		log.Println(err)
	}

	packages, _ = newPackages.(uint16)
	log.Println("basil packages:", packages)
	log.Println(int(packages))
}
