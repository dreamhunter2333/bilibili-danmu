# Bilibili Danmu

See your live broadcast room barrage at the terminal

在终端查看你的哔哩哔哩直播弹幕

## install

```bash
brew install jinmu333/dreamhunter/bilibili-danmu
```

```bash
NAME:
   Bilibili Danmu - bilibili-danmu -r 00000

USAGE:
   bilibili-danmu [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --room-id value, -r value  Bilibili Room ID
   --help, -h                 show help (default: false)
```

## build

```bash
CGO_ENABLED=0 GOOS=linux GOOS=windows go build -o out/bilibili-danmu.exe
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/bilibili-danmu-linux-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o out/bilibili-danmu-macos-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o out/bilibili-danmu-macos-arm64
```
