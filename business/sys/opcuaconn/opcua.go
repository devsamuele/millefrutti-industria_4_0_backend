package opcuaconn

import (
	"context"
	"fmt"
	"log"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type SubscribeFn func(ctx context.Context, c *opcua.Client, nodeID string, clientHandle uint32, callback func(data interface{}))

func Subscribe(ctx context.Context, c *opcua.Client, nodeID string, clientHandle uint32, callback func(data interface{})) {

	notifyCh := make(chan *opcua.PublishNotificationData)

	sub, err := c.SubscribeWithContext(ctx, &opcua.SubscriptionParameters{
		Interval: opcua.DefaultSubscriptionInterval,
	}, notifyCh)
	if err != nil {
		log.Println(err)
	}
	defer sub.Cancel(ctx)

	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		log.Println(err)
	}

	miCreateRequest := opcua.NewMonitoredItemCreateRequestWithDefaults(id, ua.AttributeIDValue, clientHandle)
	res, err := sub.Monitor(ua.TimestampsToReturnBoth, miCreateRequest)
	if err != nil || res.Results[0].StatusCode != ua.StatusOK {
		if err == nil {
			log.Println(res.Results[0].StatusCode)
		} else {
			log.Println(err)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case res := <-notifyCh:
			if res.Error != nil {
				log.Print(res.Error)
				continue
			}

			switch x := res.Value.(type) {
			case *ua.DataChangeNotification:
				for _, item := range x.MonitoredItems {
					data := item.Value.Value.Value()
					callback(data)
				}

			default:
				log.Printf("unhandled result %T", res.Value)
			}
		}
	}

}

// func Subscribe(ctx context.Context, c *opcua.Client, nodeID string) {

// 	notifyCh := make(chan *opcua.PublishNotificationData)

// 	sub, err := c.SubscribeWithContext(ctx, &opcua.SubscriptionParameters{
// 		Interval: opcua.DefaultSubscriptionInterval,
// 	}, notifyCh)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	id, err := ua.ParseNodeID(nodeID)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	miCreateRequest := opcua.NewMonitoredItemCreateRequestWithDefaults(id, ua.AttributeIDValue, uint32(22))
// 	res, err := sub.Monitor(ua.TimestampsToReturnBoth, miCreateRequest)
// 	if err != nil || res.Results[0].StatusCode != ua.StatusOK {
// 		log.Println(err)
// 	}

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case res := <-notifyCh:
// 			if res.Error != nil {
// 				log.Print(res.Error)
// 				continue
// 			}

// 			switch x := res.Value.(type) {
// 			case *ua.DataChangeNotification:
// 				for _, item := range x.MonitoredItems {
// 					data := item.Value.Value.Value()
// 					log.Printf("MonitoredItem with client handle %v = %v", item.ClientHandle, data)
// 				}

// 			default:
// 				log.Printf("what's this publish result? %T", res.Value)
// 			}
// 		}
// 	}
// }

func Write(ctx context.Context, c *opcua.Client, nodeID string, value interface{}) ([]ua.StatusCode, error) {

	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return nil, err
	}

	v, err := ua.NewVariant(value)
	if err != nil {
		return nil, err
	}

	wReq := ua.WriteRequest{
		NodesToWrite: []*ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					Value:        v,
					EncodingMask: ua.DataValueValue,
				},
			},
		},
	}

	wResp, err := c.WriteWithContext(ctx, &wReq)
	if err != nil {
		return nil, err
	}
	return wResp.Results, nil
}

func Read(ctx context.Context, c *opcua.Client, nodeID string) (interface{}, error) {

	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return nil, err
	}

	rReq := ua.ReadRequest{
		NodesToRead:        []*ua.ReadValueID{{NodeID: id}},
		MaxAge:             2000,
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	rResp, err := c.ReadWithContext(ctx, &rReq)
	if err != nil {
		return nil, err
	}

	if rResp.Results[0].Status != ua.StatusOK {
		return nil, fmt.Errorf("status not OK: %v", rResp.Results[0].Status)
	}

	return rResp.Results[0].Value.Value(), nil
}
