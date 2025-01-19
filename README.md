# Speedrunning roadmap.sh Broadcast Server

The following code is a solution to [roadmap.sh](https://roadmap.sh)'s
[Broadcast Server](https://roadmap.sh/projects/broadcast-server) project,
written in Go.

It uses the [Melody](https://github.com/olahol/melody) library for setting up
the WebSocket server, and uses the [Gorilla's WebSocket library](https://github.com/gorilla/websocket)
for the client side.

The CLI is made with [Cobra](https://cobra.dev).

## Building

To build the project, just run

```bash
go build
```

## Usage

### Server

To start the server, use the `start` command:

```bash
srlivechat start
```

By default, the server uses port `3000`. If you want to change that, use the
`--port`, or the `-p`, flag:

```bash
srlivechat start --port 3030
```

To shutdown the server, simply press <kbd>Control</kbd> + <kbd>C</kbd>. This
will also send a shutdown message to all clients.

### Client

To start the client, use the `connect` command:

```bash
srlivechat connect
```

By default, this connects you to `localhost:3000`. If you want to change that,
use the `--host` flag:

```bash
srlivechat start --host localhost:3030
```

You can also give your client a name with the `--username`, or the `-u`, flag:

```bash
srlivechat start --username nilhiu
```

To disconnect and exit, simply press <kbd>Control</kbd> + <kbd>C</kbd>. This
will also send a disconnection message to the server to broadcast.
