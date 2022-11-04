The objective of this project is create a wrapper of RSocket for Golang for JSON payloads and triggering methods on the server side, like a HTTP Route (a web request in gorilla, for example).

Example Server:
```golang
package main

import (
	"fmt"
	"rsocket_json_requests"
)

func main() {
	rsocket_json_requests.AppendFunctionHandler("execute_something", execute_something)
	rsocket_json_requests.AppendFunctionHandler("dont_execute_something", dont_execute_something)
	rsocket_json_requests.ServeCalls()
}

func execute_something(payload interface{}) interface{}{
	fmt.Println("execute_something")
	fmt.Println(payload)
	return payload
}


func dont_execute_something(payload interface{}) interface{}{
	fmt.Println("dont_execute_something")
	fmt.Println(payload)
	return payload
}

```


Example Client:
```golang
package main

import (
	"fmt"
	"rsocket_json_requests"
)

type peers_cont struct {
	Peers []peers `json:"peers"`

}

type peers struct {
	Name string `json:"name"`
	Address string `json:"address"`

}

func main(){

	list_peers := getObject() 

	rsocket_json_requests.RequestConfigs("127.0.0.1", 7878)
	result := rsocket_json_requests.RequestJSON("execute_something", list_peers)
	fmt.Println(result)

}

func getObject() peers_cont {
	var list_peers peers_cont
	var p peers
	p.Name = "Test"
	p.Address = "Street"
	list_peers.Peers = append(list_peers.Peers, p)

	var p1 peers
	p1.Name = "Test1"
	p1.Address = "Street1"
	list_peers.Peers = append(list_peers.Peers, p1)
	return list_peers
}

```
TODO:
- add TLS certificate
- ...
