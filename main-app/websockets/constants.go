package websockets

import "time"

const (
	WEBSOCKET_MAX_WRITE_TIME  = time.Second * 10
	WEBSOCKET_MAX_PACKET_SIZE = 65536
)
