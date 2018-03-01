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
    "os"
    "stringutil"
    "strings"
    "fmt"
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


func (rs *ResourceRecord) SetName(name string) {
    rs.name = strings.Split(name, ".")
}


func New(buf []byte) *MessageStream {
    return &MessageStream{
	s : bytes.NewBuffer(buf),
	r : buf,
    }
}


func (m MessageStream) WriteToFile(filename string) (err error) {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil { return }
    defer f.Close()

    _, err = f.Write(m.r)
    return 
}


func (m MessageStream) WriteToStdout() {
    os.Stdout.Write(m.r)
}


func (m MessageStream) Hexdump() string {
    return stringutil.Hexdump(m.r)
}


func ReadSOA(m MessageStream) (mname string, rname string, serial int32, refresh int32, retry int32, expire int32, minimum int32) {
    mname   = String(ExpandFQName(m))
    rname   = String(ExpandFQName(m))
    binary.Read(io.Reader(m.s), binary.BigEndian, &serial)
    binary.Read(io.Reader(m.s), binary.BigEndian, &refresh)
    binary.Read(io.Reader(m.s), binary.BigEndian, &retry)
    binary.Read(io.Reader(m.s), binary.BigEndian, &expire)
    binary.Read(io.Reader(m.s), binary.BigEndian, &minimum)

    return
}


func ReadFQNameDebug(m MessageStream) []string {
    var oct_len uint8
    var data    []string
    for {
        binary.Read(io.Reader(m.s), binary.BigEndian, &oct_len)
 fmt.Println("oct_len")
 fmt.Println(oct_len)
        if oct_len == 0 {
 fmt.Println(data)
 fmt.Println("done")
            return data
        }
        if (oct_len & 0xc0) == 0xc0 {  // it's a pointer
            offset := int16(oct_len & 0x3f) * 256
            binary.Read(io.Reader(m.s), binary.BigEndian, &oct_len)
            offset += int16(oct_len)
 fmt.Println("pointer", offset)
 fmt.Println(m.r[offset-1:])
 fmt.Println(offset)
 fmt.Println(data)
 fmt.Println("done")
            return append(data, ReadFQNameDebug(MessageStream{bytes.NewBuffer(m.r[offset:]), m.r})...)
        }
        data = append(data, string(m.s.Next(int(oct_len))))
    }
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


func ExpandFQNameDebug(m MessageStream) []byte {
    return WriteFQName(ReadFQNameDebug(m))
}


func ExpandFQName(m MessageStream) []byte {
    return WriteFQName(ReadFQName(m))
}


func String2(r []string) (value string) {
    return strings.Join(r, ".")
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
