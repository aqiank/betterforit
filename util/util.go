package util

import (
        "strings"
        "os"
)

func Strip(s, chars string) string {
        return strings.Map(func(r rune) rune {
               if strings.IndexRune(chars, r) == -1 {
                       return r
               }
               return -1
        }, s)
}

func Exists(path string) bool {
        _, err := os.Stat(path)
        return err == nil
}
