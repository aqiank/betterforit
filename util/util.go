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

func IndexString(s, substr string) int {
        if substr == "" {
                return -1
        }

        var idx, lastidx int
        for {
                lastidx = idx + len(substr)
                if len(s) < lastidx {
                        break
                }

                if s[idx:lastidx] == substr {
                        return idx
                }
                idx++
        }


        return -1
}

func RemoveStringWithPrefix(s, prefix string) string {
        idx := IndexString(s, prefix)
        if idx == -1 {
                return s
        }

        spaceidx := strings.IndexByte(s[idx:], ' ')
        lastidx := idx + spaceidx
        if spaceidx == -1 {
                lastidx = len(s)
        }

        return strings.Replace(s, s[idx:lastidx], "", -1)
}

func RemoveStringsWithPrefix(s, prefix string) string {
        for idx := IndexString(s, prefix); idx != -1; idx = IndexString(s, prefix) {
                spaceidx := strings.IndexByte(s[idx:], ' ')
                lastidx := idx + spaceidx
                if spaceidx == -1 {
                        lastidx = len(s)
                }
                s = strings.Replace(s, s[idx:lastidx], "", -1)
        }
        return s
}
