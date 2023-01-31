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
./terminatser -natsurl 172.20.20.1:4222 -natstopic terminatser -clientname eos  -gnmitarget 127.0.0.1:6030 nmiuser admin -gnmipassword admin -gnmipath /
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



