package define

import "sync"

var (
	File    string
	Url     string
	TimeOut int
	OutPut  string
)

var (
	HostList      []string
	ReqIpList     []string
	ReqDomainList []string
	ResultList    []string
)

var (
	Modes     = []string{"fofa", "quake", "hunter", "zoomeye", "vt", "netlas", "pulsedive"} //模块
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4621.0 Safari/537.36"
)

type Configure struct {
	FofaMail     string `mapstructure:"FofaMail" json:"FofaMail" yaml:"FofaMail"`
	FofaKey      string `mapstructure:"FofaKey" json:"FofaKey" yaml:"FofaKey"`
	QuakeKey     string `mapstructure:"QuakeKey" json:"QuakeKey" yaml:"QuakeKey"`
	HunterKey    string `mapstructure:"HunterKey" json:"HunterKey" yaml:"HunterKey"`
	ZoomeyeKey   string `mapstructure:"ZoomeyeKey" json:"ZoomeyeKey" yaml:"ZoomeyeKey"`
	PulsediveKey string `mapstructure:"PulsediveKey" json:"PulsediveKey" yaml:"PulsediveKey"`
	DaydaymapKey string `mapstructure:"DaydaymapKey" json:"DaydaymapKey" yaml:"DaydaymapKey"`
}

var ModeToGrammar = map[string]string{
	"fofa":   "=",
	"hunter": "=",
	"quake":  ":",
}

var Once sync.Once
