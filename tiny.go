package main

import (
  "net"
)

func serveClient(conn *net.TCPConn) {
  message := []byte("Yep, serving from ")
  conn.Write(message)
  defer conn.Close()
}

func main() {

  addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8088}
  serverListener, _ := net.ListenTCP("tcp", addr)

  for true {
    conn, _ := serverListener.AcceptTCP()
    go serveClient(conn)
  }
}
