from services_build6 import services
#from betsi_interface_build import betsi_interface

bpc = services.BatteryPowerControl()
ppc = services.PvPowerControl()
lpc = services.LoadPowerControl()
ppoc = services.Poc()



#data_in = services.ServiceData(Poc=ppoc, Battery=bpc, Pv=ppc, Load=lpc)
data_in = services.ServiceData()

data_in.Poc = ppoc
data_in.Battery = bpc
data_in.Pv = ppc
data_in.Load = lpc

po1 = services.PO1_B1P1C1()

result = po1.Execute(data_in)

print(result)
