/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */
package main

import (
    "fmt"
    "net"
    "encoding/binary"
    "bytes"
    "io"
    "strings"
    "math/rand"
)

import "dns"

type Header struct {
    id      uint16
    status  uint16
    qdcount uint16
    ancount uint16
    nscount uint16
    arcount uint16
}

type Question struct {
    qname  []string
    qtype  uint16
    qclass uint16 
}

type ResourceRecord struct {
    name     []string
    rr_type  uint16
    class    uint16 
    ttl      uint32
    rdlength uint16
    rdata    []byte
}

type Answer ResourceRecord

type Authority ResourceRecord

type Additional ResourceRecord

type Message struct {
    header     Header
    question   []Question  // allow for multiple questions even though most DNS servers don't support this
    answer     []Answer  
    authority  []Authority
    additional []Additional
}

func (header *Header) Unpack(r io.Reader) {
    binary.Read(r, binary.BigEndian, &header.id)
    binary.Read(r, binary.BigEndian, &header.status)
    binary.Read(r, binary.BigEndian, &header.qdcount)
    binary.Read(r, binary.BigEndian, &header.ancount)
    binary.Read(r, binary.BigEndian, &header.nscount)
    binary.Read(r, binary.BigEndian, &header.arcount)
}

func (header Header) Pack() []byte {
    buf := new(bytes.Buffer)

    binary.Write(buf, binary.BigEndian, header.id)
    binary.Write(buf, binary.BigEndian, header.status)
    binary.Write(buf, binary.BigEndian, header.qdcount)
    binary.Write(buf, binary.BigEndian, header.ancount)
    binary.Write(buf, binary.BigEndian, header.nscount)
    binary.Write(buf, binary.BigEndian, header.arcount)
 
    return buf.Bytes()
}

func (header *Header) Init() {
    header.id      = 0
    header.status  = 0
    header.qdcount = 0
    header.ancount = 0
    header.nscount = 0
    header.arcount = 0
}

func (header Header) GetField(field int) uint16 {
    switch field {
    case dns.ID:
        return header.id
    case dns.QR:
        return (header.status & 0x8000) >> 15
    // to be implemented 
    }
    return 0
}

func (header *Header) SetCount(field int, value uint16) {
    switch field {
    case dns.QDCOUNT:
        header.qdcount = value
    case dns.ANCOUNT:
        header.ancount = value
    case dns.NSCOUNT:
        header.nscount = value
    case dns.ARCOUNT:
        header.arcount = value
    }
}

func (header *Header) SetField(field int, value uint16) {
    switch field {
    case dns.ID:
        header.id = value
    case dns.RCODE:
        header.status |= (value & 0x0f)
    case dns.RA:
        header.status |= (value & 0x01) << 7
    case dns.RD:
        header.status |= (value & 0x01) << 8
    case dns.TC:
        header.status |= (value & 0x01) << 9
    case dns.AA:
        header.status |= (value & 0x01) << 10
    case dns.OPCODE:
        header.status |= (value & 0x0f) << 14
    case dns.QR:
        header.status |= (value & 0x01) << 15
    }
}

func ReadFQName(s *bytes.Buffer, r []byte) []string {
    var oct_len uint8
    var data    []string
    for {
        binary.Read(io.Reader(s), binary.BigEndian, &oct_len)
        if oct_len == 0 {
            return data
        }
        if (oct_len & 0xc0) == 0xc0 {  // it's a pointer
            offset := int16(oct_len & 0x3f) * 256 
            binary.Read(io.Reader(s), binary.BigEndian, &oct_len)
            offset += int16(oct_len)

            return append(data, ReadFQName(bytes.NewBuffer(r[offset:]), r)...)
        }
        data = append(data, string(s.Next(int(oct_len))))
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

func (question *Question) Unpack(s *bytes.Buffer, r []byte) {
    question.qname = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &question.qtype)
    binary.Read(io.Reader(s), binary.BigEndian, &question.qclass)
}

func (question Question) Pack() []byte {
    buf := new(bytes.Buffer)

    buf.Write(WriteFQName(question.qname))
    binary.Write(buf, binary.BigEndian, question.qtype)
    binary.Write(buf, binary.BigEndian, question.qclass)
  
    return buf.Bytes()
}

func (answer *Answer) Unpack(s *bytes.Buffer, r []byte) {
    answer.name  = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.rr_type)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.class)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.rdlength)
    answer.rdata = s.Next(int(answer.rdlength))
}

func (authority *Authority) Unpack(s *bytes.Buffer, r []byte) {
    authority.name  = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.rr_type)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.class)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.rdlength)
    authority.rdata = s.Next(int(authority.rdlength))
}

func (additional *Additional) Unpack(s *bytes.Buffer, r []byte) {
    additional.name  = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &additional.rr_type)
    binary.Read(io.Reader(s), binary.BigEndian, &additional.class)
    binary.Read(io.Reader(s), binary.BigEndian, &additional.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &additional.rdlength)
    additional.rdata = s.Next(int(additional.rdlength))
}

func (message *Message) Unpack(s *bytes.Buffer) {
    r := s.Bytes()
    message.header.Unpack(s)

    message.question = make([]Question, message.header.qdcount)
    for i:=0; i<int(message.header.qdcount); i++ {
        message.question[i].Unpack(s, r)
    }

    message.answer = make([]Answer, message.header.ancount)
    for i:=0; i<int(message.header.ancount); i++ {
        message.answer[i].Unpack(s, r)
    }

    message.authority = make([]Authority, message.header.nscount)
    for i:=0; i<int(message.header.nscount); i++ {
        message.authority[i].Unpack(s, r)
    }

    message.additional = make([]Additional, message.header.arcount)
    for i:=0; i<int(message.header.arcount); i++ {
        message.additional[i].Unpack(s, r)
    }

    fmt.Println(message)
}

func (message Message) Pack() []byte {
    var data = message.header.Pack()

    for i:=0; i<int(message.header.qdcount); i++ {
        data = append(data, message.question[i].Pack()...)
    }

    return data
}



/*
type DNSAnswer struct {
    name     []string
    atype    uint16
    class    uint16 
    ttl      uint16
    rdlength uint16
    rdata    []byte
}

type DNSAuthority struct {
    name     []string
    atype    uint16
    class    uint16 
    ttl      uint16
    rdlength uint16
    rdata    []byte
}

type DNSServer struct {
    addr      *net.UDPAddr
    header    DNSHeader
    question  DNSQuestion 
    answer    DNSAnswer
    authority DNSAuthority
}
*/

/*
func (answer *DNSAnswer) Init(s *bytes.Buffer) {
    answer.name  = ReadFQName(s)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.atype)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.class)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.rdlength)
//  answer.rdata =
}

func (authority *DNSAuthority) Init(s *bytes.Buffer) {
    authority.name = ReadFQName(s)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.atype)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.class)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &authority.rdlength)
//  authority.rdata =
}
*/


func (message Message) HandleAIn() {
    fmt.Println("Handlke A In")
    fmt.Println(message.question[0].qname)

    connection, _ := net.Dial("udp", "192.168.1.3:53")
    fmt.Println(connection)

    var buf [dns.MAX_MESSAGE_LENGTH]byte

    connection.Write(message.Pack())
    len, _ := connection.Read(buf[:])

    fmt.Println("Answer")
    fmt.Println(buf[0:len])
    //r := bytes.NewBuffer(buf[:len])

    connection.Close()
}

func (message *Message) AddQuestion(qtype uint16, qclass uint16, name string) {
    message.question = append(message.question, Question{qtype: qtype, qclass: qclass, qname: strings.Split(name, ".")})
    message.header.qdcount++
}

func (message *Message) DNSQuery(name string, query uint16) {
    message.header.Init()

    message.header.SetField(dns.ID, uint16(rand.Int31n(0xffff)))

    message.header.SetField(dns.OPCODE, dns.QUERY)
    message.header.SetField(dns.RD, 1)

    message.AddQuestion(query, dns.IN, name)
}

func (message *Message) DNSQueryA(name string) {
    message.header.Init()

    message.header.SetField(dns.ID, uint16(rand.Int31n(0xffff)))

    message.header.SetField(dns.OPCODE, dns.QUERY)
    message.header.SetField(dns.RD, 1)

    message.AddQuestion(dns.A, dns.IN, name)
}

/*
func (server DNSServer) Run() {
    var new_message Message

    new_message.DNSQueryA("home.hanskramer.com")
    
    fmt.Println(new_message.Pack())

    server.addr, _  = net.ResolveUDPAddr("udp", ":53")
    sock, _         := net.ListenUDP("udp", server.addr)

    var buf [dns.MAX_MESSAGE_LENGTH]byte

    for {
        rlen, remote, err := sock.ReadFromUDP(buf[:])
        fmt.Println(err)
        fmt.Println(remote)

        fmt.Println(rlen)
        fmt.Println(buf[0:rlen])
        r := bytes.NewBuffer(buf[:rlen])
        fmt.Println(r)

        var message Message 
        message.Unpack(r)
        
        if message.header.qdcount>0 && message.question[0].qclass == dns.IN {
            switch message.question[0].qtype {
            case dns.A:
                message.HandleAIn()
            case dns.CNAME:
                message.HandleAIn()
            }
        }
        fmt.Println(message.Pack())
        fmt.Println(message.question[0].qtype)
        fmt.Println(message.question[0].qtype)
    }
}
*/

func uncompress(buf []byte) []byte {
/*
    r := bytes.NewBuffer(buf[:])

    var header Header
    header.Unpack(r)

    for (i:=0 ; i:=int(header.qdcount); i++) {
       hk  
    }
fmt.Println(header.qdcount)
fmt.Println(header.ancount)
fmt.Println(header.nscount)
fmt.Println(header.arcount)
*/

    new_buf := buf[:12]
    return new_buf
}

func (message Message) SendDNSQuery(server string) Message {
    connection, _ := net.Dial("udp", server)

    connection.Write(message.Pack())

    var buf [dns.MAX_MESSAGE_LENGTH]byte
    len, _ := connection.Read(buf[:])
    
    var answer Message
    r := bytes.NewBuffer(buf[:len])
    
    answer.Unpack(r)

    return answer
}


func main() {
    var query Message

    query.DNSQuery("www.xs4all.nl", dns.A)
    answer := query.SendDNSQuery("192.168.1.3:53")

    fmt.Println(answer.Pack())
/*
    connection.Write(message.Pack())
    len, _ := connection.Read(buf[:])

    fmt.Println("Answer")
    fmt.Println(buf[0:len])
    //r := bytes.NewBuffer(buf[:len])

    connection.Close()
    var server DNSServer

    server.Run()
*/
}

