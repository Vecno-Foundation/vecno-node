package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://4bab00eabd6360c03960c695fa1eb629a8ee7feea86cd9e0d05d8e1dec81004ccd5168caad9c34ab6f3060f0d227152b6416b75d1e1d3ccff816debc02bc0a15@66.29.155.55:30303", // Dev
	"", // Explorer

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
