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


type MessageStream struct {
    s *bytes.Buffer      // stream pointer
    r []byte             // random access (required due to DNS compression)
}


type ResourceRecord struct {
    name     []string
    rr_type  uint16
    class    uint16
    ttl      uint32
    rdlength uint16
    rdata    []byte
}


func ReadFQName(m MessageStream) []string {
    var oct_len uint8
    var data    []string
    for {
        binary.Read(io.Reader(m.s), binary.BigEndian, &oct_len)
        if oct_len == 0 {
            return data
        }
        if (oct_len & 0xc0) == 0xc0 {  // it's a pointer
            offset := int16(oct_len & 0x3f) * 256
            binary.Read(io.Reader(m.s), binary.BigEndian, &oct_len)
            offset += int16(oct_len)

            return append(data, ReadFQName(MessageStream{bytes.NewBuffer(m.r[offset:]), m.r})...)
        }
        data = append(data, string(m.s.Next(int(oct_len))))
    }
}


func WriteFQName(name []string) []byte {
    buf := new(bytes.Buffer)

    for _, label := range name {
        var length byte = byte(len(label))
        binary.Write(buf, binary.BigEndian, length)
        buf.WriteString(label)
    }
    binary.Write(buf, binary.BigEndian, byte(0))

    return buf.Bytes()
}


func ExpandFQName(m MessageStream) []byte {
    return WriteFQName(ReadFQName(m))
}


func String(r []byte) (value string) {
    for i := 0; i < len(r) ; i++ {
	c := int(r[i])
	if c > 0 {
	    value += string(r[i+1:i+c+1]) + "."
            i += int(c)
        }
    }
    return value[:len(value)-1]
}



// move out 
/*

type Additional ResourceRecord

func (additional *Additional) Unpack(s *bytes.Buffer, r []byte) {
    additional.name  = ReadFQName(MessageStream{s, r})
    binary.Read(io.Reader(s), binary.BigEndian, &additional.rr_type)
    binary.Read(io.Reader(s), binary.BigEndian, &additional.class)
    binary.Read(io.Reader(s), binary.BigEndian, &additional.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &additional.rdlength)
    additional.rdata = s.Next(int(additional.rdlength))
}
*/
