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
    "fmt"
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


func (header Header) String() string {
    result := fmt.Sprintf("Id: %d\n",      header.GetField(ID))
    result += fmt.Sprintf("Status : %d\n", header.GetField(STATUS))
    result += fmt.Sprintf("    QR : %d\n", header.GetField(QR))
    result += fmt.Sprintf("Opcode : %d\n", header.GetField(OPCODE))
    result += fmt.Sprintf("    AA : %d\n", header.GetField(AA))
    result += fmt.Sprintf("    TC : %d\n", header.GetField(TC))
    result += fmt.Sprintf("    RD : %d\n", header.GetField(RD))
    result += fmt.Sprintf("    RA : %d\n", header.GetField(RA))
    result += fmt.Sprintf("     Z : %d\n", header.GetField(Z))
    result += fmt.Sprintf(" RCODE : %d\n", header.GetField(RCODE))
    result += fmt.Sprintf("QD Count : %d\n", header.Qdcount)
    result += fmt.Sprintf("AN Count : %d\n", header.Ancount)
    result += fmt.Sprintf("NS Count : %d\n", header.Nscount)
    result += fmt.Sprintf("AR Count : %d\n", header.Arcount)

    return result
}


func (header Header) GetField(field int) uint16 {
    switch field {
    case ID:
        return header.id
    case STATUS:
        return header.status
    case QR:
        return (header.status >> 15) & 0x01
    case OPCODE:
        return (header.status >> 11) & 0x0f
    case AA:
	return (header.status >> 10) & 0x01
    case TC:
	return (header.status >>  9) & 0x01
    case RD:
	return (header.status >>  8) & 0x01
    case RA:
	return (header.status >>  7) & 0x01
    case Z:
        return (header.status >>  4) & 0x07
    case RCODE:
	return header.status & 0x0f
    }
    return 0
}


func (header *Header) SetField(field int, value uint16) {
    switch field {
    case ID:
        header.id = value
    case STATUS:
        header.status = value
    case QR:
        header.status |= (value & 0x01) << 15
    case OPCODE:
        header.status |= (value & 0x0f) << 11
    case AA:
        header.status |= (value & 0x01) << 10
    case TC:
        header.status |= (value & 0x01) << 9
    case RD:
        header.status |= (value & 0x01) << 8
    case RA:
        header.status |= (value & 0x01) << 7
    case Z:
        header.status |= (value & 0x07) << 4
    case RCODE:
        header.status |= (value & 0x0f)
    }
}
