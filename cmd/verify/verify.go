package verify

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	pb "github.com/flood-io/go-wrenches/floodchrome"
	fcClient "github.com/flood-io/go-wrenches/floodchrome/client"
	"github.com/pkg/errors"
)

type VerifyCmd struct {
	LaunchDevtoolsMode bool
	FloodChromeChannel string

	Host    string
	DevMode string
	Verbose bool
}

func (b *VerifyCmd) floodchromeClient(token string) (client *fcClient.Client, err error) {
	client = fcClient.New(b.Host, token)
	return
}

func (b *VerifyCmd) Run(authToken string, scriptFile string) (err error) {
	var test Test = NewLoggingTest(&StatefulTest{})

	// test := (Test)(nil)
	// &StatefulTest{
	// ScriptFile: scriptFile,
	// Channel:    b.FloodChromeChannel,
	// }

	// ui.SetStatus("Flood Chrome Verify")

	// ui.Log("flood chrome channel: ", b.FloodChromeChannel)
	// ui.Log("script file:", scriptFile)

	test.AssertConfigured()

	f, err := os.Open(scriptFile)
	if err != nil {
		err = errors.Wrapf(err, "unable to open script at %s", scriptFile)
		return
	}

	scriptBytes, err := ioutil.ReadAll(f)
	if err != nil {
		err = errors.Wrap(err, "unable to read script contents")
		return
	}

	client, err := b.floodchromeClient(authToken)
	if err != nil {
		err = errors.Wrap(err, "unable to init flood-chrome client")
		return
	}

	testRequest := &pb.TestRequest{
		Script:             string(scriptBytes),
		ScriptFilename:     scriptFile,
		FloodChromeVersion: b.FloodChromeChannel,
	}

	state := &state{
		Test: test,
	}
	state.state = state.awaitTest

	stream, err := client.Run(context.Background(), testRequest)
	if err != nil {
		err = errors.Wrap(err, "unable to call test.Run")
		return
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "error on test result stream")
		}

		// state.dumpLife(result)

		// fmt.Printf("result = %+v\n", result)
		// fmt.Printf("result = %T\n", result)
		if b.Verbose {
			fmt.Printf("message (%T) %s\n", result.Result, result.String())
		}

		err = state.next(result)
		if err != nil {
			return errors.Wrap(err, "error handling test result message")
		}
	}

	return
}
