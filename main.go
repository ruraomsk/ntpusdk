package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ntpusdk/setup"
	"github.com/ruraomsk/ntpusdk/tester"
	"github.com/ruraomsk/ntpusdk/transport"
)

var (
	//go:embed config
	config embed.FS
)

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(config, "config/config.toml", &setup.Set); err != nil {
		fmt.Println("Dissmis config.toml")
		os.Exit(-1)
		return
	}
	os.MkdirAll(setup.Set.LogPath, 0777)
}
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
	logger.Info.Println("Start ntp server for asdu ...")
	fmt.Println("\nStart ntp server for asdu ...")
	go transport.ListenExternalDevices()
	go tester.RunTester()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Wait make abort...")
	time.Sleep(3 * time.Second)
	logger.Info.Println("Exit ntp server for asdu...")
	fmt.Println("\nExit ntp server for asdu...")
}
