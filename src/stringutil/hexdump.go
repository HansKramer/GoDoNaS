package stringutil

import (
   "fmt"
//   "strconv"
   "unicode"
)

func reset(a []byte) {
}


func isPrintable(x rune) bool {
    return x<127 && unicode.IsPrint(x)
}


func Hexdump(buffer []byte) (result string) {
    var ascii [16]byte
    for i := range ascii { 
        ascii[i] = '.' 
    }
    for i, b := range buffer[:] {
        result += fmt.Sprintf("%02X ", b)
        if isPrintable(rune(b)) {
            ascii[i%16] = b
        }
        if i % 16 == 15 {
            result += fmt.Sprintf(" %s\n", ascii)
            for i := range ascii {
                ascii[i] = '.' 
            }
        }
    }
    if len(buffer) % 16 != 0 {
        for i := 0; i < 16 - (len(buffer) % 16); i++ {
            result += fmt.Sprintf("   ")
        } 
        result += fmt.Sprintf(" %s\n", ascii[:len(buffer)%16])
    }

    return
}

