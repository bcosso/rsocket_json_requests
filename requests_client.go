package rsocket_json_requests

import (
	"context"
	"log"
	"fmt"
	"encoding/json"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

type GenericList struct{
	Method string `json:"method"`
	Payload interface{} `json:"payload"`
}

var _ip string
var _port int

var _genericList GenericList


func RequestConfigs(ip string, port int){
	_ip = ip
	_port = port
}

func RequestJSON(method string, json_content interface {}) interface {} {
	// Connect to server
	var result_json interface{}
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("", "")).
		Transport(rsocket.TCPClient().SetHostAndPort(_ip, _port).Build()).
		Start(context.Background())
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	_genericList.Method = method
	method = "{\"method\":\""+ method +"\"}"
	data := []byte(method)

	
	_genericList.Payload = json_content
	
	meta_data, err := json.Marshal(_genericList)
	result, err := cli.RequestResponse(payload.New(meta_data, data)).Block(context.Background())
	if err != nil {
		panic(err)
	}
	
	err = json.Unmarshal(result.Data(), &result_json)

	if err!= nil{
		fmt.Println(err)
		log.Fatal(err)
		fmt.Println(err)
	}
	return result_json
}
