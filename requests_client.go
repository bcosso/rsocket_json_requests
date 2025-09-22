package rsocket_json_requests

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

type GenericList struct {
	Method  string      `json:"method"`
	Payload interface{} `json:"payload"`
}

var _ip string
var _port int

var _genericList GenericList
var _use_TLS bool

func RequestConfigs(ip string, port int) {
	_ip = ip
	_port = port
}

func UseTLS() {
	_use_TLS = true
}

type ClientRsock struct {
	Cli rsocket.Client
}

var _cli map[string]ClientRsock

// Starts singleton connection with server
//
// Parameters:
//   - name: unique name of the instance
//   - host: the host or ip address
//   - port: the port of the host
func InitConn(name string, host string, port int) error {
	if _cli == nil {
		_cli = make(map[string]ClientRsock)
	}
	_tc := &tls.Config{
		InsecureSkipVerify: true,
	}
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("", "")).
		Transport((func() *rsocket.TCPClientBuilder {

			var builder_result *rsocket.TCPClientBuilder
			builder_result = rsocket.TCPClient()
			if _use_TLS == true {
				builder_result = builder_result.SetTLSConfig(_tc)
			}
			return builder_result

		}()).SetHostAndPort(host, port).Build()).
		Start(context.Background())
	if err != nil {
		//panic(err)
		return err
	}
	var cliStruct ClientRsock
	cliStruct.Cli = cli
	_cli[name] = cliStruct

	return nil

}

func CloseConn(name string) {
	_cli[name].Cli.Close()
}

// Creates a request with an existing connection to an instance.
//
// Parameters:
//   - method: name of the method to be called in the endpoint
//   - json_content: parameters to be sent to the endpoint in JSON
//   - name: unique name of the instance
func RequestJSONNew(method string, json_content interface{}, name string) (interface{}, error) {
	// Connect to server
	var result_json interface{}

	// defer cli.Close()
	_genericList.Method = method
	method = "{\"method\":\"" + method + "\"}"
	data := []byte(method)

	_genericList.Payload = json_content

	meta_data, err := jsonIterGlobal.Marshal(_genericList)
	result, err := _cli[name].Cli.RequestResponse(payload.New(meta_data, data)).Block(context.Background())
	if err != nil {
		//panic(err)
		return nil, err
	}

	err = jsonIterGlobal.Unmarshal(result.Data(), &result_json)

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		fmt.Println(err)
		return nil, err
	}
	return result_json, nil
}

func RequestJSON(method string, json_content interface{}) (interface{}, error) {
	// Connect to server
	var result_json interface{}
	_tc := &tls.Config{
		InsecureSkipVerify: true,
	}
	cli, err := rsocket.Connect().
		SetupPayload(payload.NewString("", "")).
		Transport((func() *rsocket.TCPClientBuilder {

			var builder_result *rsocket.TCPClientBuilder
			builder_result = rsocket.TCPClient()
			if _use_TLS == true {
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
	method = "{\"method\":\"" + method + "\"}"
	data := []byte(method)

	_genericList.Payload = json_content

	meta_data, err := jsonIterGlobal.Marshal(_genericList)
	result, err := cli.RequestResponse(payload.New(meta_data, data)).Block(context.Background())
	if err != nil {
		//panic(err)
		return nil, err
	}

	err = jsonIterGlobal.Unmarshal(result.Data(), &result_json)

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		fmt.Println(err)
		return nil, err
	}
	return result_json, nil
}






