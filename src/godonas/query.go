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
 //   "bytes"
    "fmt"
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

        //message.Query(os.Args[1], dns.PTR)
        message.Query(os.Args[1], dns.A)
        answer := message.Send("10.20.6.22:53")
        //answer := message.Send("75.75.75.75:53")
	fmt.Println("Answer: -------")
        fmt.Println(answer.String())
        fmt.Println("---------------")

        //message.Query("4.4.8.8.in-addr.arpa", dns.A)
        //answer = message.Send("75.75.75.75:53")
        //answer.Print()
    }

/*
    b := []byte{ 0xc1,  0xec,  0x01,  0x00,  0x00,  0x01,  0x00,  0x00,  0x00,  0x00,  0x00,  0x00,  0x02,  0x36,  0x38,  0x01,
    0x30,  0x01,  0x30,  0x02,  0x31,  0x30,  0x07,  0x69,  0x6e,  0x2d,  0x61,  0x64,  0x64,  0x72,  0x04,  0x61,
    0x72,  0x70,  0x61,  0x00,  0x00,  0x0c,  0x00,  0x01       }
    
    var m dns.Message;
    m.Unpack(bytes.NewBuffer(b[:]));
    m.Print()
*/
}

