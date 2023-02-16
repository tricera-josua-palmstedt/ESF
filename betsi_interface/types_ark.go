package betsi_interface

type BatteryPowerControl struct {
	P         float64 `unit:"W"` // Active power measured in W
	PSet      float64 `unit:"W"` // Active power setpoint in W
	PPriority uint    // Priority for active power setpoint
	Q         float64 `unit:"Var"`   // Reactive power measured in Var
	Soc       float64 `unit:"%/100"` // Actual state of charge in %/100
	SocMax    float64 `unit:"%/100"` // Maximum state of charge that can be used by the service in %/100
	SocMin    float64 `unit:"%/100"` // Minimum state of charge that can be used by the service in %/100
	PMax      float64 `unit:"W"`     // Maximum available active power in W
	PMin      float64 `unit:"W"`     // Minimum available active power in W
	ENom      float64 `unit:"Wh"`    // Nominal energy in Wh
}

type PvPowerControl struct {
	P         float64 `unit:"W"` // Active power measured in W
	PSet      float64 `unit:"W"` // Active power setpoint in W
	PPriority uint    // Priority for active power setpoint
	Q         float64 `unit:"Var"` // Reactive power measured in Var
	PMax      float64 `unit:"W"`   // Maximum available active power in W
	PMin      float64 `unit:"W"`   // Minimum available active power in W
}

type LoadPowerControl struct {
	P               float64 `unit:"W"` // Active power measured in W
	PSet            float64 `unit:"W"` // Active power setpoint in W
	PPriority       uint    // Priority for active power setpoint
	Q               float64 `unit:"Var"` // Reactive power measured in Var
	PMax            float64 `unit:"W"`   // Maximum available active power in W
	PMin            float64 `unit:"W"`   // Minimum available active power in W
	PFromGrid       float64 `unit:"W"`   // Maximum allowed active power for loading from grid in W
	PConsumptionMin float64 `unit:"W"`   // Minimum power consumption in W
}
