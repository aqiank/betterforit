package tw

import (
        "log"
        "net/url"
        "os"
        "strings"
        // "time"

        "github.com/jackyb/betterforit/cfg"
        "github.com/jackyb/betterforit/net"
        "github.com/jackyb/betterforit/user"
        "github.com/jackyb/betterforit/util"
        "github.com/ChimeraCoder/anaconda"
)

const (
        MaxUint = ^uint(0)
        MinUint = 0
        MaxInt = int(MaxUint >> 1)
        MinInt = -MaxInt - 1
)

var tw *anaconda.TwitterApi

func init() {
        anaconda.SetConsumerKey("A5WIV8Jto8EbpRP8xM0uMKj3W")
        anaconda.SetConsumerSecret("QODTrmmsYOmYEVaYvcIzk8FB8p3VBJc1jpx8N2EUNKbVLmEBEI")
        tw = anaconda.NewTwitterApi("7483012-1N7FQlwdAmGvAekFYDHYi3gFx8xwR41g7uun0gW5Rh", "ftvcZ1CzjeXYFiGa9gETm2aLrepv6Rz1nKTmNXOvuBCg4")
}

func Run(us []user.User, sigC chan os.Signal) {
        log.Println("twitter: starting..")
        log.Println("twitter: rate is set to", cfg.TwitterRate())
        log.Println("twitter: count is set to", cfg.TwitterCount())

        go func() {
                last := 0
                for i := 0; i < len(us); i++ {
                        last = i
                        if us[i].No == cfg.LastTwitterUser() {
                               break
                        }
                }
                for i := last; i < len(us); i++ {
                        if len(us[i].TwitterHandle) < 6 {
                                continue
                        }

                        cfg.SetLastTwitterUser(us[i].No)
                        downloadUserMedia(us[i], sigC)
                }
        }()

        sigC <- <-sigC

        log.Println("twitter: stopping..")
}

func downloadUserMedia(u user.User, sigC chan os.Signal) {
        username := strings.Trim(u.TwitterHandle, "@")
        dst := "downloads/twitter/" + u.No + " " + u.Fullname + " " + username

        params := url.Values{}
        params.Set("screen_name", username)
        params.Set("count", cfg.TwitterRate())
        params.Set("include_rts", "false")

        tu, err := tw.GetUsersShow(username, nil)
        if err != nil {
                log.Println("twitter: downloadProfilePic:", err)
                return
        }

        log.Printf("twitter: searching [%s %s %s]\n", u.No, u.Fullname, username)

        // download profile pic if available
        if !tu.DefaultProfileImage {
                downloadProfilePic(tu, dst)
        }

        maxID := ""
        cnt := cfg.TwitterCount()
        for cnt > 0 {
                tweets, err := tw.GetUserTimeline(params)
                if err != nil {
                        log.Printf("twitter: DownloadUserMedia: %v\n", err)
                        break
                }

                if len(tweets) == 0 {
                        log.Println("twitter: no tweets found")
                        break
                }

                cnt -= len(tweets)
                //log.Println("twitter: found", len(tweets), "tweets")

                for _, tweet := range tweets {
                        var hasTag bool

                        maxID = tweet.IdStr

                        for _, tag := range tweet.Entities.Hashtags {
                                if strings.EqualFold(tag.Text, "betterforit") {
                                        hasTag = true
                                }
                        }

                        if !hasTag {
                                continue
                        }

                        log.Println("twitter: found tweet with #betterforit hashtag")

                        // download tweet
                        tweetDst := dst + "/" + tweet.IdStr
                        net.SaveText(tweetDst, tweet.Text, "tweet.txt")

                        for _, media := range tweet.Entities.Media {
                                if media.Type != "photo" {
                                        continue
                                }

                                // download images
                                if err := net.Download(tweetDst, media.Media_url, ""); err != nil {
                                        log.Println("twitter: downloadUserMedia:", err)
                                }
                        }
                }

                params.Set("max_id", maxID)
                //time.Sleep(1 * time.Second)
        }
}

func downloadProfilePic(tu anaconda.User, dst string) {
        if util.Exists(dst + "/profile.jpg") {
                return
        }
        url := strings.Replace(tu.ProfileImageURL, "_normal", "", -1)
        net.Download(dst, url, "profile.jpg")
}
