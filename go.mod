module github.com/pigfall/tzzGoUtil

go 1.13

require (
	github.com/buger/goterm v0.0.0-20200322175922-2f3e71b85129
	github.com/go-kit/kit v0.11.0
	github.com/go-redis/redis/v8 v8.10.0
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/gopacket v1.1.17
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/lestrrat-go/jwx v1.1.7
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8
	github.com/streadway/amqp v1.0.0
	github.com/tealeg/xlsx/v3 v3.2.3
	github.com/vishvananda/netlink v1.1.0
	github.com/xuri/excelize/v2 v2.4.1-0.20210727163809-e9ae9b45b20a
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6
	golang.org/x/text v0.3.7 // indirect
	golang.zx2c4.com/wireguard v0.0.0-00010101000000-000000000000
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/faiface/beep => github.com/pigfall/beep v1.0.3-0.20210306054403-c0ca516fe23e

replace golang.zx2c4.com/wireguard => ../wintun-go
