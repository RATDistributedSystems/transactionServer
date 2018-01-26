package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {
    for{
  // connect to this socket

  conn, _ := net.Dial("tcp", "localhost:3334")
  //  conn, _ := net.Dial("tcp", "134.87.152.32:44441")

    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Text to send: ")
    text, _ := reader.ReadString('\n')
    // send to socket
    fmt.Fprintf(conn, text + "\n")
    // listen for reply
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Message from server: "+message)
    }
}