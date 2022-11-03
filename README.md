The objective of this project is create a wrapper of RSocket for Golang for JSON payloads and triggering methods on the server side, like a web request in gorilla.

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

func execute_something(payload interface{}) string{
	fmt.Println("execute_something")
	fmt.Println(payload)
	return "execute_something"
}


func dont_execute_something(payload interface{}) string{
	fmt.Println("dont_execute_something")
	fmt.Println(payload)
	return "dont_execute_something"
}

```


Example Client:
```golang
package main

import (
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

	var list []peers

	var p peers
	p.Name = "Test"
	p.Address = "Street"
	list = append(list, p)

	var p1 peers
	p1.Name = "Test1"
	p1.Address = "Street1"
	list = append(list, p1)
	var list_peers peers_cont
	list_peers.Peers = list

	rsocket_json_requests.RequestConfigs("127.0.0.1", 7878)
	rsocket_json_requests.RequestJSON("execute_something", list_peers)
}
```
TODO:
- add TLS certificate
- ...
