package bludev

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	floodchromeClient "github.com/flood-io/go-wrenches/floodchromeclient/client"
	floodchromeTest"github.com/flood-io/go-wrenches/floodchromeclient/test"
)

type BLUDev struct {
	LaunchDevtoolsMode bool
}

func (b *BLUDev) floodchromeClient() (client *floodchromeClient.FloodchromeProto, err error) {
	client = floodchromeClient.NewHTTPClientWithConfig(nil, &floodchromeClient.TransportConfig{
		Host:     "localhost:5000",
		BasePath: "/",
		Schemes:  []string{"https"},
	})

	bearerTokenAuth := httptransport.BearerToken(os.Getenv("API_ACCESS_TOKEN"))

	return
}

func (b *BLUDev) Run(scriptFile string) (err error) {
	fmt.Println("running dev-blu")
	fmt.Printf("scriptFile = %+v\n", scriptFile)

	f, err := os.Open(scriptFile)
	if err != nil {
		return
	}

	scriptBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	client, err := b.floodchromeClient()
	if err != nil {
		return
	}

	params:=&floodchromeTest.RunParams{}
	resp, err := client.Test.Run(Run(params, bearerToken)
		if err!=nil {
		}

	fmt.Printf("client = %+v\n", client)

	test := &pb.TestRequest{
		Script: string(scriptBytes),
	}

	fmt.Println("streaming")
	stream, err := client.Run(context.Background(), test)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream %v.Run(_) = _, %v", client, err)
		}

		// fmt.Printf("result = %+v\n", result)
		// fmt.Printf("result = %T\n", result)
		// fmt.Println("result", result.String())

		// dispatch
		if logM := result.GetLog(); logM != nil {
			fmt.Printf("[%5s] %+v\n", logM.Level, result.Message)
		} else if measurementM := result.GetMeasurement(); measurementM != nil {
			fmt.Printf("[meas ] %s - %s - %v\n", result.Message, measurementM.Measurement, measurementM.Value)
		} else if traceM := result.GetTrace(); traceM != nil {
			fmt.Printf("[trace] %s - %s\n", result.Message, traceM.ResponseCode)
			// fmt.Printf("traceM.String() = %+v\n", traceM.String())
			if networkT := traceM.GetNetwork(); networkT != nil {
				err = writeNetworkTrace(networkT)
				if err != nil {
					return err
				}
			}

		} else {
			fmt.Println(result.Message)
			fmt.Println(result.String())
		}

		// if traceM := result.GetTrace(); traceM != nil {
		// fmt.Printf("traceM.TraceDataJSON = %+v\n", traceM.TraceDataJSON)
		// }

		if completeM := result.GetComplete(); completeM != nil {
			break
		}
	}

	return
}

func writeNetworkTrace(t *pb.TestResult_Trace_Network) (err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	har := filepath.Join(cwd, "har.json")
	fmt.Println("writing network trace", har)

	f, err := os.Create(har)
	if err != nil {
		return
	}

	_, err = f.WriteString(t.Har)
	if err != nil {
		return
	}

	return
}
