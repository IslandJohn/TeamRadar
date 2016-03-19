# TeamRadar

A CLI utility to interact with TFS team rooms REST API, and some GUI tools for notifications.  

[https://www.visualstudio.com/integrate/api/chat/overview](https://www.visualstudio.com/integrate/api/chat/overview)

## Building

```
export GOPATH=$HOME/go
mkdir -p $GOPATH
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
