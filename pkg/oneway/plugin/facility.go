package plugin

const (
	Input   = "input"
	Transit = "transit"
	Process = "process"
	Output  = "output"
	CronJob = "cronjobs"
)

var (
	Service string

	ChanSize int32
	// I2TChan  = make(chan []byte)
	// T2PChan  = make(chan map[string]string)
	// P2OChan  = make(chan map[string]string)

	I2TChan = make(map[string](chan []byte))
	T2PChan = make(map[string](chan map[string]string))
	P2OChan = make(map[string](chan map[string]string))

	Metrics = &Metric{}
	Records = &Record{}

	InputStatus   = &inputStatus{Completed: false}
	TransitStatus = &transitStatus{Completed: false}
	ProcessStatus = &processStatus{Completed: false}
	OutputStatus  = &outputStatus{Completed: false}

	IsOneTimeExec = false
)

type Metric struct {
	InputOK  int32 `json:"inputOK"`
	InputErr int32 `json:"inputErr"`

	TransitOK  int32 `json:"transitOK"`
	TransitErr int32 `json:"transitErr"`

	ProcessOK  int32 `json:"processOK"`
	ProcessErr int32 `json:"processErr"`

	OutputOK  int32 `json:"outputOK"`
	OutputErr int32 `json:"outputErr"`
}

type Record struct {
	Input   [][]byte
	Transit [][]byte
	Process [][]byte
	Output  [][]byte
}

type inputStatus struct {
	Completed bool
}

type transitStatus struct {
	Completed bool
}
type processStatus struct {
	Completed bool
}

type outputStatus struct {
	Completed bool
}
