package main

import (
	"encoding/binary"
	"fmt"
	"github.com/goburrow/modbus"
	"math"
	"time"
)

const (
	LineToNeutralVolts                   uint16 = 0x00
	Current                              uint16 = 0x06
	ActivePower                          uint16 = 0x0C
	ApparentPower                        uint16 = 0x12
	ReactivePower                        uint16 = 0x18
	PowerFactor                          uint16 = 0x1E
	PhaseAngle                           uint16 = 0x24
	Frequency                            uint16 = 0x46
	ImportActiveEnergy                   uint16 = 0x48
	ExportActiveEnergy                   uint16 = 0x4A
	ImportReactiveEnergy                 uint16 = 0x4C
	ExportReactiveEnergy                 uint16 = 0x4E
	TotalSystemPowerDemand               uint16 = 0x54
	MaximumTotalSystemPowerDemand        uint16 = 0x56
	CurrentSystemPositivePowerDemand     uint16 = 0x58
	MaximumSystemPositivePowerDemand     uint16 = 0x5A
	CurrentSystemReversePowerDemand      uint16 = 0x5C
	MaximumSystemReversePowerDemand      uint16 = 0x5E
	CurrentDemand                        uint16 = 0x0102
	MaximumCurrentDemand                 uint16 = 0x0108
	TotalActiveEnergy                    uint16 = 0x0156
	TotalReactiveEnergy                  uint16 = 0x0158
	CurrentResettableTotalActiveEnergy   uint16 = 0x0180
	CurrentResettableTotalReactiveEnergy uint16 = 0x0182
)

type Eastron struct {
	client modbus.Client
}

func NewEastron(c modbus.Client) *Eastron {
	e := Eastron{
		client: c,
	}
	return &e
}

func main() {

	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "E"
	handler.StopBits = 1
	handler.SlaveId = 2
	handler.Timeout = 2 * time.Second

	err := handler.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer handler.Close()

	client := modbus.NewClient(handler)

	SDM120 := NewEastron(client)

	fmt.Println(SDM120.Read(LineToNeutralVolts))
	fmt.Println(SDM120.Read(Current))
	fmt.Println(SDM120.Read(ActivePower))
	fmt.Println(SDM120.Read(ApparentPower))
	fmt.Println(SDM120.Read(ReactivePower))

	fmt.Println(SDM120.Read(PowerFactor))
}

func (e *Eastron) Read(addr uint16) float32 {
	results, err := e.client.ReadInputRegisters(addr, 2)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	bits := binary.BigEndian.Uint32(results)
	float := math.Float32frombits(bits)
	return float
}
