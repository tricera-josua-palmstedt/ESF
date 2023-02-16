package services

import "math"

type Limit struct {
	PMax            float64
	PMin            float64
	SMax            float64
	PFromGrid       float64
	PConsumptionMin float64
}

type Poc struct {
	P         float64 `unit:"W"` // Active power measured
	PSet      float64 `unit:"W"` // Active power setpoint
	PPriority uint    // Priority for active power setpoint
	Q         float64 `unit:"Var"` // Reactive power measured
	Limit     Limit
	// PMax      float64 `unit:"W"`   // Maximum available active power
	// PMin      float64 `unit:"W"`   // Minimum available active power
	// SMax      float64 `unit:"VA"`  // Maximum available apparent power
}

type BatteryPowerControl struct {
	P         float64 `unit:"W"` // Active power measured in W
	PSet      float64 `unit:"W"` // Active power setpoint in W
	PPriority uint    // Priority for active power setpoint
	Q         float64 `unit:"Var"`   // Reactive power measured in Var
	Soc       float64 `unit:"%/100"` // Actual state of charge in %/100
	SocMax    float64 `unit:"%/100"` // Maximum state of charge that can be used by the service in %/100
	SocMin    float64 `unit:"%/100"` // Minimum state of charge that can be used by the service in %/100
	ENom      float64 `unit:"Wh"`    // Nominal energy in Wh
	Limit     Limit
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

func limit(value float64, minimum float64, maximum float64) float64 {
	return math.Min(maximum, math.Max(value, minimum))
}
