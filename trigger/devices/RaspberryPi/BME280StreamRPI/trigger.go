package BME280StreamRPI

import (
	"context"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
	"strconv"
	"time"
)

// log is the default package logger
var log = logger.GetLogger("trigger-BME280Stream-RPI")

var interval = 500

// BME280Factory My Trigger factory
type BME280Factory struct {
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &BME280Factory{metadata: md}
}

//New Creates a new trigger instance for a given id
func (t *BME280Factory) New(config *trigger.Config) trigger.Trigger {
	return &BME280Trigger{metadata: t.metadata, config: config}
}

// BME280Trigger is a stub for your Trigger implementation
type BME280Trigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

// Init implements trigger.Trigger.Init
func (t *BME280Trigger) Init(runner action.Runner) {
	t.runner = runner

	if t.config.Settings == nil {
		log.Infof("No configuration set for the trigger... Using default configuration...")
	}

	if _, ok := t.config.Settings["delay_ms"]; !ok {
		log.Infof("No delay has been set. Using default value (", interval, "ms)")
	} else {
		interval, _ = strconv.Atoi(t.config.GetSetting("delay_ms"))
	}

	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)
}

// Metadata implements trigger.Trigger.Metadata
func (t *BME280Trigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *BME280Trigger) Start() error {
	// start the trigger
	log.Debug("Start Trigger BME280Stream for Raspberry PI")
	handlers := t.config.Handlers
	//t.timers = make(map[string]*scheduler.Job)

	log.Debug("Processing handlers")
	for _, handler := range handlers {
		t.scheduleRepeating(handler)
		log.Debugf("Processing Handler: %s", handler.ActionId)
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *BME280Trigger) Stop() error {
	// stop the trigger
	return nil
}

func (t *BME280Trigger) scheduleRepeating(endpoint *trigger.HandlerConfig) {

	log.Debug("Repeating every ", interval, "ms")

	fn2 := func() {
		act := action.Get(endpoint.ActionId)

		temp, press, hum, err := t.getDataFromSensor(endpoint)
		if err != nil {
			log.Error("Error while reading sensor data. Err: ", err.Error())
		}

		data := make(map[string]interface{})
		data["Temperature"] = temp
		data["Pressure"] = press
		data["Humidity"] = hum

		log.Debug("Temperature: ", temp, " C, Pressure: ", press, " hPa, Humidity: ", hum, " %%")
		startAttrs, err := t.metadata.OutputsToAttrs(data, true)

		if err != nil {
			log.Errorf("After run error' %s'\n", err)
		}

		ctx := trigger.NewContext(context.Background(), startAttrs)
		results, err := t.runner.RunAction(ctx, act, nil)

		if err != nil {
			log.Errorf("An error occured while starting the flow. Err:", err)
		}
		log.Info("Exec: ", results)
	}

	// schedule repeating
	doEvery(time.Duration(interval)*time.Millisecond, fn2)
}

func (t *BME280Trigger) getDataFromSensor(endpoint *trigger.HandlerConfig) (temp, press, hum float64, err error) {

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, bme280.I2CAddr)
	if err != nil {
		panic(err)
	}

	b := bme280.New(d)
	err = b.Init()

	return b.EnvData()
}