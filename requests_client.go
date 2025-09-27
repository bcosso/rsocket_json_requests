package rsocket_json_requests

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
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
	// Cli      rsocket.Client
	FreeList map[string]map[string]*ClientRsockDetails
	InUse    map[string]map[string]*ClientRsockDetails
	mu       sync.RWMutex
}

type ClientRsockDetails struct {
	// Cli      rsocket.Client
	Client       *rsocket.Client
	HostName     string
	InstanceName string
	InternalId   string
	Port         int
}

var _cli ClientRsock

// Shows status of the RSsock client with the Unique Instance Name. Returns 1 if there is a free connection, 2 if there is a connection in use and -1
// if there is no connection
//
// Parameters:
//   - instanceName: unique name of the instance
func (cli *ClientRsock) GetStatusConnection(instanceName string) int {
	//Read lock maybe
	if len(cli.FreeList[instanceName]) > 0 {
		return 1
	} else if len(cli.InUse[instanceName]) > 0 {
		return 2
	} else {
		return -1
	}
}

func (cli *ClientRsock) GetExistingOrNew(host string) (*ClientRsockDetails, error) {
	var result *ClientRsockDetails
	if len(cli.FreeList[host]) > 0 {
		var firstKey string
		for k, _ := range cli.FreeList[host] {
			if cli.FreeList[host][k] != nil {
				firstKey = k
				break
			}
		}
		result = cli.FreeList[host][firstKey]
		fmt.Println(result)
		cli.mu.Lock()
		defer cli.mu.Unlock()
		cli.InUse[host][firstKey] = result
		delete(cli.FreeList[host], firstKey)
		// cli.FreeList[host] = slices.Delete(cli.FreeList[host], 0, 1)
	} else if len(cli.InUse[host]) > 0 {
		var firstKey string
		for k, _ := range cli.InUse[host] {
			if cli.InUse[host][k].Client != nil {
				firstKey = k
				break
			}
		}

		newResult := *cli.InUse[host][firstKey]
		InitConn(newResult.InstanceName, newResult.HostName, newResult.Port)
		for k, _ := range cli.FreeList[host] {
			if cli.FreeList[host][k].Client != nil {
				firstKey = k
				break
			}
		}
		cli.mu.Lock()
		defer cli.mu.Unlock()
		result = cli.FreeList[host][firstKey]
	} else {
		newError := fmt.Errorf("Unique instances not initialized, please use InitConn method")
		return nil, newError
	}
	return result, nil
}

func (cli *ClientRsock) FreeExistingSocket(client *ClientRsockDetails) error {
	cli.mu.Lock()
	defer cli.mu.Unlock()
	if _, ok := cli.InUse[client.InstanceName][client.InternalId]; ok {
		delete(cli.InUse[client.InstanceName], client.InternalId)
	}
	cli.FreeList[client.InstanceName][client.InternalId] = client

	return nil
}

func (cli *ClientRsock) AddRSockClient(client *rsocket.Client, instanceName string, hostName string, port int) {
	cli.mu.Lock()
	defer cli.mu.Unlock()

	if cli.InUse == nil {
		cli.InUse = make(map[string]map[string]*ClientRsockDetails)
	}
	if cli.FreeList == nil {
		cli.FreeList = make(map[string]map[string]*ClientRsockDetails)
	}
	if _, ok := cli.FreeList[instanceName]; !ok {
		cli.FreeList[instanceName] = make(map[string]*ClientRsockDetails)
	}
	if _, ok := cli.InUse[instanceName]; !ok {
		cli.InUse[instanceName] = make(map[string]*ClientRsockDetails)
	}

	var clientDetail ClientRsockDetails
	(*client).OnClose(func(err error) {
		if _, ok := _cli.FreeList[instanceName][clientDetail.InternalId]; ok {
			delete(_cli.FreeList[instanceName], clientDetail.InternalId)
		}
		if _, ok := _cli.InUse[instanceName][clientDetail.InternalId]; ok {
			delete(_cli.InUse[instanceName], clientDetail.InternalId)
		}
	})
	clientDetail.Client = client
	clientDetail.HostName = hostName
	clientDetail.InstanceName = instanceName
	clientDetail.Port = port
	clientDetail.InternalId = uuid.New().String()

	cli.FreeList[instanceName][clientDetail.InternalId] = &clientDetail

}

// Starts singleton connection with server
//
// Parameters:
//   - name: unique name of the instance
//   - host: the host or ip address
//   - port: the port of the host
func InitConn(name string, host string, port int) error {
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

	_cli.AddRSockClient(&cli, name, host, port)

	return nil

}

func GetStatusConn(name string) bool {

	if result := _cli.GetStatusConnection(name); result > -1 {
		return true
	}
	return false
}

func CloseAllConn() {

	for k1, _ := range _cli.FreeList {
		for k2, _ := range _cli.FreeList[k1] {
			(*_cli.FreeList[k1][k2].Client).Close()
		}
	}
	for k1, _ := range _cli.InUse {
		for k2, _ := range _cli.InUse[k1] {
			(*_cli.InUse[k1][k2].Client).Close()
		}
	}

	for k1, _ := range _cli.InUse {
		_cli.InUse[k1] = make(map[string]*ClientRsockDetails)

	}
	_cli.InUse = make(map[string]map[string]*ClientRsockDetails)

	for k1, _ := range _cli.FreeList {
		_cli.FreeList[k1] = make(map[string]*ClientRsockDetails)

	}
	_cli.FreeList = make(map[string]map[string]*ClientRsockDetails)

}

// Creates a request with an existing connection to an instance.
//
// Parameters:
//   - method: name of the method to be called in the endpoint
//   - json_content: parameters to be sent to the endpoint in JSON
//   - name: unique name of the instance
func RequestJSONNew(method string, json_content interface{}, instanceName string) (interface{}, error) {
	// Connect to server
	var result_json interface{}

	// defer cli.Close()
	_genericList.Method = method
	method = "{\"method\":\"" + method + "\"}"
	data := []byte(method)

	_genericList.Payload = json_content

	meta_data, err := jsonIterGlobal.Marshal(_genericList)

	currentCli, err := _cli.GetExistingOrNew(instanceName)
	if err != nil {
		panic(err)
	}
	result, err := (*(*currentCli).Client).RequestResponse(payload.New(meta_data, data)).Block(context.Background())
	if err != nil {
		//panic(err)
		return nil, err
	}
	_cli.FreeExistingSocket(currentCli)

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








