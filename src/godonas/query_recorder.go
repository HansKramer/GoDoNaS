/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */

package main

import (
    "log"
    "dns"
    "net"
    "fmt"
    "strconv"
//    "bytes"
)


func hex_dump(buffer []byte) string {
    result := ""

    var ascii [16]byte
    for i := range ascii { ascii[i] = '.' }
    for i, b := range buffer[:] {
        result += fmt.Sprintf("%02X ", b)
        if strconv.IsPrint(rune(b)) {
            ascii[i%16] = b
        }
        if i % 16 == 15 {
            result += fmt.Sprintf(" %s\n", ascii)
            for i := range ascii { ascii[i] = '.' }
        }
    }
    if len(buffer) % 16 != 0 {
        for i := 0; i < 16 - (len(buffer) % 16); i++ {
            result += fmt.Sprintf("   ")
        } 
        result += fmt.Sprintf(" %s\n", ascii[:len(buffer)%16])
    }

    return result
}


func first(args ...interface{}) interface{} {
    return args[0]
}


func main() {
    if sock, err := net.ListenUDP("udp", first(net.ResolveUDPAddr("udp", ":53")).(*net.UDPAddr)); err == nil {
        for {
	    var message dns.Message
	    message.Recv(sock)
            fmt.Println(message.String())

	    /*
            var buffer [dns.MAX_MESSAGE_LENGTH]byte
            if rlen, remote, err := sock.ReadFromUDP(buffer[:]); err == nil {
                log.Printf("%s %d", remote, rlen)

                fmt.Println(hex_dump(buffer[0:rlen]))

                var message dns.Message
                message.Unpack(bytes.NewBuffer(buffer[:rlen]))
                fmt.Println(message.String())
	    } else {
                log.Fatal(err)
            }
	    */
        }
    } else {
        log.Fatal(err)
    }
}
