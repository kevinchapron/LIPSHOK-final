package receivers

import (
	"github.com/kevinchapron/LIPSHOK/messaging"
)

var UDPmessagesInternal = make(chan messaging.Message)

var TCPmessagesInternal = make(chan messaging.Message)

var BLEMessagesInternal = make(chan messaging.Message)

var ZWaveMessagesInternal = make(chan messaging.Message)
