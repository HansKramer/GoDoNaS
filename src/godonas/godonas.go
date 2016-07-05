/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */
package main

import "dns"

import (
    "fmt"
)

/*
    "bytes"
    "encoding/binary"
    "io"
*/
    // "net"
    // "strings"
    // "math/rand"



/*
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
*/

/*
func (answer *Answer) Unpack(s *bytes.Buffer, r []byte) {
    answer.name  = ReadFQName(s, r)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.rr_type)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.class)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.ttl)
    binary.Read(io.Reader(s), binary.BigEndian, &answer.rdlength)
    answer.rdata = s.Next(int(answer.rdlength))
}
*/

/*
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
*/

/*
func (message *Message) Unpack(s *bytes.Buffer) {
    r := s.Bytes()
    message.header.Unpack(s)

    message.question = make([]Question, message.header.Qdcount)
    for i:=0; i<int(message.header.Qdcount); i++ {
        message.question[i].Unpack(s, r)
    }

    message.answer = make([]Answer, message.header.Ancount)
    for i:=0; i<int(message.header.Ancount); i++ {
        message.answer[i].Unpack(s, r)
    }

    message.authority = make([]Authority, message.header.Nscount)
    for i:=0; i<int(message.header.Nscount); i++ {
        message.authority[i].Unpack(s, r)
    }

    message.additional = make([]Additional, message.header.Arcount)
    for i:=0; i<int(message.header.Arcount); i++ {
        message.additional[i].Unpack(s, r)
    }

    fmt.Println(message)
}
*/



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


/*
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

func (message *Message) DNSQueryA(name string) {
    message.header.Init()

    message.header.SetField(dns.ID, uint16(rand.Int31n(0xffff)))

    message.header.SetField(dns.OPCODE, dns.QUERY)
    message.header.SetField(dns.RD, 1)

    message.AddQuestion(dns.A, dns.IN, name)
}
*/

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


func main() {
    var query dns.Message

    query.DNSQuery("www.xs4all.nl", dns.A)
    answer := query.SendDNSQuery("75.75.75.75:53")

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

