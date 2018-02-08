# `flood` CLI

The Flood IO Command Line Interface

## Installation

Download the [latest release](https://github.com/flood-io/cli/releases/latest) for your platform, then extract and install it:
```bash
# assuming you're installing version 0.2.0
cd ~/Downloads
unzip -d flood flood-0.2.0-darwin-amd64.zip

# move the file to somewhere on your $PATH:
mv flood/flood /usr/local/bin/flood

# optionally, tidy up:
rm -rf flood
```

## Writing Flood Chrome scripts

Flood IO now allows you to generate load from Flood Chrome using load scripts written in TypeScript or Javascript.

The process for running a load test using Flood Chrome is the same as for other generators (such as JMeter and Gatling):
you develop and debug a script locally before creating and running a Flood at scale via the [Flood IO](https://flood.io) web interface.

The `flood verify` command helps to streamline the development and debugging phase by running your script against a 
real-but-unscaled instance of Flood Chrome.

### 1. Create your script
Scripts are written in TypeScript (or Javascript, but we strongly recommend TypeScript to help write your scripts more quickly and robustly).
Full Flood Chrome script documentation is available [here TBA](http://help.flood.io/).

```javascript
// cart-test.ts
import { step, TestSettings, Until, By, MouseButtons, Device, Driver } from '@flood/chrome'
import * as assert from 'assert'
export const settings: TestSettings = {
  userAgent: 'my-flood-io-cart-test/1.0',
	disableCache: true,
	actionDelay: 0.5,
	stepDelay: 2.5,
}

/**
 * My Cart Script
 * This script tests out the shopping cart flow
 * for our site.
 */
export default () => {
	step('My Cart', async (browser: Driver) => {
		await browser.visit('https://challenge.flood.io')
    ...
  })
  step('My Cart: item added', async (browser: Driver) => {
  })
}
```

### 2. Run the script
```bash
flood verify cart-test.ts
```

**Note** that you must first log in using `flood login` and your Flood IO credentials ([see below](#authenticating))

### 3. Iterate - develop and debug your script

Verify that your script is compiling and running correctly, and is testing the correct things:

- Observe errors in the `flood verify` output.
- Use `console.log` and `assert` to ensure your test is working correctly.
- Set a unique value for `settings.userAgent` and search for it in your logs.

Each time fix or tweak your script then, return to step 2 and re-run it.

### 4. Launch a scaled up Flood at [Flood IO](https://flood.io]

Once you're happy with the operation of your test script, visit https://flood.io to create a Flood.
Here, you can run your script on up to 50 Flood Chrome instances per Flood node.

(**Note** we have plans to expand this CLI to allow running fully scaled Floods directly from the command line)

## Usage

```bash
flood help
```

### Authenticating

Before you can do anything, we need to authenticate your local machine with your
Flood account:

    $ flood login

    What's your username:
    Enter a value: user@exampe.com

    What's your password (masked):
    Enter a value: **************************************************

    Welcome back Ivan Vanderbyl!

This will store a temporary authentication token on your machine.

You can deregister this machine by running `flood logout` at any time.

## Development

See [Development](DEVELOPMENT.md)
