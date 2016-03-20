# TeamRadar

TeamRadar is a polyglot project to help collaborate within multiple Team Foundation Server (TFS) and Visual Studio Team Services (VSTS) team rooms simultaenously. The CLI utility polls all available team rooms via REST APIs to generate events. The GUI apps use the CLI to listen to the events and display notifications.  

[https://www.visualstudio.com/integrate/api/chat/overview](https://www.visualstudio.com/integrate/api/chat/overview)

## Command Line Interface (CLI)

The CLI is written in Go and connects to TFS/VSTS, reads commands from stdin, writes events to stdout, and logs exceptions to stderr.

### Building

```
export GOPATH=$HOME/go
mkdir -p $GOPATH
go get github.com/IslandJohn/TeamRadar/Go/teamradar
cd $GOPATH/src/github.com/IslandJohn/TeamRadar/Go/teamradar
go install
```

### Usage 

```
cd $GOPATH/bin
./teamradar https://<account>.visualstudio.com/defaultcollection <email> <token>
./teamradar https://<tfs>/<collection> <user> <password>
```

### Commands

```
join <roomid>
leave <roomid>
send <roomid> <message>
exit|logout|quit
```

## Graphical User Interface (GUI)

### OS X

It's a work in progress.

### Windows

It's planned.
