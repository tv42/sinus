# Sinus -- command-line remote control for Sonos/UPnP audio devices

## Installation

```console
go get github.com/tv42/sinus
```

## Usage

**Discovery**: First, you must choose what device to use. Run `sinus
discover`. If you have only one UPnP device in your network, it's
automatically chosen; otherwise, you see a list of options. Run `sinus
discover REGEXP..` to prune down the list, until exactly one entry
matches.

After that, you can pause, play, adjust volume and so on:

```
sinus pause
sinus play
sinus vol up
sinus vol down
```

Run just `sinus` or `sinus SUBCOMMAND -help` for help.

## Device support

Currently, Sinus works with UPnP media devices supporting the
[av1](http://upnp.org/specs/av/av1/) spec. It has been mostly tested
with [Sonos](http://www.sonos.com/) devices, and some actions (`sinus
play queue`, `sinus line-in`) have Sonos-specific details.

Better multi-device support may happen. Right now, you can use `sinus
discover` to switch between devices).
