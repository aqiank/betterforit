package cfg

import (
        "encoding/json"
        "log"
        "os"
)

var config = struct {
        LastTwitterUser string `json:last_twitter_user`
        LastInstagramUser string `json:last_instagram_user`
        TwitterRate string `json:"twitter_rate"`
        TwitterCount int `json:"twitter_count"`
        InstagramRate int `json:"instagram_rate"`
}{
        "1",
        "1",
        "2",
        33,
        1,
}

func Load() {
        file, err := os.Open("config.json")
        if err != nil {
                if !os.IsNotExist(err) {
                        log.Println("cfg: load:", err)
                        return
                }
        }
        defer file.Close()

        if err = json.NewDecoder(file).Decode(&config); err != nil {
                log.Println("cfg: load:", err)
        }
}

func Save() {
        file, err := os.Create("config.json")
        if err != nil {
                log.Println("cfg: save:", err)
                return
        }
        defer file.Close()

        if err = json.NewEncoder(file).Encode(&config); err != nil {
                log.Println("cfg: save:", err)
        }
}

func LastTwitterUser() string {
        return config.LastTwitterUser
}

func LastInstagramUser() string {
        return config.LastInstagramUser
}

func SetLastTwitterUser(n string) {
        config.LastTwitterUser = n
}

func SetLastInstagramUser(n string) {
        config.LastInstagramUser = n
}

func TwitterRate() string {
        return config.TwitterRate
}

func TwitterCount() int {
        return config.TwitterCount
}

func InstagramRate() int {
        return config.InstagramRate
}
