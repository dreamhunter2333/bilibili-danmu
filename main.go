package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	live "github.com/iyear/biligo-live"
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

func main() {
	var roomId int64
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	app := &cli.App{
		Name:  "Bilibili Danmu",
		Usage: "bilibili-danmu -r 00000",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:        "room-id",
				Aliases:     []string{"r"},
				Usage:       "Bilibili Room ID",
				Destination: &roomId,
			},
		},
		Action: func(c *cli.Context) error {
			// 获取一个Live实例
			// debug: debug模式，输出一些额外的信息
			// heartbeat: 心跳包发送间隔。不发送心跳包，70 秒之后会断开连接，通常每 30 秒发送 1 次
			// cache: Rev channel 的缓存
			// recover: panic recover后的操作函数
			l := live.NewLive(true, 30*time.Second, 0, func(err error) {
				log.Fatal(err)
			})

			// 连接ws服务器
			// dialer: ws dialer
			// host: bilibili live ws host
			if err := l.Conn(websocket.DefaultDialer, live.WsDefaultHost); err != nil {
				log.Fatal(err)
				return err
			}

			ctx, stop := context.WithCancel(context.Background())

			go func() {
				if err := l.Enter(ctx, roomId, "", 0); err != nil {
					log.Fatal(err)
					return
				}
			}()

			go rev(ctx, l)

			<-interrupt
			fmt.Println("stoping")
			// 关闭ws连接与相关协程
			stop()
			fmt.Println("stoped")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func rev(ctx context.Context, l *live.Live) {
	for {
		select {
		case tp := <-l.Rev:
			if tp.Error != nil {
				log.Println(tp.Error)
				continue
			}
			handle(tp.Msg)
		case <-ctx.Done():
			log.Println("rev func stopped")
			return
		}
	}
}

func handle(msg live.Msg) {
	// 使用 msg.(type) 进行事件跳转和处理，常见事件基本都完成了解析(Parse)功能，不常见的功能有一些实在太难抓取
	// 更多注释和说明等待添加
	switch msg.(type) {
	// 心跳回应直播间人气值
	case *live.MsgHeartbeatReply:
		log.Printf("直播间人气值: %d\n", msg.(*live.MsgHeartbeatReply).GetHot())
	// 弹幕消息
	case *live.MsgDanmaku:
		dm, err := msg.(*live.MsgDanmaku).Parse()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%s %s: %s\n", time.Unix(dm.Time/1000, 0).Format("2006-01-02 15:04:05"), PrintColor(dm.Uname), PrintColor(dm.Content))
	// 礼物消息
	case *live.MsgSendGift:
		g, err := msg.(*live.MsgSendGift).Parse()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("%s: %s %d个%s\n", g.Action, g.Uname, g.Num, g.GiftName)
	// 直播间粉丝数变化消息
	case *live.MsgFansUpdate:
		f, err := msg.(*live.MsgFansUpdate).Parse()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("room: %d, fans: %d, fansClub: %d\n", f.RoomID, f.Fans, f.FansClub)
	}
}
