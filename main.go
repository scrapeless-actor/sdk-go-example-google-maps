package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless"
	proxyModel "github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless/proxy"
	"github.com/scrapeless-ai/scrapeless-actor-sdk-go/scrapeless/storage/queue"
	"time"
)

const taskName = "google-maps-queues"

func main() {
	actor := scrapeless.New(scrapeless.WithProxy(), scrapeless.WithStorage())
	defer actor.Close()
	ctx := context.TODO()
	// Get params from environment variables
	var param = &RequestParam{}
	if err := actor.Input(param); err != nil {
		panic(err)
	}
	// Get proxy
	proxy := getProxy(actor, ctx)
	InitProxyClient(proxy)

	// Example of queue
	q := actor.Storage.GetQueue()
	// Create a queue if not exist
	name := uuid.New().String()[:8]

	queueId := createQueue(ctx, q, name)
	createdQueueId := actor.Storage.GetQueue(queueId)
	// Cyclically add tasks to the queue
	marshal, err := json.Marshal(param)
	if err != nil {
		panic(fmt.Sprintf("Error marshalling params: %s", err))
	}
	setQueue(ctx, createdQueueId, taskName, string(marshal))
	// Create Kv
	namespaceId, _, err := actor.Storage.GetKv().CreateNamespace(ctx, name)
	if err != nil {
		panic(fmt.Sprintf("Error creating namespace: %s", err))
	}
	kv := actor.Storage.GetKv(namespaceId)
	// Get one from the queue every 10s, and crawl Google Maps
	for {
		time.Sleep(5 * time.Second)
		msgs := getQueue(ctx, createdQueueId, 1)
		if len(msgs) == 0 {
			fmt.Printf("Queue %s is empty\n", queueId)
			continue
		}
		msg := msgs[0]
		fmt.Println(">==============begin crawl==============<")
		fmt.Printf("queue id:%s, msg name:%s, id:%s, payload:%s, desc:%s\n", queueId, msg.Name, msg.ID, msg.Payload, msg.Desc)
		var param = &RequestParam{}
		err := json.Unmarshal([]byte(msg.Payload), param)
		if err != nil {
			fmt.Printf("Error unmarshalling params: %s\n", err)
			continue
		}
		response, err := crawl(ctx, param)
		if err != nil {
			fmt.Printf("Failed to crawl, err:%s\n", err.Error())
			continue
		}
		marshal, _ := json.Marshal(response)
		// Use kv for storage
		ok, err := kv.SetValue(ctx, msg.ID, string(marshal), 0)
		if err != nil {
			fmt.Printf("Failed to set value, err:%s\n", err)
			continue
		}
		if !ok {
			fmt.Printf("Failed to set value, err:%s\n", err.Error())
			continue
		}
		value, err := kv.GetValue(ctx, msg.ID)
		if err != nil {
			fmt.Printf("Failed to get value, err:%s\n", err)
			continue
		}
		fmt.Printf("kv-->Get value:%s\n", value)
		fmt.Println(">==============end crawl==============<")
	}
}

func getProxy(actor *scrapeless.Actor, ctx context.Context) string {
	proxy, err := actor.Proxy.Proxy(ctx, proxyModel.ProxyActor{
		Country:         "us",
		SessionDuration: 10,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to get proxy, err:%s", err.Error()))
	}
	return proxy
}

func createQueue(ctx context.Context, q queue.Queue, name string) (queueId string) {
	queueId, _, err := q.Create(ctx, &queue.CreateQueueReq{
		Name:        name,
		Description: name,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create queue, err:%s", err.Error()))
	}
	return queueId
}

func getQueue(ctx context.Context, q queue.Queue, size int32) queue.GetMsgResponse {
	pull, err := q.Pull(ctx, size)
	if err != nil {
		panic(fmt.Sprintf("Failed to pull: %v", err))
	}
	return pull
}

func setQueue(ctx context.Context, q queue.Queue, taskName, payload string) {
	push, err := q.Push(ctx, queue.PushQueue{
		Name:     taskName,
		Payload:  []byte(payload),
		Retry:    0,
		Timeout:  0,
		Deadline: 0,
	})
	if err != nil {
		panic(fmt.Sprintf("push failed: %v", err))
	}
	fmt.Println("push queue:", push)
}

func crawl(ctx context.Context, param *RequestParam) (response *Response, err error) {
	defer func() {
		a := recover()
		if a != nil {
			fmt.Println("crawl failed:", a)
		}
	}()
	err = param.FieldValidation()
	if err != nil {
		return &Response{}, err
	}
	switch param.Engine {
	case GoogleMaps:
		response, err = DoMaps(ctx, param)
		return response, err
	case GoogleMapsAutocomplete:
		response, err = DoMapsAutocomplete(ctx, param)
		return response, err
	case GoogleMapsContributorReviews:
		response, err = DoMapsContributorReviews(ctx, param)
		return response, err
	case GoogleMapsDirections:
		response, err = DoMapsDirections(ctx, param)
		return response, err
	case GoogleMapsPhotos:
		response, err = DoMapsPhotos(ctx, param)
		return response, err
	case GoogleMapsPhotoMeta:
		response, err = DoMapsPhotoMeta(ctx, param)
		return response, err
	case GoogleMapsReviews:
		response, err = DoMapsReviews(ctx, param)
		return response, err
	default:
		return response, fmt.Errorf("unsupported engine %s", param.Engine)
	}
}
