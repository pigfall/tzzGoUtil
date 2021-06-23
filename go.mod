module github.com/Peanuttown/tzzGoUtil

go 1.13

require (
	github.com/buger/goterm v0.0.0-20200322175922-2f3e71b85129
	github.com/go-redis/redis/v8 v8.10.0
	github.com/google/gopacket v1.1.17
	github.com/lestrrat-go/jwx v1.1.7
	github.com/pkg/errors v0.9.1
	github.com/pkg/sftp v1.12.0
	github.com/streadway/amqp v1.0.0
	github.com/vishvananda/netlink v1.1.0
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/crypto v0.0.0-20201217014255-9d1352758620
	golang.org/x/sys v0.0.0-20210112080510-489259a85091
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/faiface/beep => github.com/Peanuttown/beep v1.0.3-0.20210306054403-c0ca516fe23e
