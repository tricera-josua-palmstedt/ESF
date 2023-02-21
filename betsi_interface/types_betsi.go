package betsi_interface

import "github.com/tricera-josua-palmstedt/ESF/services"

type BpcSlice []services.BatteryPowerControl
type PpcSlice []services.PvPowerControl
type LpcSlice []services.LoadPowerControl

// type BpcSlicePtr []*services.BatteryPowerControl
// type PpcSlicePtr []*services.PvPowerControl
// type LpcSicePtr []*services.LoadPowerControl

func (this BpcSlice) BpcSliceAppend(bpc services.BatteryPowerControl) BpcSlice {
	return append(this, bpc)
}

func (this PpcSlice) PpcSliceAppend(ppc services.PvPowerControl) PpcSlice {
	return append(this, ppc)
}

func (this LpcSlice) LpcSliceAppend(lpc services.LoadPowerControl) LpcSlice {
	return append(this, lpc)
}

// func (this BpcSlicePtr) BpcSlicePtrAppend(bpc *services.BatteryPowerControl) BpcSlicePtr {
// 	return append(this, bpc)
// }
