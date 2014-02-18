package bellows

type Config struct {
	Env    string `json:"env"`
	Server struct {
		SSLCertPath string `json:"ssl_cert_path"`
		SSLKeyPath  string `json:"ssl_key_path"`
	} `json:"server"`
	Storage struct {
		Driver string `json:"driver"`
		Dsn    string `json:"dsn"`
	} `json:"storage"`
	Channels struct {
		MutationQueueDepth int `json:"mutation_queue_depth"`
		EventQueueDepth    int `json:"event_queue_depth"`
	} `json:"channels"`
}
