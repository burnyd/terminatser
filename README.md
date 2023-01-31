docker run --name nats --rm -p 4222:4222 -p 8222:8222 nats --http_port 8222

go run main.go -natsurl 172.20.20.1:4222 -natstopic terminatser -clientname eos  \
-gnmitarget 172.20.20.2:6030 \
nmiuser admin -gnmipassword admin -gnmipath /

