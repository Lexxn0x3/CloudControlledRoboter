package main

type Motor struct {
	Motor1 int8 `json:"motor1" mapstructure:"motor1"`
	Motor2 int8 `json:"motor2" mapstructure:"motor2"`
	Motor3 int8 `json:"motor3" mapstructure:"motor3"`
	Motor4 int8 `json:"motor4" mapstructure:"motor4"`
}

type Lightbar struct {
	Mode   bool `json:"mode" mapstructure:"mode"`
	LedID  byte `json:"ledid" mapstructure:"ledid"`
	R      byte `json:"red" mapstructure:"red"`
	G      byte `json:"green" mapstructure:"green"`
	B      byte `json:"blue" mapstructure:"blue"`
	Effect byte `json:"effect" mapstructure:"effect"`
	Speed  byte `json:"speed" mapstructure:"speed"`
	Parm   byte `json:"parm" mapstructure:"parm"`
}

type Buzzer struct {
	Duration int `json:"buzzer" mapstructure:"buzzer"`
}

type Laser struct {
	Status bool `json:"laser" mapstructure:"laser"`
}
