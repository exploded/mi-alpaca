package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

const NumOnOffSwitch = 6 // Number of on/off switches
const NumVarSwitch = 0   // Number of variable switches
const NumSwitches = 6    // Number of all switches in total

type Device struct {
	IP         string `json:"ip"`
	Token      string `json:"token"`
	Name       string `json:"name"`
	Devicetype string `json:"devicetype"`
	Number     uint32 `json:"number"`
	Uniqueid   string `json:"uniqueid"`
	Id         uint32 `json:"id"`
	Customname string `json:"customname"`
	Min        int64  `json:"min"`
	Max        int64  `json:"max"`
	Step       int64  `json:"step"`
	Canwrite   bool   `json:"canwrite"`
	Value      int64  `json:"value"`
}

type sw struct {
	Connected bool     `json:"connected"`
	Devices   []Device `json:"devices"`
}

var s = &sw{}
var sm sync.RWMutex

func MiSetInit() {
	s.misetinit()
	// Query actual state from all devices after loading settings
	s.queryAllDeviceStates()
}

func (s *sw) misetinit() {
	if !s.miLoadSettings() {
		log.Println("No settings file found, exiting")
		os.Exit(1)
	}
}

// queryAllDeviceStates queries the actual power state from all Xiaomi devices
func (s *sw) queryAllDeviceStates() {
	log.Println("Querying actual state from all devices...")
	for i := int32(0); i < int32(len(s.Devices)); i++ {
		state, err := miQueryPower(i)
		if err != nil {
			log.Printf("Warning: Failed to query device %d (%s): %v - keeping cached value", i+1, s.Devices[i].Name, err)
			continue
		}

		sm.Lock()
		if state {
			s.Devices[i].Value = 1
		} else {
			s.Devices[i].Value = 0
		}
		sm.Unlock()
		log.Printf("Device %d (%s): %v", i+1, s.Devices[i].Name, state)
	}
	s.miSaveSettings()
	log.Println("Device state query complete")
}

func (s *sw) miSaveSettings() {
	sm.Lock()
	defer sm.Unlock()

	data, err := json.MarshalIndent(&s, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("settings.json", data, 0644)
	if err != nil {
		panic(err)
	}
}

func (s *sw) miLoadSettings() bool {
	sm.Lock()
	defer sm.Unlock()
	data, err := os.ReadFile("settings.json")
	if err != nil {
		return false
	}
	return json.Unmarshal(data, &s) == nil
}

func MiGetInit() []DeviceConfiguration {
	sm.Lock()
	defer sm.Unlock()
	var val []DeviceConfiguration
	// Return only a single Switch device that contains all 6 switches
	// Use the first device's information as the base
	if len(s.Devices) > 0 {
		val = append(val, DeviceConfiguration{
			DeviceName:   "Mi Switch Controller",
			DeviceType:   "Switch",
			DeviceNumber: 1,
			UniqueID:     s.Devices[0].Uniqueid,
		})
	}
	return val
}

func MiSetName(id int32, CustomName string) error {
	if id < 0 || id >= NumSwitches {
		return errors.New("invalid device number")
	}
	s.setname(id, CustomName)
	return nil
}

func (s *sw) setname(id int32, CustomName string) {
	sm.Lock()
	s.Devices[id].Customname = CustomName
	sm.Unlock()
	s.miSaveSettings()
}

func MiSetConnect(c bool) {
	s.setconnect(c)
	// When connecting, refresh all device states from hardware
	if c {
		s.queryAllDeviceStates()
	}
}

func (s *sw) setconnect(c bool) {
	sm.Lock()
	s.Connected = c
	sm.Unlock()
	s.miSaveSettings()
}

func MiGetConnected() bool {
	return s.getconnected()
}

func (s *sw) getconnected() bool {
	sm.Lock()
	defer sm.Unlock()
	return s.Connected
}

func MiGetName(id int32) string {
	sm.Lock()
	defer sm.Unlock()
	if s.Devices[id].Customname != "" {
		return s.Devices[id].Customname
	}
	return s.Devices[id].Name
}

func MiGetType(id int32) string {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Devicetype
}

func MiGetNumber(id uint32) uint32 {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Number
}

func MiGetUniqueID(id int32) string {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Uniqueid
}

func MiGetOnOff(id int32) (bool, error) {
	sm.Lock()
	defer sm.Unlock()
	if s.Devices[id].Max > 1 {
		return false, errors.New("device is not just an on/off switch")
	}
	return s.Devices[id].Value != 0, nil
}

func MiGetValue(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Value
}

func MiGetMax(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Max
}

func MiGetMin(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Min
}

func MiGetStep(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Step
}

func MiGetCanWrite(id int32) bool {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices[id].Canwrite
}

// MiSetOnOff sends the command to turn the switches on or off (id is from 1 to 6)
func MiSetOnOff(id int32, state bool) error {
	if id < 1 || id > NumOnOffSwitch {
		return errors.New("invalid switch number")
	}

	if err := miOnOff(id, state); err != nil {
		return err
	}
	return s.setonoff(id, state)
}

func (s *sw) setonoff(id int32, state bool) error {
	sm.Lock()
	if state {
		s.Devices[id].Value = 1
	} else {
		s.Devices[id].Value = 0
	}
	sm.Unlock()
	s.miSaveSettings()
	log.Printf("Set switch %d to %v", id+1, state)
	return nil
}

func MiGetDevices() []Device {
	sm.Lock()
	defer sm.Unlock()
	return s.Devices
}
