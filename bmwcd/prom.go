/*
MIT License

Copyright (c) 2020 tcassaert

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package bmwcd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tidwall/gjson"
)

var (
	mileage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_mileage",
			Help: "The current mileage of the car",
		},
	)
	doorLockState = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_door_lock_state",
			Help: "Doors unlocked (0) or closed (1)",
		},
	)
	doorDriverFront = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_door_driver_front_state",
			Help: "Door open (0) or closed (1)",
		},
	)
	doorDriverRear = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_door_driver_rear_state",
			Help: "Door open (0) or closed (1)",
		},
	)
	doorPassengerFront = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_door_passenger_front_state",
			Help: "Door open (0) or closed (1)",
		},
	)
	doorPassengerRear = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_door_passenger_rear_state",
			Help: "Door open (0) or closed (1)",
		},
	)
	windowDriverFront = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_window_driver_front_state",
			Help: "Window open (0) or closed (1)",
		},
	)
	windowDriverRear = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_window_driver_rear_state",
			Help: "Window open (0) or closed (1)",
		},
	)
	windowPassengerFront = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_window_passenger_front_state",
			Help: "Window open (0) or closed (1)",
		},
	)
	windowPassengerRear = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_window_passenger_rear_state",
			Help: "Window open (0) or closed (1)",
		},
	)
	trunk = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_trunk_state",
			Help: "Trunk open (0) or closed (1)",
		},
	)
	hood = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_hood_state",
			Help: "Hood open (0) or closed (1)",
		},
	)
	remainingFuel = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_remaining_fuel",
			Help: "Remaining liters of fuel in the tank",
		},
	)
	remainingRangeElectric = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_remaining_electric_range",
			Help: "Remaining kilometers of electric range",
		},
	)
	remainingRangeHybrid = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_remaining_hybrid_range",
			Help: "Remaining kilometers of hybrid range",
		},
	)
	chargeLevel = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_charge_level",
			Help: "Charge percentage",
		},
	)
	oilCheckRemainingMileage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_oil_cbs_remaining_mileage",
			Help: "Remaining kilometers before oil cbs",
		},
	)
	vehicleCheckRemainingMileage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_vehicle_check_cbs_remaining_mileage",
			Help: "Remaining kilometers before vehicle check cbs",
		},
	)
	oilCheckDueDate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_oil_check_cbs_due_date",
			Help: "Oil cbs due by date",
		},
	)
	vehicleCheckDueDate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_vehicle_check_cbs_due_date",
			Help: "Vehicle check cbs due by date",
		},
	)
	brakeFluidCheckDueDate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_brake_fluid_check_cbs_due_date",
			Help: "Brake fluid check cbs due by date",
		},
	)
)

func jsonToProm(status string) {
	// Current mileage of the car
	mileage.Set(gjson.Get(status, "mileage").Float())

	// Kilometers untill service
	oilCheckRemainingMileageValue := gjson.Get(status, "cbsData.#(cbsType==\"OIL\").cbsRemainingMileage").Float()
	oilCheckRemainingMileage.Set(oilCheckRemainingMileageValue)

	vehicleCheckRemainingMileageValue := gjson.Get(status, "cbsData.#(cbsType==\"VEHICLE_CHECK\").cbsRemainingMileage").Float()
	vehicleCheckRemainingMileage.Set(vehicleCheckRemainingMileageValue)

	// Due date service
	oilCheckDueDateValue := gjson.Get(status, "cbsData.#(cbsType==\"OIL\").cbsDueDate").String()
	oilCheckDueDate.Set(convertToEpoch(oilCheckDueDateValue))

	brakeFluidCheckDueDateValue := gjson.Get(status, "cbsData.#(cbsType==\"BRAKE_FLUID\").cbsDueDate").String()
	brakeFluidCheckDueDate.Set(convertToEpoch(brakeFluidCheckDueDateValue))

	vehicleCheckDueDateValue := gjson.Get(status, "cbsData.#(cbsType==\"VEHICLE_CHECK\").cbsDueDate").String()
	vehicleCheckDueDate.Set(convertToEpoch(vehicleCheckDueDateValue))

	// State of the door lock
	doorLockStateString := gjson.Get(status, "doorLockState").String()
	if doorLockStateString == "UNLOCKED" {
		doorLockState.Set(0)
	} else {
		doorLockState.Set(1)
	}
	// State of the doors
	doorDriverFrontValue := gjson.Get(status, "doorDriverFront").String()
	if doorDriverFrontValue == "OPEN" {
		doorDriverFront.Set(0)
	} else {
		doorDriverFront.Set(1)
	}

	doorDriverRearValue := gjson.Get(status, "doorDriverRear").String()
	if doorDriverRearValue == "OPEN" {
		doorDriverRear.Set(0)
	} else {
		doorDriverRear.Set(1)
	}

	doorPassengerFrontValue := gjson.Get(status, "doorPassengerFront").String()
	if doorPassengerFrontValue == "OPEN" {
		doorPassengerFront.Set(0)
	} else {
		doorPassengerFront.Set(1)
	}

	doorPassengerRearValue := gjson.Get(status, "doorPassengerRear").String()
	if doorPassengerRearValue == "OPEN" {
		doorPassengerRear.Set(0)
	} else {
		doorPassengerRear.Set(1)
	}

	// State of the windows
	windowDriverFrontValue := gjson.Get(status, "windowDriverFront").String()
	if windowDriverFrontValue == "OPEN" {
		windowDriverFront.Set(0)
	} else {
		windowDriverFront.Set(1)
	}

	windowDriverRearValue := gjson.Get(status, "windowDriverRear").String()
	if windowDriverRearValue == "OPEN" {
		windowDriverRear.Set(0)
	} else {
		windowDriverRear.Set(1)
	}

	windowPassengerFrontValue := gjson.Get(status, "windowPassengerFront").String()
	if windowPassengerFrontValue == "OPEN" {
		windowPassengerFront.Set(0)
	} else {
		windowPassengerFront.Set(1)
	}

	windowPassengerRearValue := gjson.Get(status, "windowPassengerRear").String()
	if windowPassengerRearValue == "OPEN" {
		windowPassengerRear.Set(0)
	} else {
		windowPassengerRear.Set(1)
	}

	// State of the hood
	hoodValue := gjson.Get(status, "hood").String()
	if hoodValue == "OPEN" {
		hood.Set(0)
	} else {
		hood.Set(1)
	}

	// State of the trunk
	trunkValue := gjson.Get(status, "trunk").String()
	if trunkValue == "OPEN" {
		trunk.Set(0)
	} else {
		trunk.Set(1)
	}

	// Remaining fuel
	remainingFuelValue := gjson.Get(status, "remainingFuel").Float()
	remainingFuel.Set(remainingFuelValue)

	// Remaining electric range
	remainingRangeElectricValue := gjson.Get(status, "remainingRangeElectric").Float()
	remainingRangeElectric.Set(remainingRangeElectricValue)

	// Remaining hybrid range
	remainingRangeHybridValue := gjson.Get(status, "remainingRangeFuel").Float()
	remainingRangeHybrid.Set(remainingRangeHybridValue)

	// Charge level
	chargeLevelValue := gjson.Get(status, "chargingLevelHv").Float()
	chargeLevel.Set(chargeLevelValue)

}
