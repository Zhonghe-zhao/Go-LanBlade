package main

import (
	"time"

	"github.com/zhaozhonghe/lanblade/discover"
	"github.com/zhaozhonghe/lanblade/lanmsg"
)

func main() {
	go discover.AnnounceSelf()
	go discover.DiscoverDevices(10)

	time.Sleep(12 * time.Second)
	discover.PrintDevices()
	lanmsg.StartMessaging("ZhaoZhongHe") // 替换为你的设备名
}
