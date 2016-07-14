/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */

package dns

import (
    "io"
    "encoding/binary"
    "bytes"
    "strings"
    "net"
)

import "fmt"

type Answer ResourceRecord


func (answer *Answer) Unpack(s *bytes.Buffer, r []byte) {
    answer.name  = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.rr_type)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.class)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.rdlength)
    answer.rdata = s.Next(int(answer.rdlength))
}


func (answer *Answer) String() string {
    result := fmt.Sprintln(strings.Join(answer.name, "."))
    result += fmt.Sprintln(Type2string(answer.rr_type))
    result += fmt.Sprintln(answer.ttl)
    result += fmt.Sprintln(answer.rdlength)
    if answer.rr_type == A {
        var ip net.IP = answer.rdata
        result += fmt.Sprintln(ip)
    } else {
        result += fmt.Sprintln(answer.rdata)
    }

    return result
}
