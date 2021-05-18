package websockets

import "time"

var listSensors ListSensor

type ListSensor []*WebSocketClient

func (l *ListSensor) UpdateSensor(sensor *WebSocketClient) {
	if listSensors == nil {
		listSensors = []*WebSocketClient{}
	}
	exists := false
	for _, sens := range *l {
		if sens.Name == sensor.Name && sens.Protocol == sensor.Protocol {
			exists = true
			break
		}
	}
	if !exists {
		listSensors = append(listSensors, sensor)
	}
	for index, sens := range *l {
		if sens.Name == sensor.Name && sens.Protocol == sensor.Protocol {
			(*l)[index].lastMessageTime = time.Now()
			break
		}
	}
}

func (l *ListSensor) ListAllSensors(filter *string) []*WebSocketClient {
	var list []*WebSocketClient
	for _, client := range *l {
		if filter == nil || client.Protocol == *filter {
			list = append(list, client)
		}
	}
	return list
}

func GetListSensor() *ListSensor {
	if listSensors == nil {
		listSensors = []*WebSocketClient{}
	}
	return &listSensors
}
