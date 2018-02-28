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
    "time"
    "log"
    "fmt"
    "stringutil"
)


// Message

type Message struct {
    header     Header
    question   []Question  // allow for multiple questions even though most, if not all, DNS servers don't support this
    answer     []Answer
    authority  []Authority
    additional []Additional
}


func (message *Message) Recv(sock *net.UDPConn) {
    var buffer [MAX_MESSAGE_LENGTH]byte
    if rlen, remote, err := sock.ReadFromUDP(buffer[:]); err == nil {
        log.Printf("%s %d", remote, rlen)
        fmt.Println(stringutil.Hexdump(buffer[0:rlen]))

	//message.Unpack(bytes.NewBuffer(buffer[:rlen]))
        fmt.Println(message.String())
    } else {
       log.Fatal(err)
    }
}


func (message Message) Pack() []byte {
    var data = message.header.Pack()

    for i:=0; i<int(message.header.Qdcount); i++ {
        data = append(data, message.question[i].Pack()...)
    }

    return data
}


func (message *Message) Unpack(m MessageStream) {
    message.header.Unpack(m.s)

    message.question = make([]Question, message.header.Qdcount)
    for i:=0; i<int(message.header.Qdcount); i++ {
        message.question[i].Unpack(m)
    }

    message.answer = make([]Answer, message.header.Ancount)
    for i:=0; i<int(message.header.Ancount); i++ {
        message.answer[i].Unpack(m)
    }

    message.authority = make([]Authority, message.header.Nscount)
    for i:=0; i<int(message.header.Nscount); i++ {
        message.authority[i].Unpack(m)
    }

    message.additional = make([]Additional, message.header.Arcount)
    for i:=0; i<int(message.header.Arcount); i++ {
        message.additional[i].Unpack(m)
    }
}


func (message *Message) AddQuestion(qtype uint16, qclass uint16, name string) {
    message.question = append(message.question, Question{qtype: qtype, qclass: qclass, qname: strings.Split(name, ".")})
    message.header.Qdcount++
}


func (message *Message) Query(name string, query uint16, class uint16) {
    message.header.Init()

    message.header.SetField(ID, uint16(rand.Int31n(0xffff)))
    message.header.SetField(OPCODE, OPCODE_QUERY)
    message.header.SetField(RD, 1)

    message.question = message.question[:0]
    message.AddQuestion(query, class, name)
}


func (message Message) SendRaw(server string) (stream MessageStream) {
    connection, err := net.Dial("udp", server)
    if err != nil { return }
    defer connection.Close()

    connection.Write(message.Pack())   // error handling ... nah

    var buf [MAX_MESSAGE_LENGTH]byte
    connection.SetReadDeadline(time.Now().Add(2*time.Second))
    len, err:= connection.Read(buf[:])
    if err != nil { return }

    return MessageStream{bytes.NewBuffer(buf[:len]), buf[:len]}
}


func (message Message) Send(server string) (answer Message) {
    if connection, err := net.Dial("udp", server); err == nil {
	defer connection.Close()
        connection.Write(message.Pack())   // error handling ... nah

        var buf [MAX_MESSAGE_LENGTH]byte
	connection.SetReadDeadline(time.Now().Add(2*time.Second))
	len, err:= connection.Read(buf[:])
	if err != nil { return }

        answer.Unpack(MessageStream{bytes.NewBuffer(buf[:len]), buf[:len]})
    }

    return
}


func (message Message) Question() ([]string, bool) {
    if message.header.Qdcount >= 1 {
        return message.question[0].Get(), true
    }
    return nil, false
}


func (message Message) Answer() (net.IP, bool) {
   fmt.Println(message.header.Ancount)
   if message.header.Ancount > 0 {
       return message.answer[0].Get(), true
   } else {
       return nil, false
   }
}


func (message Message) String() string {
    result := "HEADER\n"
    result += message.header.String()

    result += "QUESTION\n"
    for _, question := range message.question {
        result += question.String() + "\n"
    }

    result += "ANSWER\n"
    for _, answer := range message.answer {
        fmt.Println("AAA" + answer.String())
    }

    result += "AUTHORITY\n"
    for _, authority := range message.authority {
        result += authority.String()
    }

    result += "ADDITIONAL\n"
    for _, additional := range message.additional {
        result += additional.String()
    }

    return result
}
