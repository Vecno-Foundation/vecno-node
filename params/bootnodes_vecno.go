package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://1d3dc091ef0570b1cb8512dbbd0f02d92907597a04bc931ed86c336b699b7b13b12035ed240c5af215e2e910454ceb03185cd12875ea76b02b60e3afd7f5874b@81.167.190.96:30303", // Explorer
	"", // Dev

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
