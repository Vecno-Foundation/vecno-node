package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://ceefac8d745e92d2b617fa3a39dc72b510efb72a7bfc583d1d2caf248ef60769e406db028a960259a24f0b752fadbcaa948de565cf1c3a0fa85783d94e44ecf6@66.29.155.55:30303",   // Dev
	"enode://9ac03ad3a7450928430b210fd18b07424e3cbe6b15466d50f4718173c9bfec1de57b0a839e2a0e6f32b2a98269b65f5805f2cf4365c651a6c0663865ac7c9726@203.161.54.144:30303", // Explorer

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
