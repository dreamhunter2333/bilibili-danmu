package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	Black   = Color("\033[1;30m%s\033[0m")
	Red     = Color("\033[1;31m%s\033[0m")
	Green   = Color("\033[1;32m%s\033[0m")
	Yellow  = Color("\033[1;33m%s\033[0m")
	Purple  = Color("\033[1;34m%s\033[0m")
	Magenta = Color("\033[1;35m%s\033[0m")
	Teal    = Color("\033[1;36m%s\033[0m")
	White   = Color("\033[1;37m%s\033[0m")
)

var colors = []func(...interface{}) string{Black, Red, Green, Yellow, Purple}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func PrintColor(args ...interface{}) string {
	return colors[rand.Intn(len(colors))](args...)
}

type DanmuRes struct {
	Data DanmuData `json:"data"`
}

type DanmuData struct {
	Room []DanmuItem `json:"room"`
}

type DanmuItem struct {
	Text      string        `json:"text"`
	Nickname  string        `json:"nickname"`
	Timeline  string        `json:"timeline"`
	CheckInfo CheckInfoItem `json:"check_info"`
}

type CheckInfoItem struct {
	Time int64 `json:"ts"`
}

func main() {
	var roomId string
	var refreshInterval int
	var lastRnd int64 = 0

	app := &cli.App{
		Name:  "Bilibili Danmu",
		Usage: "bilibili-danmu -r 00000",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "room-id",
				Aliases:     []string{"r"},
				Usage:       "Bilibili Room ID",
				Destination: &roomId,
			},
			&cli.IntFlag{
				Name:        "refresh",
				Usage:       "refresh interval",
				Value:       10,
				Destination: &refreshInterval,
			},
		},
		Action: func(c *cli.Context) error {
			client := &http.Client{}
			for {
				// fmt.Printf("\x1bc")
				req, err := http.NewRequest("GET", "https://api.live.bilibili.com/xlive/web-room/v1/dM/gethistory?roomid="+roomId, nil)
				if err != nil {
					return err
				}
				req.Header.Add("authority", "api.live.bilibili.com")
				res, err := client.Do(req)
				if err != nil {
					return err
				}
				defer func() { _ = res.Body.Close() }()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return err
				}
				var result DanmuRes
				err = json.Unmarshal(body, &result)
				if err != nil {
					return err
				}
				for _, danmuItem := range result.Data.Room {
					if danmuItem.CheckInfo.Time > lastRnd {
						fmt.Printf("%s %s: %s", danmuItem.Timeline, PrintColor(danmuItem.Nickname), PrintColor(danmuItem.Text))
						fmt.Println()
						lastRnd = danmuItem.CheckInfo.Time
					}
				}
				time.Sleep(time.Duration(refreshInterval) * time.Second)
			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
