/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */

package dns

import (
    "bytes"
    "strings"
    "math/rand"
    "net"
//    "fmt"
)

// Message


type Message struct {
    header     Header
    question   []Question  // allow for multiple questions even though most DNS servers don't support this
    answer     []Answer
    authority  []Authority
    additional []Additional
}


func (message Message) Pack() []byte {
    var data = message.header.Pack()

    for i:=0; i<int(message.header.Qdcount); i++ {
        data = append(data, message.question[i].Pack()...)
    }

    return data
}


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

//    fmt.Println(message)
}


func (message *Message) AddQuestion(qtype uint16, qclass uint16, name string) {
    message.question = append(message.question, Question{qtype: qtype, qclass: qclass, qname: strings.Split(name, ".")})
    message.header.Qdcount++
}


func (message *Message) Query(name string, query uint16) {
    message.header.Init()

    message.header.SetField(ID, uint16(rand.Int31n(0xffff)))
    message.header.SetField(OPCODE, QUERY)
    message.header.SetField(RD, 1)

    message.AddQuestion(query, IN, name)
}


func (message Message) Send(server string) Message {
    connection, _ := net.Dial("udp", server)

    connection.Write(message.Pack())

    var buf [MAX_MESSAGE_LENGTH]byte
    len, _ := connection.Read(buf[:])

    var answer Message
    r := bytes.NewBuffer(buf[:len])

    answer.Unpack(r)

    return answer
}


func (message Message) Print() {
    for _, answer := range message.answer {
        answer.Print()
    }
}
