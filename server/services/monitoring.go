package services

import "github.com/prometheus/client_golang/prometheus"

var connectedClients = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "chatapp_connected_clients",
	Help: "Current number of connected clients.",
})

func InitMonitoring() {
	// register the connected clients gauge
	prometheus.MustRegister(connectedClients)
}
