package video

import "net"

type Video struct {
	Name      string
	SignedURL net.Addr
}
