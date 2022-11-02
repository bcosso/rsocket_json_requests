package rsocket_json_requests

import (
	"context"
	"fmt"
	"log"
	"encoding/json"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"

)


func AppendFunctionHandler(name string, function func() string ) {
	var func1 FunctionName;
	func1.Function = function
	func1.Name = name
	func_list = append(func_list, func1)
}

type FunctionName struct{
	Function func() string;
	Name string 
}

var func_list []FunctionName

func ServeCalls(){
	err := rsocket.Receive().
	Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
		// bind responder
		return rsocket.NewAbstractSocket(
			rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
				var mt interface{}
				fmt.Println(msg)
				err := json.Unmarshal(msg.Data(), &mt)


				if err!= nil{
					fmt.Println(err)
					log.Fatal(err)
					fmt.Println(err)
				}

				fmt.Println(mt)

				parsed_document, ok :=  mt.(map[string] interface{})
				//err = json.Unmarshal(current_document, &parsed_document)
		
				if !ok{
					fmt.Println("ERROR!")
					
				}
				for _, f := range func_list{
					
					if parsed_document["method"] == f.Name{
						fmt.Println(f)
						f.Function()
					}
				}
				
				return mono.Just(msg)
			}),
		), nil
	}).
	Transport(rsocket.TCPServer().SetAddr(":7878").Build()).
	Serve(context.Background())
	log.Fatalln(err)
}