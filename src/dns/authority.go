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
    "net"
)


type Authority ResourceRecord


func (authority *Authority) Unpack(m MessageStream) {
    authority.name = ReadFQName(m)
    binary.Read(io.Reader(m.s), binary.BigEndian, &authority.rr_type)
    binary.Read(io.Reader(m.s), binary.BigEndian, &authority.class)
    binary.Read(io.Reader(m.s), binary.BigEndian, &authority.ttl)
    binary.Read(io.Reader(m.s), binary.BigEndian, &authority.rdlength)
    if authority.rr_type == NS && authority.class == IN {
        authority.rdata = ExpandFQName(m)
    } else {
        authority.rdata = m.s.Next(int(authority.rdlength))
    }
}


func (authority *Authority) String() string {
    result := fmt.Sprintf("%s ",       strings.Join(authority.name, "."))
    result += fmt.Sprintf("type=%s ",  Type2string(authority.rr_type))
    result += fmt.Sprintf("type=%d ",  authority.rr_type)
    result += fmt.Sprintf("ttl=%d ",   authority.ttl)
    result += fmt.Sprintf("class=%s ", ClassMap[authority.class])
    if authority.rr_type == A && authority.class == IN {
        var ip net.IP = authority.rdata
        result += fmt.Sprintf("IP=%s\n", ip)
    } else if authority.rr_type == NS && authority.class == IN {
        result += fmt.Sprintf("domain=%s\n", String(authority.rdata[:]))
    } else if authority.rr_type == SOA && authority.class == IN {
	ms := New(authority.rdata[:])
	mname, rname, serial, refresh, retry, expire, minimum := ReadSOA(*ms)
	result += fmt.Sprintf("soa=%s %s %d %d %d %d %d\n", mname, rname, serial, refresh, retry, expire, minimum)
    } else {
        result += fmt.Sprintf("size=%d ",  authority.rdlength)
        result += fmt.Sprintln(authority.rdata)
    }

    return result
}
