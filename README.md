## Pi Monitor

This project was build to monitor the temperature of a garage freezer. 
It's being used to read the values from DHT22 sensors, and expose these values as metrics to Prometheus. 
I am using this app with a Pi Zero W, that has two sensors connected to it, to monitor the internal and external temperatures/humidities.


### Building the binary:

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

In order to read the values from the sensors, the library [d2r2/go-dht](https://github.com/d2r2/go-dht). It uses a C implementation, which is important, since I am using a Pi Zero for this monitoring (Single core).

### Grafana 

In order to display the metrics on a nice dashboard, I am using Grafana. You can start it with this command:

`docker run -d --name=grafana -p 3000:3000 grafana/grafana`

### Prometheus

In order to collect the metrics and make them available to Grafana, you should run Prometheus. You can start it with this command:

`docker run --name prometheus -d -p 9090:9090 -v $(pwd)/examples/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus`

The above command doesn't take into consideration the data storage. It should be only used to validate the whole process.