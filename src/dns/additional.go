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
    "fmt"
    "strings"
)


type Additional ResourceRecord


func (additional *Additional) Unpack(m MessageStream) {
    additional.name = ReadFQName(m)
    binary.Read(io.Reader(m.s), binary.BigEndian, &additional.rr_type)
    binary.Read(io.Reader(m.s), binary.BigEndian, &additional.class)
    binary.Read(io.Reader(m.s), binary.BigEndian, &additional.ttl)
    binary.Read(io.Reader(m.s), binary.BigEndian, &additional.rdlength)
    additional.rdata = m.s.Next(int(additional.rdlength))
}


func (additional *Additional) String() string {
    result := fmt.Sprintf("%s ",        strings.Join(additional.name, "."))
    result += fmt.Sprintf("type=%s ",   Type2string(additional.rr_type))
    result += fmt.Sprintf("ttl=%d ",    additional.ttl)
    result += fmt.Sprintf("class=%s ",  ClassMap[additional.class])
    result += fmt.Sprintf("size=%d \n", additional.rdlength)

    return result
}
