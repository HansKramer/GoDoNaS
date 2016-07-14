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
    "fmt"
)

type Question struct {
    qname  []string
    qtype  uint16
    qclass uint16
}


func (question Question) Pack() []byte {
    buf := new(bytes.Buffer)

    buf.Write(WriteFQName(question.qname))
    binary.Write(buf, binary.BigEndian, question.qtype)
    binary.Write(buf, binary.BigEndian, question.qclass)

    return buf.Bytes()
}


func (question *Question) Unpack(s *bytes.Buffer, r []byte) {
    question.qname = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &question.qtype)
    binary.Read(io.Reader(s), binary.BigEndian, &question.qclass)
}


func (question *Question) String() string {
    result := fmt.Sprintln(question.qname)
    result += fmt.Sprintln(question.qtype)
    result += fmt.Sprintln(question.qclass)

    return result
} 
