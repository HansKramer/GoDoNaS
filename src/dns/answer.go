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
    "strings"
    "net"
    "bytes"
)

import "fmt"

type Answer ResourceRecord


func (answer *Answer) Unpack(m MessageStream) {
    answer.name = ReadFQName(m)
    binary.Read(io.Reader(m.s), binary.BigEndian, &answer.rr_type)
    binary.Read(io.Reader(m.s), binary.BigEndian, &answer.class)
    binary.Read(io.Reader(m.s), binary.BigEndian, &answer.ttl)
    binary.Read(io.Reader(m.s), binary.BigEndian, &answer.rdlength)
    if answer.rr_type == PTR && answer.class == IN {
        answer.rdata = ExpandFQName(m)
    } else {
        answer.rdata = m.s.Next(int(answer.rdlength))
    }
}


func (answer *Answer) Pack() []byte {
    buf := new(bytes.Buffer)

    return buf.Bytes()
}


func (answer Answer) Get() (net.IP) {
    return answer.rdata
}


func (answer *Answer) String() string {
    result := fmt.Sprintf("%s ",       strings.Join(answer.name, "."))
    result += fmt.Sprintf("type=%s ",  Type2string(answer.rr_type))
    result += fmt.Sprintf("ttl=%d ",   answer.ttl)
    result += fmt.Sprintf("class=%s ", ClassMap[answer.class])
    result += fmt.Sprintf("size=%d ",  answer.rdlength)
    if answer.rr_type == A && answer.class == IN {
        var ip net.IP = answer.rdata
        result += fmt.Sprintf("IP=%s\n", ip)
    } else if answer.rr_type == PTR && answer.class == IN {
        result += fmt.Sprintf("domain=%s\n", String(answer.rdata[:]))
    } else {
        result += fmt.Sprintln(answer.rdata)
    }

    return result
}
