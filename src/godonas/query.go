/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */

package main


import (
    "dns"
    "flag"
    "fmt"
)


func main() {
    server := flag.String("s", "8.8.8.8",  "DNS server to connect to")
    port   := flag.String("p", "53",       "DNS Server port")
    query  := flag.String("q", "A",        "DNS query type")
    class  := flag.String("c", "Internet", "DNS query class")

    flag.Parse()

    query_type := dns.String2type(*query)
    if query_type == 0 {
        fmt.Printf("Unknown query type: %s\n", *query)
        return
    }

    query_class := dns.Value2Key(dns.ClassMap, *class)
    if query_class == 0 {
        fmt.Printf("Unknown class type: %s\n", *class)
        return
    }

    var message dns.Message
    for _, question := range(flag.Args()) {
	fmt.Println(question)
        message.Query(question, query_type, query_class)
        answer := message.Send(*server + ":" + *port)
        fmt.Print(answer.String())
    }
}
