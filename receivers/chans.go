package receivers

import (
	"github.com/kevinchapron/FSHK-final/messaging"
)

var UDPmessagesInternal = make(chan messaging.Message)

var TCPmessagesInternal = make(chan messaging.Message)
