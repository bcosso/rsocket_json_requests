package rsocket_json_requests

import (
	"context"
	"log"
	"encoding/json"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

var _ip string
var _port int

func RequestConfigs(ip string, port int){
	_ip = ip
	_port = port
}

func RequestJSON(method string, json_content interface {}) {
	// Connect to server
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("", "")).
		Transport(rsocket.TCPClient().SetHostAndPort(_ip, _port).Build()).
		Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	method = "{\"method\":\""+ method +"\"}"
	data := []byte(method)
	meta_data, err := json.Marshal(json_content)
	result, err := cli.RequestResponse(payload.New(meta_data, data)).Block(context.Background())
	//result, err := cli.RequestResponse(payload.NewString("{\"method\":\""+ method +"\"}", "{\"tesssst\":\"field\"}")).Block(context.Background())
	if err != nil {
		panic(err)
	}
	log.Println("response:", result)
}