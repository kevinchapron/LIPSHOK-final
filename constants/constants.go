package constants

const (
	ROUTER_ADDRESS = "127.0.0.1"
	INTERFACE_PORT = 5002
)

const (
	SENSOR_WEBSOCKET_PATH = "/sensors"
	SENSOR_WEBSOCKET_PORT = 5001
	SENSOR_WEBSOCKET_NAME = "SENSORS_SETTINGS"
)

const (
	UDP_ADDR            = "127.0.0.1"
	UDP_PORT            = 5010
	MAX_UDP_PACKET_SIZE = 1000
)

const (
	TCP_ADDR            = "127.0.0.1"
	TCP_PORT            = 5020
	MAX_TCP_PACKET_SIZE = 1000
)

const (
	MESSAGING_DATATYPE_AUTH = 0x01
	MESSAGING_DATATYPE_DATA = 0x00
)

const PROTOCOLS_JSON_CONF = "protocols.json"
