package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://37c42a2fc8a7d84e379b2cc9d7df19aa9d9c59885c24b30a07811ba02f1b0110d72640a97ecc45ffa9213b26618f647eef53c8748a9026a724e637d87714bd0c@172.29.144.1:30303", // Dev
	"", // Explorer

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
