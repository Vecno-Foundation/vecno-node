package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://c788d2afcf3c1d2a8bd23e92eb0e6306712ebb6a34c8026e217c9ba05eb9d580b8dcf39859bd4ddee80434f2793b2574f0baddcc4976d979f5d6e30247476591@203.161.54.144:30303", // Dev
	"enode://3ebd502effd02de6e1c4e034e145051cba90c586bc3e11db64b1d985b27bd8a294bf477d92cb067ad2d389c5233c749f348610d99b70d15d26c510e79eb7d212@66.29.155.55:30303",   // Explorer

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
