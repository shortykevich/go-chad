# Little chat app

Simple CLI for creating chat server/connecting to already active server.

# Usage

Make sure you have [Go](https://go.dev/dl/) installed

```
git clone git@github.com:shortykevich/go-chad.git
cd go-chad
make build-exec
```

Now _app_ executable should be created and you can see all possible configurations with command:

```
./app -h
```

To run on/connect to localhost:8554:

```
# To create server
./app -mode=server
# To connect as client
./app -mode=client
```

### Origin

Project idea inspired [roadmap.sh](https://roadmap.sh/)
