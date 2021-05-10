package main

import (
	"encoding/binary"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/kevinchapron/FSHK-final/constants"
	"github.com/kevinchapron/FSHK-final/messaging"
	"github.com/kevinchapron/FSHK-final/security"
	"gobot.io/x/gobot/platforms/ble"
	"net/url"
	"strconv"
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

const TIMESTAMP_STOPPER = time.Second * 5

func ExtractFloatValuesFromBLE(b []byte) (float64, float64, float64, float64, float64, float64) {
	ax := (float64(binary.LittleEndian.Uint16(b[0:2])) / 1000) - 4
	ay := (float64(binary.LittleEndian.Uint16(b[2:4])) / 1000) - 4
	az := (float64(binary.LittleEndian.Uint16(b[4:6])) / 1000) - 4

	gx := (float64(binary.BigEndian.Uint32(append([]byte{0x00}, b[6:9]...))) / 1000) - 2000
	gy := (float64(binary.BigEndian.Uint32(append([]byte{0x00}, b[9:12]...))) / 1000) - 2000
	gz := (float64(binary.BigEndian.Uint32(append([]byte{0x00}, b[12:15]...))) / 1000) - 2000
	return ax, ay, az, gx, gy, gz
}

type InertialData struct {
	AX float64
	AY float64
	AZ float64

	GX float64
	GY float64
	GZ float64
}

func main() {
	Logging.SetLevel(Logging.DEBUG)
	lastTimestamp := time.Now()

	// Connecting to WS
	ipLIPSHOK := "192.168.5.226"
	u := url.URL{Scheme: "ws", Host: ipLIPSHOK + ":" + strconv.Itoa(constants.WEBSOCKET_INNER_BLE_PORT), Path: constants.WEBSOCKET_INNER_BLE_PATH}
	Logging.Debug("Trying to connect to inner WS : " + u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		Logging.Error(err)
		return
	}
	defer conn.Close()
	Logging.Debug("Connected to INNER WEBSOCKET")

	addr := "EF:3A:4E:89:07:0D"
	adaptor := ble.NewClientAdaptor(addr)
	for {
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
		bleEncryptor := security.GetBLEEncryptionSystem(addr)

		adaptor.WithoutResponses(false)
		err = adaptor.WriteCharacteristic(DEVICE_WRISTBAND_ENCRYPTION_KEY_SETUP, bleEncryptor.GenerateMasterKey())
		if err != nil {
			Logging.Error(err)
			return
		}
		Logging.Debug("Encryption set.")

		// Suscribe RAW DATA from IMU
		err = adaptor.Subscribe(DEVICE_WRISTBAND_IMU_RAW_DATA_CHARACTERISTIC, func(data []byte, err error) {
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
			inertialData := InertialData{
				AX: ax, AY: ay, AZ: az,
				GX: gx, GY: gy, GZ: gz,
			}

			msg := messaging.Message{}
			msg.DataType = constants.MESSAGING_DATATYPE_DATA
			msg.AesIV = security.RandomKey()
			msg.Data, _ = json.Marshal(inertialData)

			err = conn.WriteMessage(websocket.BinaryMessage, msg.ToBytes())
			if err != nil {
				Logging.Error(err)
				return
			}

		})
		if err != nil {
			Logging.Error(err)
			return
		}

		for {
			if time.Since(lastTimestamp) > TIMESTAMP_STOPPER {
				Logging.Debug("Disconnected. Will try to reconnect.")
				break
			}
			time.Sleep(time.Second)
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
