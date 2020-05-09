# `huectl`

`huectl` is a simple CLI tool to manage a [Philips Hue](https://meethue.com) installation. Once connected to a Hue bridge, it can control lights from a terminal.

# Installation

## With Go Installed

Assuming `$GOPATH/bin` is already on your path, you can run:

```
go get github.com/skwair/huectl
```

## Binary Releases

Binary releases are available [here](https://github.com/skwair/huectl/releases) for Linux, MacOS and Windows. Download the one matching your system and move it somewhere in your system path.

# Connecting to a Hue Bridge

`huectl` can automatically register a new user on a bridge for you using the `init` command. All you need to do is press the Link Button on the Bridge when asked. This command only needs to be run once.

```
$> huectl init
Searching for a Hue bridge on your local network...
Found Hue bridge "Philips hue" at: 192.168.1.50
Registering new user, please press the button on the bridge then press `Enter`

Saving configuration to "/home/user/.config/huectl/config.yml"
```

All requests to the bridge are using HTTPS, but Philips only provides self-signed certificates, so for additionnal security and when making the first connecting to the bridge, `huectl` will save its certificate fingerprint and will check that is has not changed when running other commands.

# CLI Examples

To list available lights:

```
$> huectl lights list
ID    NAME          ON       REACHABLE    BRIGHTNESS (%)    HUE
1     Kitchen       false    true         83                7170
2     Leaving room  false    true         35                7170
3     Bedroom       false    true         100               8418
```

To set the state of a light:

```
$> huectl light set 1 --on --bri=75
```

To simply toggle a light:

```
$> huectl light toggle 1
```

# License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/skwair/huectl/blob/master/LICENSE) file for details.
