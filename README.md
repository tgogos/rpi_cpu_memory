# How to get CPU and Memory usage from a Raspberry Pi
This code was tested on a Raspberry Pi 3 with Raspian Jessie, values were "pushed" to an [influxDB](https://www.influxdata.com/time-series-platform/influxdb/) and then were visualized with [Grafana](http://grafana.org/).

# How to run
```bash
# first "go get" the goprocinfo library:
go get github.com/c9s/goprocinfo

# then run it with
go run main.go
```

# Where do these stats come from?
[goprocinfo](https://github.com/c9s/goprocinfo) uses `/proc/stat` and `/proc/meminfo` pseudo-files.

## How is the CPU usage calculated by those values?
I follow the approach proposed here: [Accurate calculation of CPU usage given in percentage in Linux](http://stackoverflow.com/questions/23367857/accurate-calculation-of-cpu-usage-given-in-percentage-in-linux)

## What about the Memory?
I follow the approach proposed here (by the `htop` command author): [How to calculate memory usage from /proc/meminfo (like htop)](http://stackoverflow.com/questions/41224738/how-to-calculate-memory-usage-from-proc-meminfo-like-htop?noredirect=1&lq=1)
