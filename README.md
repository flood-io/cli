# `flood` CLI

The Flood IO Command Line Interface

## Installation

On macOS, install using homebrew:
```bash
brew install flood-io/taps/flood
```

On linux, download the [latest release](https://github.com/flood-io/cli/releases/latest) for your platform, then extract and install it:
```bash
# assuming you're installing version 1.0.0 on linux
cd ~/Downloads
mkdir flood
tar zxvf flood-1.0.0-linux-amd64.zip -C flood

# move the file to somewhere on your $PATH:
mv flood/flood /usr/local/bin/flood

# optionally, tidy up:
rm -rf flood
```

## Writing Flood Chrome scripts

Flood IO now allows you to generate load from Flood Chrome using load scripts written in [TypeScript](https://www.typescriptlang.org) or 
[type-checked Javascript](https://www.typescriptlang.org/docs/handbook/type-checking-javascript-files.html).

The process for running a load test using Flood Chrome is the same as for other generators (such as JMeter and Gatling):
you develop and debug a script locally before creating and running a Flood at scale via the [Flood IO](https://flood.io) web interface.

The `flood verify` command helps to streamline the development and debugging phase by running your script against a 
real-but-unscaled instance of Flood Chrome.

### 1. Create your script
Scripts are written in TypeScript (or Javascript, but we strongly recommend TypeScript to help write your scripts more quickly and robustly).
Full Flood Chrome script documentation is available [here TBA](http://help.flood.io/).

```bash
cd ~/dev
flood init cart-test

# flood init creates some files for you to get started

cd cart-test
```

Now edit `test.ts` to implement your test (*Hint*: if you don't already have a go-to code editor, we recommend [Microsoft Visual Studio Code](https://code.visualstudio.com/) as it works
particularly well with TypeScript as used by Flood Chrome.)

### 2. Run the script
```bash
flood verify test.ts
```

**Note** that you must first log in using `flood login` and your Flood IO credentials ([see below](#authenticating))

### 3. Iterate - develop and debug your script

Verify that your script is compiling and running correctly, and is testing the correct things:

- Observe errors in the `flood verify` output.
- Use `console.log` and `assert` to ensure your test is working correctly.
- Set a unique value for `settings.userAgent` and search for it in your logs.

Each time fix or tweak your script then, return to step 2 and re-run it.

### 4. Launch a scaled up Flood at Flood IO

Once you're happy with the operation of your test script, visit https://flood.io to create a Flood.
Here you can run your script on up to 50 Flood Chrome instances per Flood node for simulating thousands of browser users.

(**Note** we plan to expand this CLI to allow running fully scaled Floods directly from the command line)

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
