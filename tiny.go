package main

import (
  "os"
  "net"
  "net/textproto"
  "bufio"
  "strings"
  /* "fmt"*/
)

type Request struct {
  // request method (i.e GET, POST ...)
  Method string

  // path to the resource
  Path string

  // Protocl version
  HTTPVersion string
}

/*
 * Here starts our program
*/
func main() {
  addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8088}
  serverListener, _ := net.ListenTCP("tcp", addr)

  // listens for incomming connections continuously
  for {
    conn, _ := serverListener.AcceptTCP()
    // sets one go routine per connection established
    go serveClient(conn)
  }
}

/*
 * This methods handles client request
 * 1. Parses the request
 * 2. Decides on validity of request
 * 3. Decides on the action to be taken
*/
func serveClient(conn *net.TCPConn) {
  defer conn.Close() // we will close the connection once we exit this function

  // making buffered IO reader and simple echo server
  var reader = textproto.NewReader(bufio.NewReader(conn))

  // making new request object
  request, err := parseRequest(reader)
  
  if err != nil {
    conn.Write([]byte("Invalid request"))
    return
  }

  // writing back to client
  conn.Write([]byte(request.Method))
  conn.Write([]byte(request.Path))
  conn.Write([]byte(request.HTTPVersion))
}

/*
 * Parse the client request headers and build request object.
*/
func parseRequest(r *textproto.Reader) (*Request, os.Error) {

  // create new request object
  request := new(Request)

  methodLine, _ := r.ReadLine()
  methodLineElements := strings.Split(methodLine, " ")
  
  if (len(methodLineElements) != 3) {
    return request, os.NewError("Invalid request")
  }

  request.Method = methodLineElements[0]
  request.Path = methodLineElements[1]
  request.HTTPVersion = methodLineElements[2]
  return request, nil
}
