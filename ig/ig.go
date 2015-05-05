package ig

import (
        "log"
        "os"
        "path"
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
        log.Println("instagram: rate is set to", cfg.InstagramRate())
        log.Println("instagram: minDate is set to", cfg.InstagramMinDate())

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
        dst := "downloads/instagram/" + u.No + " " + u.Fullname + " " + username

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

        // download profile pic
        downloadProfilePic(igUsers[0], dst)

        prevMaxID := ""
        minTimestamp, _ := time.Parse("2006-01-02", cfg.InstagramMinDate())
        params := &instagram.Parameters{Count: 1, MinTimestamp: minTimestamp.Unix()}
        for {
                medias, next, err := ig.Users.RecentMedia(igUsers[0].ID, params)
                if err != nil {
                       log.Printf("instagram: downloadUserMedia: %s: %v\n", username, err)
                       return
                }

                // look for photos that have #betterforit hashtag
                for _, m := range medias {
        		if m.Type != "image" {
        			continue
        		}
                        for _, tag := range m.Tags {
                                if tag == "betterforit" {
                                        // save caption
                                        saveCaption(m, dst, m.Images.StandardResolution.URL)

                                        // download image
                        		if err := net.Download(dst, m.Images.StandardResolution.URL, ""); err != nil {
                                                log.Println("instagram: downloadUserMedia:", err)
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

        if err := net.Download(dst, u.ProfilePicture, "profile.jpg"); err != nil {
                log.Println("instagram: downloadProfilePic:", err)
        }
}

func saveCaption(media instagram.Media, dst, mediaUrl string) {
        name := path.Base(mediaUrl) + ".txt"
        net.SaveText(dst, media.Caption.Text, name)
}
