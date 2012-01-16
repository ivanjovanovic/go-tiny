package main

import (
  "os"
  "net"
  "bufio"
)

type Request struct {
  // request method (i.e GET, POST ...)
  method []byte

  // host
  host []byte

  // path to the resource
  path []byte
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
 * One go routing per client connection
*/
func serveClient(conn *net.TCPConn) {
  defer conn.Close() // we will close the connection once we exit this function

  // making buffered IO reader and simple echo server
  var reader = bufio.NewReader(conn)

  // making new request object
  request, _ := parseRequest(reader)

  // writing back to client
  conn.Write(request.method)
}

/*
 * Parse the client request headers and build request object.
*/
func parseRequest(r *bufio.Reader) (request *Request, err os.Error) {
  return new(Request), nil
}


