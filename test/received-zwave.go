package main

import (
	"flag"
	"github.com/kevinchapron/BasicLogger/Logging"
	"github.com/stampzilla/gozwave"
	"github.com/stampzilla/gozwave/events"
	"strconv"
	"time"
)

func main() {
	Logging.SetLevel(Logging.DEBUG)
	Logging.Debug("Starting ...")

	var port string
	flag.StringVar(&port, "port", "/dev/ttyACM0", "SerialAPI Communication Port")
	flag.Parse()

	z, err := gozwave.Connect(port, "")
	if err != nil {
		Logging.Error(err)
		return
	}

	go keepPrintingNodes(z)

	for {
		select {
		case event := <-z.GetNextEvent():
			switch e := event.(type) {
			case events.NodeDiscoverd:
				Logging.Info("New Node discovered #" + strconv.Itoa(e.Address))
			case events.NodeUpdated:
				Logging.Info("Node #" + strconv.Itoa(e.Address) + " updated")
			}
		}
	}

}
func keepPrintingNodes(z *gozwave.Controller) {
	for {
		<-time.After(time.Second * 5)
		for _, node := range z.Nodes.All() {
			idNode := node.Id

			availableCommands := node.CommandClasses
			for index, availableCommand := range availableCommands {
				if len(availableCommand.ID.String()) == 0 {
					continue
				}
				Logging.Debug("NODE #"+strconv.Itoa(idNode)+": Command #"+strconv.Itoa(index)+" ("+strconv.Itoa(int(availableCommand.ID))+"):", availableCommand.ID.String())
			}

			Logging.Debug("Trying to put sensor OFF...")
			node.Off()
			time.Sleep(time.Second * 5)
			Logging.Debug("Trying to put sensor ON ...")
			node.On()

		}
	}
}
