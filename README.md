# Little chat app

Simple CLI for creating chat server/connecting to already active server.

# Usage

Make sure you have [Go](https://go.dev/dl/) installed

```bash
git clone git@github.com:shortykevich/go-chad.git
cd go-chad
make build-exec
```

Now _app_ executable should be created and you can see all possible configurations with command:

```bash
./app -h
```

To run on/connect to localhost:8554:

```bash
# To create server
./app -mode=server
# To connect as client
./app -mode=client
```

### Origin

Idea of application inspired by list of backend projects on [roadmap.sh](https://roadmap.sh/)
