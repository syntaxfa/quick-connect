reference:

https://github.com/gorilla/websocket/tree/main/examples/filewatch

````
$ go get github.com/gorilla/websocket
$ cd `go list -f '{{.Dir}}' github.com/gorilla/websocket/examples/filewatch`
$ go run main.go <name of file to watch>
# Open http://localhost:8080/ .
# Modify the file to see it update in the browser.
````
