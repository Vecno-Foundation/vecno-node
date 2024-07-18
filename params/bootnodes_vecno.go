package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://7cc15bfbd2050042229655007935427312deaf186e460dea25c627b5880c79f35b362a96fd70df5bcd9001cc50354c251f86020e2092931d7af331adfb8fb3b5@203.161.54.144:30303", // Dev
	"enode://0cc5d3fb8623e88fbaa1f63b1e5a7512b4bb135a45dc7423f50e76b73a50b8411c0cbaf61bd618d40a913ef0df327470b631943f8f3dbe5a901d1be4ab679e07@66.29.155.55:30303",   // Explorer

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
