package bludev

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/flood-io/cli/config"
	pb "github.com/flood-io/go-wrenches/floodchrome"
	fcClient "github.com/flood-io/go-wrenches/floodchrome/client"
)

type BLUDev struct {
	LaunchDevtoolsMode bool
	FloodChromeChannel string
}

func (b *BLUDev) floodchromeClient() (client *fcClient.Client, err error) {
	host := "http://localhost:5000"

	token := config.DefaultConfig().APIToken()
	client = fcClient.New(host, token)
	return
}

func (b *BLUDev) Run(scriptFile string) (err error) {
	fmt.Println("running dev-blu")
	fmt.Println("flood chrome channel: ", b.FloodChromeChannel)
	fmt.Println("script file:", scriptFile)

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

	test := &pb.TestRequest{
		Script:             string(scriptBytes),
		FloodChromeVersion: b.FloodChromeChannel,
	}

	currentStep := "<setup>"

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
			if currentStep != measurementM.Label {
				currentStep = measurementM.Label
				fmt.Println("")
				fmt.Println("==================")
				fmt.Println("[ step]", currentStep)
			}
			fmt.Printf("[ meas] %s - %s - %v\n", measurementM.Label, measurementM.Measurement, measurementM.Value)

		} else if traceM := result.GetTrace(); traceM != nil {
			fmt.Printf("[trace] %s - response code %s\n", result.Message, traceM.ResponseCode)
			if networkT := traceM.GetNetwork(); networkT != nil {
				err = writeNetworkTrace(networkT)
				if err != nil {
					return err
				}
			}
		} else if errorM := result.GetError(); errorM != nil {
			fmt.Printf("[error] %s\n", result.Message)
			fmt.Println(errorM.Stack)

		} else if completeM := result.GetComplete(); completeM != nil {
			break

		} else {
			fmt.Println("--- unhandled type ---")
			fmt.Printf("result = %+T\n", result)
			fmt.Println(result.Message)
			fmt.Println(result.String())
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
	fmt.Println("[trace] writing network trace", har)

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
