/*

    Author : Hans Kramer

      Date : July 2016

      Code :

 */

package dns


import (
    "encoding/binary"
    "io"
    "bytes"
)


type Header struct {
    id      uint16
    status  uint16
    Qdcount uint16
    Ancount uint16
    Nscount uint16
    Arcount uint16
}


func (header *Header) Init() {
    header.id      = 0
    header.status  = 0
    header.Qdcount = 0
    header.Ancount = 0
    header.Nscount = 0
    header.Arcount = 0
}


func (header *Header) Unpack(r io.Reader) {
    binary.Read(r, binary.BigEndian, &header.id)
    binary.Read(r, binary.BigEndian, &header.status)
    binary.Read(r, binary.BigEndian, &header.Qdcount)
    binary.Read(r, binary.BigEndian, &header.Ancount)
    binary.Read(r, binary.BigEndian, &header.Nscount)
    binary.Read(r, binary.BigEndian, &header.Arcount)
}


func (header Header) Pack() []byte {
    buf := new(bytes.Buffer)

    binary.Write(buf, binary.BigEndian, header.id)
    binary.Write(buf, binary.BigEndian, header.status)
    binary.Write(buf, binary.BigEndian, header.Qdcount)
    binary.Write(buf, binary.BigEndian, header.Ancount)
    binary.Write(buf, binary.BigEndian, header.Nscount)
    binary.Write(buf, binary.BigEndian, header.Arcount)

    return buf.Bytes()
}


func (header Header) GetField(field int) uint16 {
    switch field {
    case ID:
        return header.id
    case QR:
        return (header.status & 0x8000) >> 15
    // to be implemented
    }
    return 0
}


func (header *Header) SetField(field int, value uint16) {
    switch field {
    case ID:
        header.id = value
    case RCODE:
        header.status |= (value & 0x0f)
    case RA:
        header.status |= (value & 0x01) << 7
    case RD:
        header.status |= (value & 0x01) << 8
    case TC:
        header.status |= (value & 0x01) << 9
    case AA:
        header.status |= (value & 0x01) << 10
    case OPCODE:
        header.status |= (value & 0x0f) << 14
    case QR:
        header.status |= (value & 0x01) << 15
    }
}

