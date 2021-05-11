package main

import (
	"encoding/binary"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/LIPSHOK/constants"
	"github.com/kevinchapron/LIPSHOK/messaging"
	"github.com/kevinchapron/LIPSHOK/security"
	"gobot.io/x/gobot/platforms/ble"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const DEVICE_WRISTBAND_BATTERY_SERVICE = "c2211ac0-41e2-11e9-b210-d663bd873d93"
const DEVICE_WRISTBAND_BATTERY_CHARACTERISTIC_VALUE = "2a19"
const DEVICE_WRISTBAND_CLASSIFIER_SERVICE = "fd3744da-49d1-11e9-8646-d663bd873d93"
const DEVICE_WRISTBAND_CLASSIFIER_CHARACTERISTIC = "fd374958-49d1-11e9-8646-d663bd873d93"
const DEVICE_WRISTBAND_IMU_RAW_DATA_CHARACTERISTIC = "fd374972-49d1-11e9-8646-d663bd873d93"
const DEVICE_WRISTBAND_DEBUGGING_CHARACTERISTIC = "e1ab4d54-b549-11e9-a2a3-2a2ae2dbcce4"
const DEVICE_WRISTBAND_CLASSIFICATIONS_OF_THE_DAY_CHARACTERISTIC = "261f6059-ace7-497c-b7b0-6b3d80efe6bf"
const DEVICE_WRISTBAND_ENCRYPTION_KEY_SETUP = "244cc952-c29a-4bb6-84cc-1d35d05fcd89"

const TIMESTAMP_STOPPER = time.Second * 30

func ExtractFloatValuesFromBLE(b []byte) (float64, float64, float64, float64, float64, float64) {
	ax := (float64(binary.LittleEndian.Uint16(b[0:2])) / 1000) - 4
	ay := (float64(binary.LittleEndian.Uint16(b[2:4])) / 1000) - 4
	az := (float64(binary.LittleEndian.Uint16(b[4:6])) / 1000) - 4

	gx := (float64(binary.BigEndian.Uint32(append([]byte{0x00}, b[6:9]...))) / 1000) - 2000
	gy := (float64(binary.BigEndian.Uint32(append([]byte{0x00}, b[9:12]...))) / 1000) - 2000
	gz := (float64(binary.BigEndian.Uint32(append([]byte{0x00}, b[12:15]...))) / 1000) - 2000
	return ax, ay, az, gx, gy, gz
}

const _defaultBLEAddr = "EF:3A:4E:89:07:0D"
const _defaultIPAddr = "192.168.0.1"
const _defaultIPPort = constants.WEBSOCKET_INNER_BLE_PORT
const _defaultPath = constants.WEBSOCKET_INNER_BLE_PATH

func main() {
	Logging.SetLevel(Logging.DEBUG)
	lastTimestamp := time.Now()
	addr := flag.String("addr", _defaultBLEAddr, "Addr of wristband")
	ipAddr := flag.String("ip", _defaultIPAddr, "IP of the websocket")
	wsPort := flag.String("port", strconv.Itoa(_defaultIPPort), "Port of the websocket")
	wsPath := flag.String("path", _defaultPath, "Path in the websocket")

	flag.Parse()
	Logging.Info("Launching wristbandBLE with following parameters : ")
	_m := map[bool]string{true: "(default)", false: ""}
	Logging.Info(" -addr=\"" + (*addr) + "\" " + _m[*addr == _defaultBLEAddr])
	Logging.Info(" -ip=\"" + (*ipAddr) + "\" " + _m[*ipAddr == _defaultIPAddr])
	Logging.Info(" -port=\"" + (*wsPort) + "\" " + _m[*wsPort == strconv.Itoa(_defaultIPPort)])
	Logging.Info(" -path=\"" + (*wsPath) + "\" " + _m[*wsPath == _defaultPath])
	Logging.Info()

	// Connecting to WS
	u := url.URL{Scheme: "ws", Host: (*ipAddr) + ":" + *wsPort, Path: *wsPath, RawQuery: "Name=Wristband&Protocol=BLE"}
	Logging.Debug("Trying to connect to inner WS : " + u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		Logging.Error(err)
		return
	}
	defer conn.Close()
	Logging.Debug("Connected to INNER WEBSOCKET")

	var breaker = make(chan int)
	for {
		var quitted = make(chan int)
		go func() {
			adaptor := ble.NewClientAdaptor(*addr)
			defer func() {
				if adaptor != nil {
					adaptor.Disconnect()
				}
				quitted <- 0
			}()
			err := adaptor.Connect()
			if err != nil {
				Logging.Error(err)
				return
			}
			Logging.Debug("Connected to wristband")

			/*
					=====================================
				=============================================
				============ WRITE CHARACTERISTICS ==========
				=============================================
					=====================================
			*/
			Logging.Debug("Sending Encryption.")
			bleEncryptor := security.GetBLEEncryptionSystem(*addr)

			adaptor.WithoutResponses(false)
			err = adaptor.WriteCharacteristic(DEVICE_WRISTBAND_ENCRYPTION_KEY_SETUP, bleEncryptor.GenerateMasterKey())
			if err != nil {
				Logging.Error(err)
				return
			}
			Logging.Debug("Encryption set.")

			/*
					=====================================
				=============================================
				========== SUSCRIBE CHARACTERISTICS =========
				=============================================
					=====================================
			*/

			Logging.Debug("Suscribing to IMU raw data")
			// Suscribe RAW DATA from IMU
			err = adaptor.Subscribe(DEVICE_WRISTBAND_IMU_RAW_DATA_CHARACTERISTIC, func(data []byte, err error) {
				defer func() {
					lastTimestamp = time.Now()
				}()
				//Logging.Debug("New Value coming from characteristics.")
				if err != nil {
					Logging.Error(err)
					return
				}

				decrypted_value, err := bleEncryptor.Decrypt(data)
				if err != nil {
					Logging.Error(err)
					return
				}

				jsonString := BytesU8ToJSON("imu_raw_data", decrypted_value)
				var m map[string]interface{}
				json.Unmarshal([]byte(jsonString), &m)

				ax, ay, az, gx, gy, gz := ExtractFloatValuesFromBLE(decrypted_value)
				inertialData := messaging.InertialData{
					AX: ax, AY: ay, AZ: az,
					GX: gx, GY: gy, GZ: gz,
				}

				if inertialData.HasWeirdData() {
					return
				}
				inertialData.RoundUpData()

				msg := messaging.Message{}
				msg.DataType = constants.MESSAGING_DATATYPE_DATA
				msg.AesIV = security.RandomKey()
				msg.Data, _ = json.Marshal(inertialData)

				err = conn.WriteMessage(websocket.BinaryMessage, msg.ToBytes())
				if err != nil {
					if strings.Contains(err.Error(), "reset by peer") {
						breaker <- 0
					}
					return
				}
			})
			if err != nil {
				Logging.Error(err)
				return
			}

			Logging.Debug("All setup. Loop started.")
			for {
				time.Sleep(time.Millisecond * 500)
				if time.Since(lastTimestamp) > TIMESTAMP_STOPPER {
					return
				}
			}
		}()
		select {
		case <-quitted:
			Logging.Debug("Disconnected. Will try to reconnect.")
		case <-breaker:
			//Logging.Debug("No connection with the breaker. Panicing.")
			panic("No connection with the breaker.")
		}
	}

	select {}
}

func BytesU8ToJSON(key string, values []uint8) string {
	m := make(map[string][]int)
	m[key] = make([]int, len(values))
	for i, v := range values {
		m[key][i] = int(v)
	}
	_b, _ := json.Marshal(m)
	return string(_b)
}

func BytesU16ToJSON(key string, values []uint16) string {
	m := make(map[string][]int)
	m[key] = make([]int, len(values))
	for i, v := range values {
		m[key][i] = int(v)
	}
	_b, _ := json.Marshal(m)
	return string(_b)
}
