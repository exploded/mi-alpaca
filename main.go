package main

const (
	apiPort              = 8080
	DiscoveryPort        = 32227
	DefaultAlpacaApiPort = 11111
	ListenIP             = "127.0.0.1"
	Location             = "Earth"
)

func main() {
	// Load initial switch values and query device states
	MiSetInit()

	// Start discovery server for ASCOM Alpaca device discovery
	discovery := NewDiscoveryServer(DiscoveryPort, apiPort)
	go discovery.Start()
	defer discovery.Close()

	// Start API server for ASCOM Alpaca device control
	api := NewApiServer(apiPort)
	go api.Start()

	// Block forever
	select {}
}
