# boost

A [CoreDNS](https://github.com/coredns/coredns) plugin.

一个尝试返回访问最快的解析结果的 CoreDNS 插件。

## 警告

就是个随便练手写着玩的项目，不保证任何东西。

## 用法

最基本：

```
boost
```

所有配置：

```
boost [ZONE] {
    method ping
    ping_count 3
    ping_interval 0.025
    ping_timeout 0.5
}
```

## 实例

配合 [pforward](https://github.com/microdog/pforward) 使用。

`plugin.cfg`

```
# ...

cache:cache

# Must be placed after cache plugin
boost:github.com/microdog/boost

# ...

# Replace the built-in forward plugin
forward:github.com/coredns/coredns/plugin/pforward
#forward:forward

# ...
```

---

`Corefile`

```
.:55 {
    debug
    errors
    log
    cache
    boost
    forward . 119.29.29.29:53 223.5.5.5:53 1.0.0.1:53 8.8.4.4:53
}
```

---

```
$ ./coredns
.:55
CoreDNS-1.6.5
darwin/amd64, go1.13.4, 2503df90-dirty
[DEBUG] plugin/boost: collected answers:
	www.a.shifen.com.	230	IN	A	180.101.49.11
	www.a.shifen.com.	230	IN	A	180.101.49.12
	www.wshifen.com.	61	IN	A	104.193.88.123
	www.wshifen.com.	61	IN	A	104.193.88.77
[DEBUG] plugin/boost: ping stats: &{PacketsRecv:3 PacketsSent:3 PacketLoss:0 IPAddr:180.101.49.12 Addr:180.101.49.12 Rtts:[8.032ms 8.012ms 8.29ms] MinRtt:8.012ms MaxRtt:8.29ms AvgRtt:8.111333ms StdDevRtt:126.599µs}
[DEBUG] plugin/boost: ping stats: &{PacketsRecv:3 PacketsSent:3 PacketLoss:0 IPAddr:180.101.49.11 Addr:180.101.49.11 Rtts:[8.504ms 8.122ms 8.568ms] MinRtt:8.122ms MaxRtt:8.568ms AvgRtt:8.398ms StdDevRtt:196.902µs}
[DEBUG] plugin/boost: ping stats: &{PacketsRecv:3 PacketsSent:3 PacketLoss:0 IPAddr:104.193.88.123 Addr:104.193.88.123 Rtts:[152.56ms 152.564ms 152.079ms] MinRtt:152.079ms MaxRtt:152.564ms AvgRtt:152.401ms StdDevRtt:227.694µs}
[DEBUG] plugin/boost: ping stats: &{PacketsRecv:2 PacketsSent:3 PacketLoss:33.33333333333333 IPAddr:104.193.88.77 Addr:104.193.88.77 Rtts:[147.761ms 148.082ms] MinRtt:147.761ms MaxRtt:148.082ms AvgRtt:147.9215ms StdDevRtt:160.5µs}
[DEBUG] plugin/boost: sorted results
	pingStats(IP=180.101.49.12, Sent=3, Loss=0.000000, Avg=8.111333ms, StdDev=126.599µs)
	pingStats(IP=180.101.49.11, Sent=3, Loss=0.000000, Avg=8.398ms, StdDev=196.902µs)
	pingStats(IP=104.193.88.123, Sent=3, Loss=0.000000, Avg=152.401ms, StdDev=227.694µs)
	pingStats(IP=104.193.88.77, Sent=3, Loss=33.333333, Avg=147.9215ms, StdDev=160.5µs)
[DEBUG] plugin/boost: final response:
	;; opcode: QUERY, status: NOERROR, id: 49707
	;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

	;; QUESTION SECTION:
	;www.baidu.com.	IN	 A

	;; ANSWER SECTION:
	www.a.shifen.com.	230	IN	A	180.101.49.12
[INFO] 127.0.0.1:64337 - 49707 "A IN www.baidu.com. udp 42 false 4096" NOERROR qr,rd 63 0.684812174s
[INFO] 127.0.0.1:53541 - 55208 "A IN www.baidu.com. udp 42 false 4096" NOERROR qr,aa,rd 63 0.000073673s
[INFO] 127.0.0.1:57047 - 36400 "A IN www.baidu.com. udp 42 false 4096" NOERROR qr,aa,rd 63 0.000072611s
```
