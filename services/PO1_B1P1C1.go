package services

import (
	"fmt"
	"math"
)

type esf_service interface {
	Execute()
}

type PO1_B1P1C1 struct {
	Params Parameter_PO1_B1P1C1
	state  socState
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

type socState struct {
	socBelowTarget  bool
	socBelowReserve bool
}

type ServiceInit struct {
	Name      string `json:"name" yaml:"name"`
	Parameter Parameter_PO1_B1P1C1
}

type ServiceData struct {
	Poc       Poc
	Batteries []*BatteryPowerControl
	Pvs       []*PvPowerControl
	Loads     []*LoadPowerControl
}

// type ServiceData struct {
// 	Poc     Poc
// 	Battery BatteryPowerControl
// 	Pv      PvPowerControl
// 	Load    LoadPowerControl
// }

func (s *PO1_B1P1C1) Execute(input ServiceData) (output ServiceData) {

	// Still missing:
	// - check whether input data is valid
	//   (number of peripheries, etc.)
	// - iterate over PowerControls?

	data := input

	var bpc *BatteryPowerControl
	var ppc *PvPowerControl
	var lpc *LoadPowerControl

	var poc Poc

	if len(data.Batteries) == 1 && len(data.Pvs) == 1 && len(data.Loads) == 1 {
		bpc = data.Batteries[0]
		ppc = data.Pvs[0]
		lpc = data.Loads[0]
		// Also check if there is POC
		poc = data.Poc
	} else {
		// Raise error. How?
		fmt.Printf("PO1_B1P1C1 needs exactly 1 of each powercontrols!")
	}

	pResidual := poc.P - bpc.P

	// Hysteresis definition
	if bpc.Soc >= s.Params.SocTarget+s.Params.DsocTargetHyst {
		s.state.socBelowTarget = false
	} else if bpc.Soc < s.Params.SocTarget {
		s.state.socBelowTarget = true
	}

	if bpc.Soc >= s.Params.SocReserve {
		s.state.socBelowReserve = false
	} else if bpc.Soc < s.Params.SocReserve {
		s.state.socBelowReserve = true
	}

	var pPocMaxTarget float64
	var pPocMinTarget float64

	if !s.state.socBelowTarget {
		pPocMaxTarget = 0
		pPocMinTarget = 0
	} else if !s.state.socBelowReserve {
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
	if !s.state.socBelowTarget {
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

	ppc.PPriority = 1
	lpc.PPriority = 2
	poc.PPriority = 3
	bpc.PPriority = 4

	output = ServiceData{}

	output.Poc = poc
	output.Batteries[0] = bpc
	output.Pvs[0] = ppc
	output.Loads[0] = lpc

	return
}
