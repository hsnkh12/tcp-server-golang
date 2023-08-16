# tcp-server-golang
Basic TCP server with Go language


## Donwload 
- Clone the repository to your local machine:

```shell

git clone https://github.com/hsnkh12/tcp-server-golang/
```
- To run the app:
```shell

go run main.go
```
- To compile the app:
```shell 

go build
```

## Server methods explanation
Here are some highlights of your code:

- CreateTCPServer Function: Creates a new Server instance with all necessary channels and maps.

- Listen Method: Sets up the server to listen on the provided address. It handles signal notifications for graceful shutdown and starts accepting connections in a separate goroutine.

- AcceptConnections Method: Accepts incoming connections and tracks them in the ActiveConnections map, while also starting a goroutine to handle each connection.

- ReadConneciton Method: Reads data from the connection, processes received messages, and removes the connection from the ActiveConnections map when the connection is closed.

- CloseAllConnections Method: Closes all active connections gracefully during server shutdown.


