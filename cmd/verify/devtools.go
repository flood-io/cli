package verify

func (b *VerifyCmd) LaunchDevtools() (err error) {
	// chrome, err := exec.LookPath("google-chrome")
	// if err != nil {
	// return
	// }
	// fmt.Printf("chrome = %+v\n", chrome)

	// client, err := b.floodChromeClient()
	// if err != nil {
	// return
	// }
	// defer client.Close()

	// wsEndpoint, err := client.WsEndpoint()
	// if err != nil {
	// return
	// }
	// fmt.Printf("wsEndpoint = %+v\n", wsEndpoint)

	// // chromeTempDir, err := ioutil.TempDir("", "flood-cli")
	// // if err != nil {
	// // return
	// // }
	// // defer os.RemoveAll(chromeTempDir)

	// appArg := fmt.Sprintf("--app='chrome-devtools://devtools/bundled/inspector.html?experiments=true&v8only=true&ws=%s'", wsEndpoint)

	// // cmd := exec.Command(chrome, "--auto-open-devtools-for-tabs", "--user-data-dir", chromeTempDir, "--new-window", url)
	// cmd := exec.Command(chrome, "--profile-directory=FloodChromeDevtools", appArg)

	// err = cmd.Run()
	// if err != nil {
	// return
	// }

	// stdoutStderr, err := cmd.CombinedOutput()
	// if err != nil {
	// return
	// }

	// fmt.Printf("stdoutStderr = %+v\n", stdoutStderr)
	return
}
