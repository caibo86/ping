// -------------------------------------------
// @file      : ping.go
// @author    : bo cai
// @contact   : caibo923@gmail.com
// @time      : 2024/12/27 下午6:18
// -------------------------------------------

package misc

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// ICMP ICMP头
type ICMP struct {
	Type     uint8  // 类型
	Code     uint8  // 代码
	CheckSum uint16 // 校验和
	ID       uint16 // 标识符
	Seq      uint16 // 序号
}

// Ping Ping工具
type Ping struct {
	Destination string         // 目标地址
	Help        bool           // 是否显示帮助
	Timeout     int64          // 超时时间
	Size        int            // 数据包大小
	Count       int            // 发送次数
	SendCount   int            // 已发送次数
	RecvCount   int            // 已接收次数
	MaxTime     float64        // 最大耗时
	MinTime     float64        // 最小耗时
	TotalTime   float64        // 总耗时
	Stop        chan os.Signal // 停止信号
	ID          int            // 当前ID
}

// NewPing 创建Ping工具
func NewPing() *Ping {
	ping := &Ping{
		Stop: make(chan os.Signal),
	}
	return ping
}

// 打印统计信息
func (ping *Ping) printStatistics() {
	fmt.Printf("--- %s ping statistics ---\n", ping.Destination)
	fmt.Printf("%d packets transmitted, %d received, %d%% packet loss, time %dms\n",
		ping.SendCount, ping.RecvCount, (ping.SendCount-ping.RecvCount)*100/ping.SendCount, int64(ping.TotalTime))
	fmt.Printf("rtt min/avg/max = %.3f/%.3f/%.3f ms\n",
		ping.MinTime, ping.TotalTime/float64(ping.RecvCount), ping.MaxTime)
}

// ParseArgs 解析命令行参数
func (ping *Ping) ParseArgs() {
	if len(os.Args) < 2 {
		fmt.Println("ping: usage error: destination address required")
		os.Exit(1)
	}
	ping.Destination = os.Args[len(os.Args)-1]
	flag.Int64Var(&ping.Timeout, "w", 1000, "")
	flag.IntVar(&ping.Size, "s", 64, "")
	flag.IntVar(&ping.Count, "c", 4, "")
	flag.BoolVar(&ping.Help, "h", false, "")
	flag.Parse()
}

// Run 运行
func (ping *Ping) Run() {
	// 连接
	conn, err := net.DialTimeout("ip:icmp", ping.Destination, time.Duration(ping.Timeout)*time.Millisecond)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	// 远程地址
	rAddr := conn.RemoteAddr()
	fmt.Printf("PING %s (%s) %d(%d) bytes of data.\n", ping.Destination, rAddr, ping.Size, ping.Size+28)
	defer func() {
		ping.printStatistics()
	}()
	var buffer bytes.Buffer
	icmp := &ICMP{
		Type:     8,
		Code:     0,
		CheckSum: uint16(0),
	}
	data := make([]byte, ping.Size)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ping.Stop:
			return
		case <-ticker.C:
			ping.ID += 1
			icmp.ID = uint16(ping.ID)
			icmp.Seq = uint16(ping.ID)
			buffer.Reset()
			err = binary.Write(&buffer, binary.BigEndian, icmp)
			if err != nil {
				log.Println(err)
				break
			}
			buffer.Write(data)
			msg := buffer.Bytes()
			sum := checkSum(msg)
			msg[2] = byte(sum >> 8)
			msg[3] = byte(sum)
			// 设置超时时间
			err = conn.SetDeadline(time.Now().Add(time.Duration(ping.Timeout) * time.Millisecond))
			if err != nil {
				log.Println(err)
				break
			}
			// 发送数据
			var n int
			startTime := time.Now()
			n, err = conn.Write(msg)
			if err != nil {
				log.Println(err)
				break
			}
			ping.SendCount += 1
			buf := make([]byte, 128)
			n, err = conn.Read(buf)
			if err != nil {
				log.Println(err)
				if ping.SendCount >= ping.Count {
					return
				}
				break
			}
			ping.RecvCount += 1
			t := float64(time.Since(startTime).Microseconds()) / 1000
			fmt.Printf("%d bytes from %d.%d.%d.%d: icmp_seq=%d ttl=%d time=%.1f ms\n",
				n-28, buf[12], buf[13], buf[14], buf[15], icmp.Seq, buf[8], t)
			ping.MaxTime = max(ping.MaxTime, t)
			ping.MinTime = min(ping.MinTime, t)
			ping.TotalTime += t
			if ping.SendCount >= ping.Count {
				return
			}
		}
	}
}

// 计算校验和
func checkSum(data []byte) uint16 {
	length := len(data)
	index := 0
	var sum uint32
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	hi := sum >> 16
	for hi != 0 {
		sum = hi + (sum & 0xff)
		hi = sum >> 16
	}
	return uint16(^sum)
}
