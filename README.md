# TeamRadar

A poller to interact with TFS Team Rooms api.  

[https://www.visualstudio.com/integrate/api/chat/overview](https://www.visualstudio.com/integrate/api/chat/overview)

## Building

```
export GOPATH=$HOME/go
mkdir -p $HOME/go
go get github.com/IslandJohn/TeamRadar/Go/teamradar
cd $GOPATH/src/github.com/IslandJohn/TeamRadar/Go/teamradar
go install
```

## Usage 

```
cd $GOPATH/bin
./teamradar https://<account>.visualstudio.com/defaultcollection <email> <token>
./teamradar https://<tfs>/<collection> <user> <password>
```
