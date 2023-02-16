package services

import (
	"math"
)

type esf_service interface {
	Execute()
}

type PO1_B1P1C1 struct {
	Params Parameter_PO1_B1P1C1
	state  SocState
}

type Parameter_PO1_B1P1C1 struct {
	SocTarget      float64
	SocReserve     float64
	DsocTargetHyst float64
	//DsocReserveHyst			  float64
	//PCutInfeed				  float64
	PCutConsumption           float64
	PReserve                  float64
	BatInjectionToGridAllowed bool
	DpPvControl               float64
}

type SocState struct {
	SocBelowTarget  bool
	SocBelowReserve bool
}

type ServiceInit struct {
	Name      string `json:"name" yaml:"name"`
	Parameter Parameter_PO1_B1P1C1
}

//	type ServiceData struct {
//		poc       Poc
//		Batteries []*energy.BatteryPowerControl
//		Pvs       []*energy.PvPowerControl
//		Loads     []*energy.LoadPowerControl
//	}
type ServiceData struct {
	Poc     Poc
	Battery BatteryPowerControl
	Pv      PvPowerControl
	Load    LoadPowerControl
}

func (s *PO1_B1P1C1) Execute(input ServiceData) (output ServiceData) {

	// Still missing:
	// - check whether input data is valid
	//   (number of peripheries, etc.)
	// - iterate over PowerControls?

	data := input

	poc := data.Poc

	bpc := data.Battery
	ppc := data.Pv
	lpc := data.Load

	pResidual := poc.P - bpc.P

	// Hysteresis definition
	if bpc.Soc >= s.Params.SocTarget+s.Params.DsocTargetHyst {
		s.state.SocBelowTarget = false
	} else if bpc.Soc < s.Params.SocTarget {
		s.state.SocBelowTarget = true
	}

	if bpc.Soc >= s.Params.SocReserve {
		s.state.SocBelowReserve = false
	} else if bpc.Soc < s.Params.SocReserve {
		s.state.SocBelowReserve = true
	}

	var pPocMaxTarget float64
	var pPocMinTarget float64

	if !s.state.SocBelowTarget {
		pPocMaxTarget = 0
		pPocMinTarget = 0
	} else if !s.state.SocBelowReserve {
		pPocMaxTarget = 0
		pPocMinTarget = s.Params.PCutConsumption
	} else {
		pPocMaxTarget = pResidual + s.Params.PReserve
		pPocMaxTarget = limit(pPocMaxTarget, s.Params.PCutConsumption, 0)
		pPocMinTarget = s.Params.PCutConsumption
	}

	var pBatSet float64

	if pResidual > pPocMaxTarget {
		pBatSet = pPocMaxTarget - pResidual
	} else if pResidual < pPocMinTarget {
		pBatSet = pPocMinTarget - pResidual
	} else {
		pBatSet = 0
	}

	pBatSet = limit(pBatSet, bpc.Limit.PMin, bpc.Limit.PMax)

	//weiter mit enfluri
	var pBatMaxEnfluri float64
	var pBatMinEnfluri float64

	if !s.Params.BatInjectionToGridAllowed {
		pBatMaxEnfluri = math.Max(-pResidual, 0)
		pBatSet = math.Min(pBatSet, pBatMaxEnfluri)
	} else {
		pBatMinEnfluri = math.Min(-pResidual, 0)
		pBatSet = math.Max(pBatSet, pBatMinEnfluri)
	}

	var pPvSet float64
	pPvSet = ppc.P + s.Params.DpPvControl
	pPvSet = limit(pPvSet, ppc.PMin, ppc.PMax)

	var pClSet float64
	if !s.state.SocBelowTarget {
		pClSet = lpc.PFromGrid - ppc.P
	} else {
		pClSet = lpc.PFromGrid - ppc.P - bpc.Limit.PMax
	}

	var pPocSet float64
	pPocSet = bpc.Limit.PMax + ppc.PMax + lpc.PMax

	poc.PSet = pPocSet
	bpc.PSet = pBatSet
	ppc.PSet = pPvSet
	lpc.PSet = pClSet

	//s.storageVars.pResidual = pResidual

	bpc.PPriority = 2
	ppc.PPriority = 1
	lpc.PPriority = 3
	poc.PPriority = 4

	output = ServiceData{}

	output.Poc = poc
	output.Battery = bpc
	output.Pv = ppc
	output.Load = lpc

	return
}
