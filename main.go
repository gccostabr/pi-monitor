package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d2r2/go-dht"
	"gopkg.in/yaml.v2"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DEFAULT_PORT    = 8080
	DEFAULT_INTEVAL = 5
)

type Config struct {
	Port                int             `yaml:"port"`
	SensorCheckInterval int             `yaml:"interval"`
	SensorMetrics       []SensorMetrics `yaml:"sensor_metrics"`
}

type SensorMetrics struct {
	Pin             int    `yaml:"pin"`
	TemperatureName string `yaml:"temperature_name"`
	TemperatureHelp string `yaml:"temperature_help"`
	HumidityName    string `yaml:"humidity_name"`
	HumidityHelp    string `yaml:"humidity_help"`
}

var config Config

type gauge struct {
	Pin              int
	TemperatureGauge prometheus.Gauge
	HumidityGauge    prometheus.Gauge
}

var gauges []gauge

func updateValues() {

	ticker := time.NewTicker(getInterval())

	for {

		select {
		case <-ticker.C:

			for _, gauge := range gauges {
				temperature, humidity, _, err := dht.ReadDHTxxWithRetry(dht.DHT22, gauge.Pin, false, 10)
				if err != nil {
					log.Printf("Unable to read sensor on pin %d: %+v", gauge.Pin, err)
				} else {
					gauge.HumidityGauge.Set(float64(humidity))
					gauge.TemperatureGauge.Set(float64(temperature))
				}
			}
		}

	}
}

func getInterval() time.Duration {

	checkInterval := config.SensorCheckInterval
	if checkInterval <= 0 {
		log.Printf("Invalid sensor check interval [%d], will use the default value [%d].", checkInterval, DEFAULT_INTEVAL)
		checkInterval = DEFAULT_INTEVAL
	}

	duration, err := time.ParseDuration(fmt.Sprintf("%ds", checkInterval))
	if err != nil {
		log.Fatalf("%v [Please check your interval value]", err)
	}

	log.Printf("Will check the sensors every %s...", duration.String())

	return duration
}

func loadConfig() {
	log.Print("Loading configuration...")
	filename := "/etc/pi-monitor.yml"
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("Unable to load config file [%s]: %+v", filename, err)
		config = Config{Port: DEFAULT_PORT}
		return
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error processing config file [%s]: %+v", filename, err)
	}

	log.Print("Config loaded.")
}

func initGauges() {

	for _, sensormetric := range config.SensorMetrics {
		gauges = append(gauges, gauge{
			Pin: sensormetric.Pin,
			TemperatureGauge: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: sensormetric.TemperatureName,
				Help: sensormetric.TemperatureHelp,
			}),
			HumidityGauge: prometheus.NewGauge(prometheus.GaugeOpts{
				Name: sensormetric.HumidityName,
				Help: sensormetric.HumidityHelp,
			}),
		})
	}

	log.Printf("%d gauges created", len(gauges))
}

func init() {
	loadConfig()
	initGauges()
}

func main() {

	go updateValues()

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}

	http.Handle("/metrics", promhttp.Handler())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting prometheus metrics listener on port %d...", config.Port)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Unable to listen on port %d: %+v", config.Port, err)
		}
	}()

	<-done
	log.Print("Stopping listener...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to stop listener: %+v", err)
	}
	log.Print("Listener properly stopped, good bye!")

}
