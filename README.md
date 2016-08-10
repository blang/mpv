mpv (remote) control library
======

This library provides everything needed to (remote) control the [mpv media player](https://mpv.io/).

It provides an easy api, a json api and rpc functionality.

Usecases: Remote control your mediaplayer running on a raspberry pi or laptop or build a http interface for mpv

Usage
-----

```bash
$ go get github.com/blang/mpv
```
Note: Always vendor your dependencies or fix on a specific version tag.

Start mpv:
```bash
$ mpv --idle --input-ipc-server=/tmp/mpvsocket
```

Remote control:
```go
import github.com/blang/mpv

ipcc := mpv.NewIPCClient("/tmp/mpvsocket") // Lowlevel client
c := mpv.NewClient(ipcc) // Highlevel client, can also use RPCClient

c.LoadFile("movie.mp4", mpv.LoadFileModeReplace)
c.SetPause(true)
c.Seek(600, mpv.SeekModeAbsolute)
c.SetFullscreen(true)
c.SetPause(false)

pos, err := c.Position()
fmt.Printf("Position in Seconds: %.0f", pos)
```

Also check the [GoDocs](http://godoc.org/github.com/blang/mpv).


Features
-----

- Low-Level and High-Level API
- RPC Server and Client (fully transparent)
- HTTP Handler exposing lowlevel API (json)


Contribution
-----

Feel free to make a pull request. For bigger changes create a issue first to discuss about it.


License
-----

See [LICENSE](LICENSE) file.
