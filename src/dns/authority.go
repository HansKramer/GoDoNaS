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
)


type Authority ResourceRecord


func (authority *Authority) Unpack(s *bytes.Buffer, r []byte) {
    authority.name  = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.rr_type)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.class)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.rdlength)
    authority.rdata = s.Next(int(authority.rdlength))
}


