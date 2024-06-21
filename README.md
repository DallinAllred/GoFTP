# Introduction

This is simple FTP-like server and client created as a practice exercise in Go.

# Usage
## Server
```sh
go run server/ftpServer.go -p <port>
# Default port: 2020
```

## Client
```sh
go run server/ftpClient.go -h <host> -p <port>
# Default host: localhost
# Default port: 2020
```
### Client Commands
#### Local
```sh
lcd directory # Change local working directory
lls <directory> # List contents of local working directory or, optionally, specify an alternate directory
lpwd # Print the local working directory
```

#### Remote
```sh
cd directory # Change remote working directory
ls <directory> # List contents of remote working directory or, optionally, specify an alternate directory
lpwd # Print the remote working directory
exit # Terminates client
```

#### File Transfer*

```sh
get filename # Copies a file from the server to the client
put filename # Sends a file from the client to the server
```
<span style="color: orange">* Functionality Pending</span>