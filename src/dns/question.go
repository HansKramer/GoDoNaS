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
    "strings"
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


func (question *Question) Unpack(m MessageStream) {
    question.qname = ReadFQName(m)
    binary.Read(io.Reader(m.s), binary.BigEndian, &question.qtype)
    binary.Read(io.Reader(m.s), binary.BigEndian, &question.qclass)
}


func (question *Question) String() string {
    return fmt.Sprintf("question: %s type=%s class=%s",
                       strings.Join(question.qname, "."),
		       Type2string(question.qtype),
		       ClassMap[question.qclass])
}
