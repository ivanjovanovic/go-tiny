/*
 * Simple, dummy HTTP server that just can serve a given path if 
 * present in the current directory. 
 */
package main

import (
  "os"
  "net"
  "net/textproto"
  "bufio"
  "strings"
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

  // making new request object
  request, err := parseRequest(bufio.NewReader(conn))
  
  if err != nil {
    conn.Write([]byte("HTTP/1.1 400 Bad Request"))
    return
  }

  // return error if path contains ../
  if strings.Contains(request.Path, "../") {
    conn.Write([]byte("HTTP/1.1 403 Forbidden"))
    return
  }

  dir, _ := os.Getwd()
  filePath := strings.Join([]string{dir, request.Path}, "")

  /* Checking if file exists
   * We do type assertion (announcing that error is of type *PathError)
   * in order to get the proper type
   * @see http://golang.org/doc/go_spec.html#Type_assertions
   */
  if _, err := os.Stat(filePath); err != nil && err.(*os.PathError).Error == os.ENOENT {
    conn.Write([]byte("HTTP/1.1 404 Not Found"))
    return
  }

  // now when I have request for a certain path I can serve
  file, err := os.Open(filePath)

  if err != nil {
    conn.Write([]byte("HTTP/1.1 500 Internal Server Error"))
    return
  }

  // writing file to connection by transfering chunk by chunk
  var buffer = make([]byte, 1024)
  for {
    n, err := file.Read(buffer)
    conn.Write(buffer[0:n])
    if err == os.EOF {
      break
    }
  }
}

/*
 * Parse the client request headers and build request object.
*/
func parseRequest(reader *bufio.Reader) (*Request, os.Error) {

  var r = textproto.NewReader(reader)

  // create new request object
  request := new(Request)

  methodLine, _ := r.ReadLine()
  methodLineElements := strings.Split(methodLine, " ")
  
  if (len(methodLineElements) != 3) {
    return request, os.NewError("Invalid request")
  }

  request.Method = methodLineElements[0]

  if methodLineElements[1] == "/" {
    request.Path = "index.html"
  } else {
    request.Path = methodLineElements[1]
  }
  request.HTTPVersion = methodLineElements[2]
  return request, nil
}
