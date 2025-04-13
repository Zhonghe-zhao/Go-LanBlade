在局域网下 发现局域网下的设备 并 实现共享 在局域网内的设备间快速、无需配置地传输文件

mDNS 使用组播避免了广播的泛滥问题！

使用了mDNS库 [mDNS](github.com/hashicorp/mdns)实现了对于局域网内的设备发现

利用了go的交叉编译，致使程序在跨平台使用非常方便！

### windows编译
`go build -o device_scanner_windows.exe main.go`

直接将打包好的可执行文件发送到另一台设备即可实现通信

### linux 编译
`go build -o devices_scanner_linux main.go`

使用 `scp` 指令 快速实现传送文件

`scp .\devices_scanner xxx@192.168.1.xxx:~/`

之后直接运行 ./devices_scanner 即可！



后续功能还在实现中......
