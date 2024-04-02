package main

//Peer represents the response containing peer information.
type Peer struct {
	Address      string `json:"address"`
	DownloadPort string `json:"download_port"`
	DelayPort    string `json:"delay_port"`
}
