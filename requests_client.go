package rsocket_json_requests

import (
	"context"
	"log"
	"fmt"
	"encoding/json"
	"crypto/tls"
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
var _use_TLS bool


func RequestConfigs(ip string, port int){
	_ip = ip
	_port = port
}

func UseTLS(){
	_use_TLS = true;
}


func RequestJSON(method string, json_content interface {}) (interface {}, error) {
	// Connect to server
	var result_json interface{}
	_tc := &tls.Config{
		InsecureSkipVerify: true,
	}
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("", "")).
		Transport((func() *rsocket.TCPClientBuilder{
			
				var builder_result *rsocket.TCPClientBuilder
				builder_result = rsocket.TCPClient()
				if _use_TLS == true{
					builder_result = builder_result.SetTLSConfig(_tc)	
				}
				return builder_result
				
		}()).SetHostAndPort(_ip, _port).Build()).
		Start(context.Background())
	if err != nil {
		//panic(err)
		return nil, err
	}
	defer cli.Close()
	_genericList.Method = method
	method = "{\"method\":\""+ method +"\"}"
	data := []byte(method)

	
	_genericList.Payload = json_content
	
	meta_data, err := json.Marshal(_genericList)
	result, err := cli.RequestResponse(payload.New(meta_data, data)).Block(context.Background())
	if err != nil {
		//panic(err)
		return nil, err
	}
	
	err = json.Unmarshal(result.Data(), &result_json)

	if err!= nil{
		fmt.Println(err)
		log.Fatal(err)
		fmt.Println(err)
		return nil, err
	}
	return result_json, nil
}

