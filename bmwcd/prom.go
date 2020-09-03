package bmwcd

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/tidwall/gjson"
)

var (
	brakeFluidCheckDueDate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_brake_fluid_check_cbs_due_date",
			Help: "Brake fluid check cbs due by date",
		},
	)
	chargingStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_charging_status",
			Help: "Not charging (0), charging (1), fully charged (2)",
		},
	)
	chargeLevel = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_charge_level",
			Help: "Charge percentage",
		},
	)
	connectionStatus = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_connection_status",
			Help: "Charging cable connected (1) or disconnected (0)",
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
	hood = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_hood_state",
			Help: "Hood open (0) or closed (1)",
		},
	)
	mileage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_mileage",
			Help: "The current mileage of the car",
		},
	)
	oilCheckDueDate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_oil_check_cbs_due_date",
			Help: "Oil cbs due by date",
		},
	)
	oilCheckRemainingMileage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_oil_cbs_remaining_mileage",
			Help: "Remaining kilometers before oil cbs",
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
	trunk = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_trunk_state",
			Help: "Trunk open (0) or closed (1)",
		},
	)
	vehicleCheckDueDate = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_vehicle_check_cbs_due_date",
			Help: "Vehicle check cbs due by date",
		},
	)
	vehicleCheckRemainingMileage = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "bmwcd_vehicle_check_cbs_remaining_mileage",
			Help: "Remaining kilometers before vehicle check cbs",
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
)

func jsonToProm(status string) {

	brakeFluidCheckDueDateValue := gjson.Get(status, "cbsData.#(cbsType==\"BRAKE_FLUID\").cbsDueDate").String()
	brakeFluidCheckDueDate.Set(convertToEpoch(brakeFluidCheckDueDateValue))

	chargeLevelValue := gjson.Get(status, "chargingLevelHv").Float()
	chargeLevel.Set(chargeLevelValue)

	chargingStatusValue := gjson.Get(status, "chargingStatus").String()
	if chargingStatusValue == "FINISHED_FULLY_CHARGED" {
		chargingStatus.Set(2)
	} else if chargingStatusValue == "CHARGING" {
		chargingStatus.Set(1)
	} else {
		chargingStatus.Set(0)
	}

	connectionStatusValue := gjson.Get(status, "connectionStatus").String()
	if connectionStatusValue == "CONNECTED" {
		connectionStatus.Set(1)
	} else {
		connectionStatus.Set(0)
	}

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

	doorLockStateValue := gjson.Get(status, "doorLockState").String()
	if doorLockStateValue == "UNLOCKED" {
		doorLockState.Set(0)
	} else {
		doorLockState.Set(1)
	}

	hoodValue := gjson.Get(status, "hood").String()
	if hoodValue == "OPEN" {
		hood.Set(0)
	} else {
		hood.Set(1)
	}

	mileage.Set(gjson.Get(status, "mileage").Float())

	oilCheckDueDateValue := gjson.Get(status, "cbsData.#(cbsType==\"OIL\").cbsDueDate").String()
	oilCheckDueDate.Set(convertToEpoch(oilCheckDueDateValue))

	oilCheckRemainingMileageValue := gjson.Get(status, "cbsData.#(cbsType==\"OIL\").cbsRemainingMileage").Float()
	oilCheckRemainingMileage.Set(oilCheckRemainingMileageValue)

	remainingFuelValue := gjson.Get(status, "remainingFuel").Float()
	remainingFuel.Set(remainingFuelValue)

	remainingRangeElectricValue := gjson.Get(status, "remainingRangeElectric").Float()
	remainingRangeElectric.Set(remainingRangeElectricValue)

	remainingRangeHybridValue := gjson.Get(status, "remainingRangeFuel").Float()
	remainingRangeHybrid.Set(remainingRangeHybridValue)

	trunkValue := gjson.Get(status, "trunk").String()
	if trunkValue == "OPEN" {
		trunk.Set(0)
	} else {
		trunk.Set(1)
	}

	vehicleCheckRemainingMileageValue := gjson.Get(status, "cbsData.#(cbsType==\"VEHICLE_CHECK\").cbsRemainingMileage").Float()
	vehicleCheckRemainingMileage.Set(vehicleCheckRemainingMileageValue)

	vehicleCheckDueDateValue := gjson.Get(status, "cbsData.#(cbsType==\"VEHICLE_CHECK\").cbsDueDate").String()
	vehicleCheckDueDate.Set(convertToEpoch(vehicleCheckDueDateValue))

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
}
