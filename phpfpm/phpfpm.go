package phpfpm

import (
	"encoding/json"
	fastCgiClient "github.com/tomasen/fcgi_client"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

type URL struct {
	Scheme  string
	Address string
	Path    string
	Timeout time.Duration
}

type FPMMetrics struct {
	PoolName            string `json:"pool"`
	ProcessManager      string `json:"process manager"`
	StartTime           int64  `json:"start time"`
	StartSince          int64  `json:"start since"`
	AcceptedConnections int64  `json:"accepted conn"`
	ListenQueue         int64  `json:"listen queue"`
	MaxListenQueue      int64  `json:"max listen queue"`
	ListenQueueLength   int64  `json:"listen queue len"`
	IdleProcesses       int64  `json:"idle processes"`
	ActiveProcesses     int64  `json:"active processes"`
	TotalProcesses      int64  `json:"total processes"`
	MaxActiveProcesses  int64  `json:"max active processes"`
	MaxChildrenReached  int64  `json:"max children reached"`
	SlowRequests        int64  `json:"slow requests"`
}

var ENV = make(map[string]string)

func (u *URL) GenClient() {
	ENV["SCRIPT_FILENAME"] = u.Path
	ENV["SCRIPT_NAME"] = u.Path
	ENV["SERVER_SOFTWARE"] = "go / php-fpm_exporter"
	ENV["REMOTE_ADDR"] = "127.0.0.1"
	ENV["QUERY_STRING"] = "json"
}

func (u *URL) QueryStatus() (metrics *FPMMetrics, err error) {
	fcgi, err := fastCgiClient.DialTimeout(u.Scheme, u.Address, u.Timeout)
	if err != nil {
		zap.L().Error("QueryStatus fastCgiClient.DialTimeout", zap.Error(err))
		return nil, err
	}

	resp, err := fcgi.Get(ENV)
	if err != nil {
		zap.L().Error("QueryStatus fcgi.Get", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("QueryStatus ioutil.ReadAll", zap.Error(err))
		return nil, err
	}
	metrics = new(FPMMetrics)
	if err = json.Unmarshal(content, metrics); err != nil {
		zap.L().Error("QueryStatus json.Unmarshal", zap.Error(err))
		return nil, err
	}

	return metrics, err
}
