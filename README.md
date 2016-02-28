# TeamRadar

A poller to interact with TFS Team Rooms api.  

[https://www.visualstudio.com/integrate/api/chat/overview](https://www.visualstudio.com/integrate/api/chat/overview)


## Building

```
mkdir $HOME/go
export GOPATH=$HOME/go
go get github.com/IslandJohn/TeamRadar/Go/teamradar
cd $HOME/go/src/github.com/IslandJohn/TeamRadar/Go/teamradar
go install
```

## Usage 

```
cd $GOPATH/bin

./teamradar https://vso_account.visualstudio.com/defaultcollection youremail@example.com xnOb.2nqlNDbPgRw9cxnOb.2nqlNDbPgRw9c
```
