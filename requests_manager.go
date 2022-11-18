package rsocket_json_requests

import (
	"context"
	"fmt"
	"log"
	"encoding/json"
	"crypto/tls"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
)


func AppendFunctionHandler(name string, function func(pload interface{}) interface{} ) {
	var func1 FunctionName;
	func1.Function = function
	func1.Name = name
	func_list = append(func_list, func1)
}

var tc * tls.Config

type FunctionName struct{
	Function func(pload interface{}) interface{};
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
				var resp interface{}
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
					
					if parsed_document["method"].(string) == f.Name{
						fmt.Println(f)
						resp = f.Function(parsed_document["payload"].(interface{}))
					}
				}
				method := "{\"status\":\"200\"}"
				data := []byte(method)

				
				meta_data, err := json.Marshal(resp)
				if err!= nil{
					fmt.Println(err)
					log.Fatal(err)
					fmt.Println(err)
				}

				return mono.Just(payload.New(meta_data, data))
			}),
		), nil
	}).
	Transport((func() *rsocket.TCPServerBuilder {
		var builder_result *rsocket.TCPServerBuilder
		builder_result = rsocket.TCPServer()
			if tc != nil{
				builder_result = builder_result.SetTLSConfig(tc)	
			}
			return builder_result
		}()).SetAddr(":7878").Build()).
	Serve(context.Background())
	log.Fatalln(err)
}

func SetTLSConfig(cert_path string, key_path string){
	cert, err := tls.LoadX509KeyPair(cert_path, key_path)
	if err != nil {
		panic(err)
	}
	// Init TLS configuration.
	tc = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: []tls.Certificate{cert},
	}
}
