package main

import (
	"time"

	"github.com/zhaozhonghe/lanblade/discover"
)

func main() {
	go discover.AnnounceSelf()
	go discover.DiscoverDevices(10)

	// 等一会，打印发现的设备
	time.Sleep(12 * time.Second)
	discover.PrintDevices()
}
