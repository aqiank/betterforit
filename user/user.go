package user

import(
        "encoding/csv"
        "log"
        "os"
        "strings"

        "github.com/jackyb/betterforit/util"
)

type User struct {
        No string
        RegistrationID string
        Fullname string
        Email string
        TwitterHandle string
        FacebookHandle string
        InstagramHandle string
        CommitMessage string
}

func Load(filename string) []User {
        var us []User

        file, err := os.Open(filename)
        if err != nil {
                log.Println("Load:", err)
                return nil
        }

        r := csv.NewReader(file)
        records, err := r.ReadAll()
        if err != nil {
                log.Println("Load:", err)
                return nil
        }

        for _, v := range records[1:] {
                var u User

                u.No = strings.Trim(v[0], " \n\r")
                if u.No == "" {
                        // no record number, skip
                        continue
                }

                u.RegistrationID = strings.Trim(v[1], " \n\r")
                if u.RegistrationID == "" {
                        log.Println("Load:", err)
                        return nil
                }

                u.TwitterHandle = util.Strip(v[4], " \n\r")
                u.InstagramHandle = util.Strip(v[7], " \n\r")
                u.FacebookHandle = util.Strip(v[5], " ")
                u.Fullname = strings.Trim(v[2], " ")
                //u.Email = strings.Trim(v[3], " ")
                u.CommitMessage = strings.Trim(v[8], " ")

                us = append(us, u)
        }

        return us
}
