## Pi Monitor

This project was build to monitor the temperature of a garage freezer. 
It's being used to read the values from DHT22 sensors, and expose these values as metrics to Prometheus. 
I am using this app with a Pi Zero W, that has two sensors connected to it, to monitor the internal and external temperatures/humidities.


### Binaty build:

In order to build the binary to run on the Pi, you need to execute the command `make build`.
It will create a Docker builder image, and the build process will use that to run the compilation process inside this container.
The result will be a file called `pi-monitor`, that should be copied to the Pi.

### Configuration files examples:

The folder `./examples` contains examples for the configurations that were used in this project.

#### pi-monitor configuration

The file `./examples/pi-monitor.yml` should be added to the `/etc` folder, and will configure the sensors and its gauge names, as well as the interval to check them, and the port to expose the `/metrics` API for Prometheus.

#### Prometeus

The file `./example/prometheus.yml` contains the example to collect the metrics from your Pi. In my case, I have one that collects the internal and external temperatures and humidities from my garage freezer.

### Sensor Library:

https://github.com/d2r2/go-dht

### Grafana 

docker run -d --name=grafana -p 3000:3000 grafana/grafana

### Prometheus

docker run --name prometheus -d -p 9090:9090 -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus