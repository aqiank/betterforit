package ig

import (
        "log"
        "os"
        "strings"
        "time"

        "github.com/jackyb/betterforit/cfg"
        "github.com/jackyb/betterforit/net"
        "github.com/jackyb/betterforit/user"
        "github.com/jackyb/betterforit/util"
        "github.com/carbocation/go-instagram/instagram"
)

var ig = instagram.NewClient(nil)

func init() {
	ig.ClientID = "348f92826e204525a149ab4c3ff3d180"
}

func Run(us []user.User, sigC chan os.Signal) {
        log.Println("instagram: starting..")

        go func() {
                last := 0
                for i := 0; i < len(us); i++ {
                        last = i
                        if us[i].No == cfg.LastInstagramUser() {
                               break
                        }
                }
                for i := last; i < len(us); i++ {
                        if len(us[i].InstagramHandle) < 6 {
                                continue
                        }

                        cfg.SetLastInstagramUser(us[i].No)
                        downloadUserMedia(us[i], sigC)
                }
        }()

        sigC <- <-sigC

        log.Println("instagram: stopping..")
}

func downloadUserMedia(u user.User, sigC chan os.Signal) {
        username := strings.Trim(u.InstagramHandle, "@")

        igUsers, _, err := ig.Users.Search(username, &instagram.Parameters{Count: uint64(cfg.InstagramRate())})
        if err != nil {
               log.Printf("instagram: downloadUserMedia: %s: %v\n", username, err)
               return
        }

        if len(igUsers) == 0 {
                log.Printf("instagram: couldn't find user called %s\n", username)
                return
        }

        log.Printf("instagram: searching [%s %s %s]\n", u.No, u.Fullname, username)

        prevMaxID := ""
        minTimestamp, _ := time.Parse("2006-01-02", "2015-02-01")
        params := &instagram.Parameters{Count: 1, MinTimestamp: minTimestamp.Unix()}
        for {
                medias, next, err := ig.Users.RecentMedia(igUsers[0].ID, params)
                if err != nil {
                       log.Printf("instagram: downloadUserMedia: %s: %v\n", username, err)
                       return
                }

                for _, m := range medias {
        		if m.Type != "image" {
        			continue
        		}
                        for _, tag := range m.Tags {
                                if tag == "betterforit" {
                                        dst := "downloads/instagram/" + u.No + " " + u.Fullname + " " + username
                                        os.MkdirAll(dst, os.ModeDir | 0700)
                        		if err = net.Download(dst, m.Images.StandardResolution.URL); err == nil {
                                                downloadProfilePic(igUsers[0], dst)
                                        }
                                        break
                                }
                        }
        	}

		if next.NextMaxID == "" {
			break
		}

                prevMaxID = params.MaxID
		params.MaxID = next.NextMaxID
                if params.MaxID == prevMaxID {
                        break
                }

                time.Sleep(1 * time.Second)
        }
}


func downloadProfilePic(u instagram.User, dst string) {
        if util.Exists(dst + "/profile.jpg") {
                return
        }
        net.DownloadAs(dst, u.ProfilePicture, "profile.jpg")
}
