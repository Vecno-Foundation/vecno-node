package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://c788d2afcf3c1d2a8bd23e92eb0e6306712ebb6a34c8026e217c9ba05eb9d580b8dcf39859bd4ddee80434f2793b2574f0baddcc4976d979f5d6e30247476591@203.161.54.144:30303", // Dev
	"", // Explorer

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
