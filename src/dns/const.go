/*

    Author : Hans Kramer

      Date : Jan 2015

      Code : Go implementation of (limited) DNS server
             Port from my Python Code

 */

package dns


const MAX_MESSAGE_LENGTH int = 512


// TYPE Values
const (
    A     = 1   // a host address
    NS    = 2   // an authoritative name server
    MD    = 3   // a mail destination (Obsolete - use MX)
    MF    = 4   // a mail forwarder (Obsolete - use MX)
    CNAME = 5   // the canonical name for an alias
    SOA   = 6   // marks the start of a zone of authority
    MB    = 7   // a mailbox domain name (EXPERIMENTAL)
    MG    = 8   // a mail group member (EXPERIMENTAL)
    MR    = 9   // a mail rename domain name (EXPERIMENTAL)
    NULL  = 10  // a null RR (EXPERIMENTAL)
    WKS   = 11  // a well known service description
    PTR   = 12  // a domain name pointer
    HINFO = 13  // host information
    MINFO = 14  // mailbox or mail list information
    MX    = 15  // mail exchange
    TXT   = 16  // text strings
    AXFR  = 252 // A request for a transfer of an entire zone
    MAILB = 253 // A request for mailbox-related records (MB, MG or MR)
    MAILA = 254 // A request for mail agent RRs (Obsolete - see MX)
    ALL   = 255 // A request for all records
)


var TypeMap = map[uint16]string {
    A     : "A",
    NS    : "NS",
    MX    : "MX",
    SOA   : "SOA",
    PTR   : "PTR",
    TXT   : "TXT",
    CNAME : "CNAME",
}


func Type2string(rr_type uint16) string {
    return TypeMap[rr_type]
}


func String2type(value string) uint16 {
    for k, v := range TypeMap {
	if v == value {
            return k
        }
    }

    return 0
}


// CLASS values
const (
    IN = 1 // the Internet
    CS = 2 // the CSNET class (Obsolete - used only for examples in some obsolete RFCs)
    CH = 3 // the CHAOS class
    HS = 4 // Hesiod [Dyer 87]
)

var ClassMap = map[uint16]string {
    IN : "Internet",
    CS : "CSNET",
    CH : "CHAOS",
    HS : "Hesiod",
}


const (
    OPCODE_QUERY  = 0 // a standard query (QUERY)
    OPCODE_IQUERY = 1 // an inverse query (IQUERY)
    OPCODE_STATUS = 2 // a server status request (STATUS)
)

const (
    ID     = iota
    STATUS = iota
    QR     = iota
    OPCODE = iota
    AA     = iota
    TC     = iota
    RD     = iota
    RA     = iota
    Z      = iota
    RCODE  = iota
)

/*
const (
    QDCOUNT = iota
    ANCOUNT = iota
    NSCOUNT = iota
    ARCOUNT = iota
)
*/
