package discover

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/zhaozhonghe/lanblade/util"
)

type Device struct {
	Name string
	IP   string
	Port int
	Host string
}

var Devices = make([]Device, 0)

func PrintDevices() {
	if len(Devices) == 0 {
		println("No devices found")
		return
	}
	for _, device := range Devices {
		log.Printf("Name: %s, IP: %s, Port: %d, Host: %s", device.Name, device.IP, device.Port, device.Host)
	}

}

func DiscoverDevices(timeoutSec int) {
	customLogger := log.New(os.Stdout, "[mDNS] ", log.LstdFlags|log.Lmicroseconds)
	entriesCh := make(chan *mdns.ServiceEntry, 32)
	devicesMap := make(map[string]*mdns.ServiceEntry)

	go func() {
		timeout := time.After(time.Duration(timeoutSec) * time.Second)
		for {
			select {
			case entry := <-entriesCh:
				if entry == nil || entry.AddrV4 == nil {
					continue
				}
				key := entry.AddrV4.String() + ":" + strconv.Itoa(entry.Port) + "|" + entry.Name
				if _, exists := devicesMap[key]; !exists {
					devicesMap[key] = entry
					Devices = append(Devices, Device{
						Name: entry.Name,
						IP:   entry.AddrV4.String(),
						Port: entry.Port,
						Host: entry.Host,
					})
					customLogger.Printf("发现 %-40s IP: %-15s 端口: %-5d 主机: %s",
						util.TruncateString(entry.Name, 35),
						entry.AddrV4,
						entry.Port,
						util.TruncateString(entry.Host, 25),
					)
				}
			case <-timeout:
				customLogger.Printf("\n✅ 扫描完成! 共发现 %d 个设备", len(Devices))
				close(entriesCh)
				return
			}
		}
	}()

	services := []string{"_workstation._tcp", "_smb._tcp", "_ssh._tcp"}
	interfaces, err := net.Interfaces()
	if err != nil {
		customLogger.Fatalf("获取网络接口失败: %v", err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagMulticast == 0 {
			continue
		}
		for _, service := range services {
			go func(svc string, iface net.Interface) {
				param := &mdns.QueryParam{
					Service:             svc,
					Domain:              "local",
					Timeout:             2 * time.Second,
					Entries:             entriesCh,
					Interface:           &iface,
					WantUnicastResponse: true,
					DisableIPv4:         false,
					DisableIPv6:         true,
					Logger:              customLogger,
				}
				if err := mdns.Query(param); err != nil {
					customLogger.Printf("接口 %s 查询 %s 失败: %v", iface.Name, svc, err)
				}
			}(service, iface)
		}
	}
}

func AnnounceSelf() {
	host, _ := os.Hostname()
	info := []string{"MyLANBox"}

	service, err := mdns.NewMDNSService(
		host, "_workstation._tcp", "local.", "",
		8000, nil, info,
	)
	if err != nil {
		log.Fatalf("注册服务失败: %v", err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		log.Fatalf("启动 mDNS Server 失败: %v", err)
	}

	log.Printf("已广播服务: %s._workstation._tcp.local", host)

	// 阻止退出
	go func() {
		select {}
	}()
	_ = server
}
