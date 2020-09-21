package bmwcd

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/tidwall/gjson"
)

const (
	namespace = "bmwcd"
)

type Collector struct {
	brakeFluidCheckDueDate       *prometheus.Desc
	chargingStatus               *prometheus.Desc
	chargeLevel                  *prometheus.Desc
	collectDuration              *prometheus.Desc
	collectFailures              prometheus.Counter
	connectionStatus             *prometheus.Desc
	doorLockState                *prometheus.Desc
	doorDriverFront              *prometheus.Desc
	doorDriverRear               *prometheus.Desc
	doorPassengerFront           *prometheus.Desc
	doorPassengerRear            *prometheus.Desc
	hood                         *prometheus.Desc
	mileage                      *prometheus.Desc
	mutex                        sync.Mutex
	oilCheckDueDate              *prometheus.Desc
	oilCheckRemainingMileage     *prometheus.Desc
	password                     string
	remainingFuel                *prometheus.Desc
	remainingRangeElectric       *prometheus.Desc
	remainingRangeHybrid         *prometheus.Desc
	trunk                        *prometheus.Desc
	up                           *prometheus.Desc
	username                     string
	vehicleCheckDueDate          *prometheus.Desc
	vehicleCheckRemainingMileage *prometheus.Desc
	windowDriverFront            *prometheus.Desc
	windowDriverRear             *prometheus.Desc
	windowPassengerFront         *prometheus.Desc
	windowPassengerRear          *prometheus.Desc
}

func NewCollector(username, password string) *Collector {
	return &Collector{
		brakeFluidCheckDueDate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "brake_fluid_check_cbs_due_date"),
			"Brake fluid check cbs due by date",
			nil,
			nil),
		chargingStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "charging_status"),
			"Not charging (0), charging (1), fully charged (2)",
			nil,
			nil),
		chargeLevel: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "charge_level"),
			"Charge percentage",
			nil,
			nil),
		collectDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "collect_duration_seconds"),
			"The time it took to collect the metrics in seconds",
			nil,
			nil),
		collectFailures: prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "collect_failures",
				Help:      "The number of collection failures since the exporter was started",
			},
		),
		connectionStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connection_status"),
			"Charging cable connected (1) or disconnected (0)",
			nil,
			nil),
		doorLockState: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "door_lock_state"),
			"Doors unlocked (0) or closed (1)",
			nil,
			nil),
		doorDriverFront: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "door_driver_front_state"),
			"Door open (0) or closed (1)",
			nil,
			nil),
		doorDriverRear: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "door_driver_rear_state"),
			"Door open (0) or closed (1)",
			nil,
			nil),
		doorPassengerFront: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "door_passenger_front_state"),
			"Door open (0) or closed (1)",
			nil,
			nil),
		doorPassengerRear: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "door_passenger_rear_state"),
			"Door open (0) or closed (1)",
			nil,
			nil),
		hood: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "hood_state"),
			"Hood open (0) or closed (1)",
			nil,
			nil),
		mileage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "mileage"),
			"The current mileage of the car",
			nil,
			nil),
		oilCheckDueDate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "oil_check_cbs_due_date"),
			"Oil cbs due by date",
			nil,
			nil),
		oilCheckRemainingMileage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "oil_cbs_remaining_mileage"),
			"Remaining kilometers before oil cbs",
			nil,
			nil),
		password: password,
		remainingFuel: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "remaining_fuel"),
			"Remaining liters of fuel in the tank",
			nil,
			nil),
		remainingRangeElectric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "remaining_electric_range"),
			"Remaining kilometers of electric range",
			nil,
			nil),
		remainingRangeHybrid: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "remaining_hybrid_range"),
			"Remaining kilometers of hybrid range",
			nil,
			nil),
		trunk: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "trunk_state"),
			"Trunk open (0) or closed (1)",
			nil,
			nil),
		up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"Status of BMW Connected Drive API",
			nil,
			nil),
		username: username,
		vehicleCheckDueDate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "vehicle_check_cbs_due_date"),
			"Vehicle check cbs due by date",
			nil,
			nil),
		vehicleCheckRemainingMileage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "vehicle_check_cbs_remaining_mileage"),
			"Remaining kilometers before vehicle check cbs",
			nil,
			nil),
		windowDriverFront: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "window_driver_front_state"),
			"Window open (0) or closed (1)",
			nil,
			nil),
		windowDriverRear: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "window_driver_rear_state"),
			"Window open (0) or closed (1)",
			nil,
			nil),
		windowPassengerFront: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "window_passenger_front_state"),
			"Window open (0) or closed (1)",
			nil,
			nil),
		windowPassengerRear: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "window_passenger_rear_state"),
			"Window open (0) or closed (1)",
			nil,
			nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.collectFailures.Describe(ch)
	ch <- c.brakeFluidCheckDueDate
	ch <- c.chargingStatus
	ch <- c.chargeLevel
	ch <- c.collectDuration
	ch <- c.connectionStatus
	ch <- c.doorLockState
	ch <- c.doorDriverFront
	ch <- c.doorDriverRear
	ch <- c.doorPassengerFront
	ch <- c.doorPassengerRear
	ch <- c.hood
	ch <- c.mileage
	ch <- c.oilCheckDueDate
	ch <- c.oilCheckRemainingMileage
	ch <- c.remainingFuel
	ch <- c.remainingRangeElectric
	ch <- c.remainingRangeHybrid
	ch <- c.trunk
	ch <- c.up
	ch <- c.vehicleCheckDueDate
	ch <- c.vehicleCheckRemainingMileage
	ch <- c.windowDriverFront
	ch <- c.windowDriverRear
	ch <- c.windowPassengerFront
	ch <- c.windowPassengerRear
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Infoln("Started collection")

	startTime := time.Now()

	defer func() {
		duration := time.Since(startTime).Seconds()
		log.Infof("Collection completed in %f seconds", duration)
		ch <- prometheus.MustNewConstMetric(c.collectDuration, prometheus.GaugeValue, duration)
	}()

	vin, err := getVehicleVin(getOAuthToken(c.username, c.password))

	if err != nil {
		log.Errorln("Failed to get the vehicle VIN")
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		c.collectFailures.Inc()
		return
	}

	token, err := getOAuthToken(c.username, c.password)

	if err != nil {
		log.Errorln("Failed to get an access token")
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		c.collectFailures.Inc()
		return
	}

	status, err := getVehicleStatus(string(token), vin)

	if err != nil {
		log.Errorln("Failed to get vehicle status")
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		c.collectFailures.Inc()
		return
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)

	brakeFluidCheckDueDate := convertToEpoch(gjson.Get(status, "cbsData.#(cbsType==\"BRAKE_FLUID\").cbsDueDate").String())
	ch <- prometheus.MustNewConstMetric(c.brakeFluidCheckDueDate, prometheus.GaugeValue, brakeFluidCheckDueDate)

	chargeLevel := gjson.Get(status, "chargingLevelHv").Float()
	ch <- prometheus.MustNewConstMetric(c.chargeLevel, prometheus.GaugeValue, chargeLevel)

	chargingStatus := gjson.Get(status, "chargingStatus").String()
	if chargingStatus == "FINISHED_FULLY_CHARGED" {
		ch <- prometheus.MustNewConstMetric(c.chargingStatus, prometheus.GaugeValue, 2)
	} else if chargingStatus == "CHARGING" {
		ch <- prometheus.MustNewConstMetric(c.chargingStatus, prometheus.GaugeValue, 1)
	} else {
		ch <- prometheus.MustNewConstMetric(c.chargingStatus, prometheus.GaugeValue, 0)
	}

	connectionStatus := gjson.Get(status, "connectionStatus").String()
	if connectionStatus == "CONNECTED" {
		ch <- prometheus.MustNewConstMetric(c.connectionStatus, prometheus.GaugeValue, 1)
	} else {
		ch <- prometheus.MustNewConstMetric(c.connectionStatus, prometheus.GaugeValue, 0)
	}

	doorDriverFront := gjson.Get(status, "doorDriverFront").String()
	if doorDriverFront == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.doorDriverFront, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.doorDriverFront, prometheus.GaugeValue, 1)
	}

	doorDriverRear := gjson.Get(status, "doorDriverRear").String()
	if doorDriverRear == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.doorDriverRear, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.doorDriverRear, prometheus.GaugeValue, 1)
	}

	doorPassengerFront := gjson.Get(status, "doorPassengerFront").String()
	if doorPassengerFront == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.doorPassengerFront, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.doorPassengerFront, prometheus.GaugeValue, 1)
	}

	doorPassengerRear := gjson.Get(status, "doorPassengerRear").String()
	if doorPassengerRear == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.doorPassengerRear, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.doorPassengerRear, prometheus.GaugeValue, 1)
	}

	doorLockState := gjson.Get(status, "doorLockState").String()
	if doorLockState == "UNLOCKED" {
		ch <- prometheus.MustNewConstMetric(c.doorLockState, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.doorLockState, prometheus.GaugeValue, 1)
	}

	hood := gjson.Get(status, "hood").String()
	if hood == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.hood, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.hood, prometheus.GaugeValue, 1)
	}

	mileage := gjson.Get(status, "mileage").Float()
	ch <- prometheus.MustNewConstMetric(c.mileage, prometheus.GaugeValue, mileage)

	oilCheckDueDate := convertToEpoch(gjson.Get(status, "cbsData.#(cbsType==\"OIL\").cbsDueDate").String())
	ch <- prometheus.MustNewConstMetric(c.oilCheckDueDate, prometheus.GaugeValue, oilCheckDueDate)

	oilCheckRemainingMileage := gjson.Get(status, "cbsData.#(cbsType==\"OIL\").cbsRemainingMileage").Float()
	ch <- prometheus.MustNewConstMetric(c.oilCheckRemainingMileage, prometheus.GaugeValue, oilCheckRemainingMileage)

	remainingFuel := gjson.Get(status, "remainingFuel").Float()
	ch <- prometheus.MustNewConstMetric(c.remainingFuel, prometheus.GaugeValue, remainingFuel)

	remainingRangeElectric := gjson.Get(status, "remainingRangeElectric").Float()
	ch <- prometheus.MustNewConstMetric(c.remainingRangeElectric, prometheus.GaugeValue, remainingRangeElectric)

	remainingRangeHybrid := gjson.Get(status, "remainingRangeFuel").Float()
	ch <- prometheus.MustNewConstMetric(c.remainingRangeHybrid, prometheus.GaugeValue, remainingRangeHybrid)

	trunk := gjson.Get(status, "trunk").String()
	if trunk == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.trunk, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.trunk, prometheus.GaugeValue, 1)
	}

	vehicleCheckRemainingMileage := gjson.Get(status, "cbsData.#(cbsType==\"VEHICLE_CHECK\").cbsRemainingMileage").Float()
	ch <- prometheus.MustNewConstMetric(c.vehicleCheckRemainingMileage, prometheus.GaugeValue, vehicleCheckRemainingMileage)

	vehicleCheckDueDate := convertToEpoch(gjson.Get(status, "cbsData.#(cbsType==\"VEHICLE_CHECK\").cbsDueDate").String())
	ch <- prometheus.MustNewConstMetric(c.vehicleCheckDueDate, prometheus.GaugeValue, vehicleCheckDueDate)

	windowDriverFront := gjson.Get(status, "windowDriverFront").String()
	if windowDriverFront == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.windowDriverFront, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.windowDriverFront, prometheus.GaugeValue, 1)
	}

	windowDriverRear := gjson.Get(status, "windowDriverRear").String()
	if windowDriverRear == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.windowDriverRear, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.windowDriverRear, prometheus.GaugeValue, 1)
	}

	windowPassengerFront := gjson.Get(status, "windowPassengerFront").String()
	if windowPassengerFront == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.windowPassengerFront, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.windowPassengerFront, prometheus.GaugeValue, 1)
	}

	windowPassengerRear := gjson.Get(status, "windowPassengerRear").String()
	if windowPassengerRear == "OPEN" {
		ch <- prometheus.MustNewConstMetric(c.windowPassengerRear, prometheus.GaugeValue, 0)
	} else {
		ch <- prometheus.MustNewConstMetric(c.windowPassengerRear, prometheus.GaugeValue, 1)
	}
}
