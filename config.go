package ipcarta

type Config struct {
	ElasticSearchHost string `json:"host"`

	Port string `json:"post"`
}

var config Config
