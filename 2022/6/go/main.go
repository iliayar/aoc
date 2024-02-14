package main

import (
    "fmt"
    "os"
)

const W = 14

func main() {
    content, err := os.ReadFile("input.txt")
    if err != nil {
        panic(err)
    }
    
    curChars := make(map[byte]int)

    var l, r int
    
    for ; r < W && r < len(content); r += 1 {
        curChars[content[r]] += 1
    }

    for r < len(content) {
        chR := content[r]
        chL := content[l]

        if len(curChars) == W {
            fmt.Println(r)
            break
        }

        curChars[chL] -= 1
        if curChars[chL] == 0 {
            delete(curChars, chL)
        }

        curChars[chR] += 1

        r += 1
        l += 1
    }
}
