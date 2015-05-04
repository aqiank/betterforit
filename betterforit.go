package main

import (
        "log"
        "sync"
        "os"
        "os/signal"

        "github.com/jackyb/betterforit/cfg"
        "github.com/jackyb/betterforit/ig"
        "github.com/jackyb/betterforit/tw"
        "github.com/jackyb/betterforit/user"
)

func main() {
        cfg.Load()

        // load list of users from data.csv
        us := user.Load("data.csv")
        if us == nil {
                log.Println("main: failed to load users")
                return
        }
        log.Println("loaded", len(us), "users with social network handle(s)")

        // download user twitter and instagram #betterforit posts
        var wg sync.WaitGroup

        wg.Add(2)

        sigC := make(chan os.Signal, 1)
        signal.Notify(sigC, os.Interrupt, os.Kill)

        // twitter
        go func() {
                defer wg.Done()
                tw.Run(us, sigC)
        }()

        // instagram
        go func() {
                defer wg.Done()
                ig.Run(us, sigC)
        }()

        wg.Wait()

        cfg.Save()
}
