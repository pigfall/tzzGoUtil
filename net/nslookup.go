package net


import(
	"context"
	"net"
		"time"
)

func LookupHost(ctx context.Context,timeout time.Duration,host string)([]string,error){
	 r := &net.Resolver{
        PreferGo: true,
        Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
            d := net.Dialer{
                Timeout: timeout,
            }
            return d.DialContext(ctx, network, "8.8.8.8:53")
        },
    }
		return r.LookupHost(ctx,host)
}
