# P2P module
P2P module implements the p2p communication of kkchain. 

## Requirements
1. install protoc
brew install protobuf@3.1

2. generate files 
```
go get -u github.com/gogo/protobuf/protoc-gen-gogofaster
    # tested with version v1.1.1 (636bf0302bc95575d69441b25a2603156ffdddf1)
go get -u github.com/golang/mock/mockgen
    # tested with v1.1.1 (c34cdb4725f4c3844d095133c6e40e448b86589b)
go generate ./...
```

## Network
Network implements network stack and wraps all the other modules.

## Host
Host is an object participating in a p2p network, which
implements protocols or provides services. It handles
requests like a Server, and issues requests like a Client.
It is called Host because it is both Server and Client.

## Connection
Connection is a connection to a remote peer. It multiplexes streams.
Usually there is no need to use a Conn directly, but it may
be useful to get information about the peer on the other side:
stream.Conn().RemotePeer()

## Stream
Stream represents a bidirectional channel between two agents in
the kkchain network. "agent" is as granular as desired, potentially
being a "request -> reply" pair, or whole protocols.

## Protocols
Protocols defines the specifications for p2p messages, which are defined 
in different packages according to different purposes.

### Handshake
Handshake messages are defined for handshaking. Ususally, a handshake message will be delivered once connection is established.

### DHT
DHT messages are defined to discover and manage peers. Currently, we implemented a kademlia DHT.

### Chain
Chain messages are defined to fetch or broadcast chain related messages, such as block/tx messages.