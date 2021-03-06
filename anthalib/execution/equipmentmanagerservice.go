// equipmentmanagerservice/equipmentmanagerservice.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK

package execution

import "github.com/antha-lang/antha/anthalib/liquidhandling"

// holds channels for communicating with the equipment manager
type EquipmentManagerService struct {
	RequestsIn       chan EquipmentManagerRequest
	RequestsOut      map[string]chan EquipmentManagerRequest
	devicelist       map[string][]string
	deviceproperties map[string]*liquidhandling.LHProperties
	devicequeues     map[string][]*DeviceManager
}

// get properties for a device
func (ems *EquipmentManagerService) GetEquipmentProperties(deviceclass string) interface{} {
	return ems.deviceproperties[ems.devicelist[deviceclass][0]]
}

// bit of a short-term fix

//returns a list of devices known to the system
func (ems *EquipmentManagerService) GetDeviceListByClass(class string) []string {
	return ems.devicelist[class]
}

// get properties files describing the handlers themselves
func (ems *EquipmentManagerService) GetLiquidHandlerProperties(devname string) *liquidhandling.LHProperties {
	return ems.deviceproperties[devname]
}

// ask for some equipment
func (ems *EquipmentManagerService) RequestEquipment(rin EquipmentManagerRequest) chan EquipmentManagerRequest {
	ems.RequestsIn <- rin
	ch := make(chan EquipmentManagerRequest)
	ems.RequestsOut[rin["ID"].(string)] = ch
	return ch
}

// initialize the equipment manager
// needs to read config from somewhere
func (ems *EquipmentManagerService) Init() {
	ems.RequestsIn = make(chan EquipmentManagerRequest, 5)
	ems.RequestsOut = make(map[string]chan EquipmentManagerRequest, 100)

	ems.devicelist = make(map[string][]string)

	ems.devicelist["liquidhandler"] = make([]string, 1)
	ems.devicelist["liquidhandler"][0] = "ALiquidHandler"

	ems.deviceproperties = make(map[string]*liquidhandling.LHProperties)
	ems.deviceproperties["ALiquidHandler"] = makepropertiesbodge()

	go func() {
		equipmentmanagerDaemon(ems)
	}()
}

func makepropertiesbodge() *liquidhandling.LHProperties {
	// make a liquid handling structure

	lhp := liquidhandling.NewLHProperties(12, "ALiquidHandler", "ACMEliquidhandlers", "discrete", "disposable", []string{"plate"})

	// I suspect this might need to be in the constructor
	// or at least wrapped into a factory method

	lhp.Tip_preferences = []int{1, 5, 3}
	lhp.Input_preferences = []int{10, 11, 12}
	lhp.Output_preferences = []int{7, 8, 9, 2, 4}

	// need to add some configs

	hvconfig := liquidhandling.NewLHParameter("HVconfig", 10, 250, "ul")

	cnfvol := lhp.Cnfvol
	cnfvol[0] = hvconfig
	lhp.Cnfvol = cnfvol

	lhp.CurrConf = hvconfig

	return lhp
}

// Daemon for passing requests through to the service
// the new pattern is:
// request comes in, channel comes out. manager stores channel
// when request is serviced the output channel is retrieved and
// fed the output
func equipmentmanagerDaemon(ems *EquipmentManagerService) {
	for {
		rin := <-ems.RequestsIn
		rin = handleRequest(rin)

		id := rin["ID"].(string)

		rout, ok := ems.RequestsOut[id]

		if !ok {
			panic("No corresponding output channel for request")
		}
		rout <- rin
	}
}

func MakeDeviceRequest(devicetype string) EquipmentManagerRequest {
	emr := NewEquipmentManagerRequest()

	emr["requestType"] = "DeviceRequest"
	emr["deviceType"] = devicetype

	return emr
}

func handleRequest(emr EquipmentManagerRequest) EquipmentManagerRequest {
	switch emr["requestType"] {
	case "DeviceRequest":

	}
	return emr
}
