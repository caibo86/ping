// -------------------------------------------
// @file      : main.go
// @author    : bo cai
// @contact   : caibo923@gmail.com
// @time      : 2024/12/26 下午4:54
// -------------------------------------------

package main

import (
	"fmt"
	"github.com/caibo86/ping/misc"
	"os"
	"os/signal"
)

func main() {
	ping := misc.NewPing()
	signal.Notify(ping.Stop, os.Interrupt)
	ping.ParseArgs()
	if ping.Help {
		help()
		os.Exit(0)
	}
	ping.Run()
}

// 打印帮助信息
func help() {
	fmt.Println("Usage")
	fmt.Println("\t ping [options] <destination>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("\t-h\t\tprint help and exit")
	fmt.Println("\t-w <timeout>\ttime to wait for response")
	fmt.Println("\t-s <size>\tuse <size> as number of data bytes to be sent")
	fmt.Println("\t-c <count>\tstop after <count> packets sent")
}
