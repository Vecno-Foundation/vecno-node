package params

// VecnoBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the Vecno network.
var VecnoBootnodes = []string{
	// Foundation Bootnodes
	"enode://1d3dc091ef0570b1cb8512dbbd0f02d92907597a04bc931ed86c336b699b7b13b12035ed240c5af215e2e910454ceb03185cd12875ea76b02b60e3afd7f5874b@81.167.190.96:30303", // Explorer
	"enode://8a6059cd88ee9ae479ea457049eaf34bce3a76e715daf1c98137c6919f690222a658395a377923dec2d1c24d27389e427df9dd9e16132a21cdbe6eab496bec1f@66.29.155.55:30303",  // Dev

	// Communtiy Bootnodes
	"enode://31d92d58a107fe2abb03eebe0dda2b934a0d736b3a74cf7e94bf4d62898b0f296e0354d2a797241d48b75f609d4a1161b13fc6a68eca4019a12e4bc10506bc5e@203.161.54.144:30303",
	"",
	"",
}

// Once Vecno network has DNS discovery set up,
// this value can be configured.
// var VecnoDNSNetwork = "enrtree://@example"
