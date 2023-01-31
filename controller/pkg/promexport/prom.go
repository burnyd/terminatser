package promexport

import (
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var gNMI_metrics = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "gNMI_Metrics",
		Help: "Metrics for gNMI",
	},
	[]string{"Path", "Timestamp", "Target", "Value"},
)

//func RegisterProm() {
//	err := prometheus.Register(gNMI_metrics)
//	handleErr(err)
//}

func ExportToProm(path, ts, target, value string) {
	gNMI_metrics.With(prometheus.Labels{
		"Path":      path,
		"Timestamp": ts,
		"Target":    target,
		"Value":     value,
	})
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Init() {
	prometheus.Register(gNMI_metrics)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(9090), nil))
	}()
}
