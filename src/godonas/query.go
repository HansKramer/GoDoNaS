/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */

package main


import (
    "dns"
    "os"
//    "fmt"
)


func uncompress(buf []byte) []byte {
/*
    r := bytes.NewBuffer(buf[:])

    var header Header
    header.Unpack(r)

    for (i:=0 ; i:=int(header.qdcount); i++) {
       hk  
    }
fmt.Println(header.qdcount)
fmt.Println(header.ancount)
fmt.Println(header.nscount)
fmt.Println(header.arcount)
*/

    new_buf := buf[:12]
    return new_buf
}


func main() {
    if len(os.Args) ==2 {
        var message dns.Message

        message.Query(os.Args[1], dns.A)
        answer := message.Send("75.75.75.75:53")

        answer.Print()
    }
}

