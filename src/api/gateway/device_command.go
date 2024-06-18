package gateway

const PressPowerSwitchOpcode int = 1
const PressResetSwitchOpcode int = 2
const HardPowerOffOpcode int = 3

type CommandMessage struct {
	Opcode int `json:"op"`
}
