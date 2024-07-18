package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://c55080f3a03a23cf6e09bb3c351dc910af5859a124f5880603fcb1fef7becdee5249abc61853d031498121e6a42e53215b242f32fffeadd6a96ef156f219c85c@203.161.54.144:30303", // Dev
	"enode://f22eac4fe23c95923a90e6a7c32a5040a5759f2dbfda89d87d718f8f125283d53f3719a4a393b3829bb4dc569f5c78c93653095410fd16278a00c9b69c7473ef@66.29.155.55:30303",   // Explorer

	// Communtiy Bootnodes
	"",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
