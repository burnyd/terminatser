## TermiNATser

#### Summary

This is a small example of how one would take a onbox agent for Arista EOS and take collective metrics and stream telemetry data to a NATs pub/sub bus for telemetry data,

Currently this can take any openconfig origin data and push it out to a NATs bus.

This makes the assumption that there is a nats server currently running somewhere within the environment. 

This is how I ran it on my demo machine

```
docker run --name nats --rm -p 4222:4222 -p 8222:8222 nats --http_port 8222

```

##### Switch config.

```
management api gnmi
   transport grpc eos
   provider eos-native
```

Move the /bin/terminatster to the switch of /mnt/flash

This can either be ran from bash for testing purposes or as a daemon

##### From bash for testing
```
./terminatser -natsurl 172.20.20.1:4222 -natstopic terminatser -clientname eos  -gnmitarget 127.0.0.1:6030 -gnmiuser admin -gnmipassword admin -gnmipath /
```

```
daemon TerminAttr
 exec /mnt/flash/terminatser -natsurl 172.20.20.1:4222 -natstopic terminatser -clientname eos  -gnmitarget 127.0.0.1:6030 nmiuser admin -gnmipassword admin -gnmipath /
 no shutdown
```


#### Output data
cd controller/
go run main.go 
2023/01/31 11:49:21 Listening on [terminatser]

Taking a look on the nats bus.  This can be done with any nats client.

```
2023/01/31 11:50:25 [#1169] Received on [terminatser]: '{"path":"/netconf-state/schemas/schema[format=yang][identifier=openconfig-rib-bgp-attributes][version=2022-06-06]/location","target":"(172.20.20.2:6030) ","ts":"2023-01-30T16:31:26.266008565Z","value":"[NETCONF]"}'
```

This will also export metrics to prometheus from the controller or whomever runs the controller app.

```
gNMI_Metrics{Path="/interfaces/interface[name=Ethernet1]/config/mtu",Target="(172.20.20.2:6030) ",Timestamp="2023-01-30T16:31:27.127445838Z",Value="0"} 0
gNMI_Metrics{Path="/interfaces/interface[name=Ethernet1]/config/name",Target="(172.20.20.2:6030) ",Timestamp="2023-01-30T16:31:27.466111982Z",Value="Ethernet1"} 0
```

#### CLI Commands

Yes, we have the ability to send CLI commands to the devices.

in the cli-test directory there is a cli binary.

The way it works is nats will listen on the switch to a topic of its own hostname.  For example, if the hostname is ceos1 then ceos1 will listen on that nats topic.  Anything with data included with the commands key map it will respond with the CLI output.  This effectively will proxy eAPI over Terminattr.

Example.

```
cd cli-test
./cli -natsurl 172.20.20.1:4222 -devicename ceos1 -commands "show ip route"
```

Response 
```
2023/02/22 11:59:32 Reply: [map[output:
VRF: default
Codes: C - connected, S - static, K - kernel, 
       O - OSPF, IA - OSPF inter area, E1 - OSPF external type 1,
       E2 - OSPF external type 2, N1 - OSPF NSSA external type 1,
       N2 - OSPF NSSA external type2, B - Other BGP Routes,
       B I - iBGP, B E - eBGP, R - RIP, I L1 - IS-IS level 1,
       I L2 - IS-IS level 2, O3 - OSPFv3, A B - BGP Aggregate,
       A O - OSPF Summary, NG - Nexthop Group Static Route,
       V - VXLAN Control Service, M - Martian,
       DH - DHCP client installed default route,
       DP - Dynamic Policy Route, L - VRF Leaked,
       G  - gRIBI, RC - Route Cache Route

Gateway of last resort:
 S        0.0.0.0/0 [1/0] via 172.20.20.1, Management0

 C        1.1.1.1/32 is directly connected, Loopback0
 C        172.20.20.0/24 is directly connected, Management0

]]
```