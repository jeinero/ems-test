package main

import (
	"fmt"
	"time"
)

// ESS variables
type ESS struct {
	Pess, Pmaxch, Pmaxdisch, Eess float64
}

func main() {
	EMS(-500) // PmaxSite is set to -500 kW
}

// Get ess measure
func GetEssMeasure() *ESS {
	//TODO: hardware implementation
	return &ESS{
		Pess:      120,
		Pmaxch:    -200,
		Pmaxdisch: 100,
		Eess:      500,
	} // Pess, Pmaxch, Pmaxdisch in kW and Eess in KWh
}

// Get Ppv measure
func GetPvMeasure() float64 {
	//TODO: hardware implementation
	return 200000 // Ppv in W
}

// Get Ppoc measure
func GetPocMeterMeasure() float64 {
	//TODO: hardware implementation
	return -600 // Ppoc in kW
}

// Set ESS instruction
func SetEssSetpoint(setpointPess float64) {
	//TODO: hardware implementation
	fmt.Printf("Set ESS setpoint to %f kW\n", setpointPess)
}

// EMS logic
func EMS(PmaxSite float64) {
	for {
		// Get measurements
		essMeasure := GetEssMeasure()
		Ppv := GetPvMeasure() / 1000 // W to KW
		Ppoc := GetPocMeterMeasure()
		fmt.Printf("Ppoc = %fKW\n", Ppoc)

		// Calculate Pload
		Pload := Ppoc - essMeasure.Pess - Ppv
		fmt.Printf("Pload = %fKW\n", Pload)

		// If Ppoc is more than the maximum limit
		if PmaxSite > Ppoc {
			// Calculate the power excess
			powerExcess := PmaxSite - Ppoc
			fmt.Printf("powerExcess = %fKW\n", powerExcess)

			// Check if the ESS has enough energy to discharge
			if essMeasure.Eess > 0 {
				// If excess power is less than the ESS's maximum discharge capacity
				if powerExcess < essMeasure.Pmaxdisch {
					// Set the ESS to discharge the excess power
					SetEssSetpoint(powerExcess)
				} else {
					// If not, set the ESS to discharge at its maximum capacity
					SetEssSetpoint(essMeasure.Pmaxdisch)
				}
			} else {
				// If the ESS has no more energy to discharge
				SetEssSetpoint(0)
			}
		} else if Ppoc == PmaxSite {
			// If Ppoc is equal to the maximum limit, do nothing
			SetEssSetpoint(0)
		} else {
			// If Ppoc is less than the maximum limit

			// Calculate the power deficit
			powerDeficit := PmaxSite - Ppoc
			fmt.Printf("powerDeficit = %fKW\n", powerDeficit)

			// If the power deficit is less than the ESS's maximum charge capacity
			if powerDeficit > essMeasure.Pmaxch {
				// Set the ESS to charge the power deficit
				SetEssSetpoint(powerDeficit)
			} else {
				// If not, set the ESS to charge at its maximum capacity
				SetEssSetpoint(essMeasure.Pmaxch)
			}
		}

		// wait before the next control cycle
		time.Sleep(1 * time.Second)
	}
}
