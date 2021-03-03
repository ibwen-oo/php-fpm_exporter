package phpfpm

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"sync"
)

const (
	FPMUP = 1
	FPMDOWN = 0
)

type PHPCollector struct {
	mutex sync.Mutex
	url *URL
	// php-fpm 状态
	up prometheus.Gauge
	// 当前池子接受的请求数
	acceptedConn *prometheus.Desc
	// 请求等待队列,如果值不为0,说明php进程数太小或php处理请求太慢
	listenQueue *prometheus.Desc
	// 请求等待队列最高的数量
	maxListenQueue *prometheus.Desc
	// socket 等待队列长度
	listenQueueLen *prometheus.Desc
	// 空闲进程数
	idleProcesses *prometheus.Desc
	// 活跃进程数
	activeProcesses *prometheus.Desc
	// 总进程数
	totalProcesses *prometheus.Desc
	// 最大的活跃进程数(从FPM启动开始基础)
	maxActiveProcesses *prometheus.Desc
	// 达到进程最大数量限制的次数,如果这个值不为0,说明最大进程数配置的太小了
	maxChildrenReached *prometheus.Desc
	// 启用了php-fpm slog.log, php慢请求的数量
	slowRequests *prometheus.Desc
}

// NewPHPCollector 创建 PHPCollector 实例
func NewPHPCollector(namespace string, u *URL) *PHPCollector {
	return &PHPCollector{
		url: u,
		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "up", Help: "php-fpm status(up or down)",
				Namespace: namespace, ConstLabels: map[string]string{"app": "php-fpm"},
			}),
		acceptedConn: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "accepted_conn"),
			"The number of requests currently accepted by the pool.",
			[]string{"app"}, nil),
		listenQueue: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "listen_queue"),
			"The number of requests in the queue of pending connections.",
			[]string{"app"}, nil),
		maxListenQueue: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_listen_queue"),
			"The maximum number of requests in the queue of pending connections since FPM has started.",
			[]string{"app"}, nil),
		listenQueueLen: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "listen_queue_len"),
			"The size of the socket queue of pending connections.",
			[]string{"app"}, nil),
		idleProcesses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "idle_processes"),
			"The number of idle processes.",
			[]string{"app"}, nil),
		activeProcesses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "active_processes"),
			"The number of active processes.",
			[]string{"app"}, nil),
		totalProcesses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "total_processes"),
			"The number of idle + active processes.",
			[]string{"app"}, nil),
		maxActiveProcesses: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_active_processes"),
			"The maximum number of active processes since FPM has started.",
			[]string{"app"}, nil),
		maxChildrenReached: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "max_children_reached"),
			"The number of times, the process limit has been reached, when pm tries to start more children (works only for pm 'dynamic' and 'ondemand').",
			[]string{"app"}, nil),
		slowRequests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "slow_requests"),
			"The number of requests that exceeded your 'request_slowlog_timeout' value.",
			[]string{"app"}, nil),
	}
}


func (p *PHPCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.up.Desc()
	ch <- p.acceptedConn
	ch <- p.listenQueue
	ch <- p.maxListenQueue
	ch <- p.listenQueueLen
	ch <- p.idleProcesses
	ch <- p.activeProcesses
	ch <- p.totalProcesses
	ch <- p.maxActiveProcesses
	ch <- p.maxChildrenReached
	ch <- p.slowRequests
}

func (p *PHPCollector) Collect(ch chan<- prometheus.Metric) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	metrics, err := p.url.QueryStatus()
	if err != nil {
		zap.L().Error("Collect p.url.QueryStatus", zap.Error(err))
		p.up.Set(FPMDOWN)
		ch <- p.up
		return
	}

	p.up.Set(FPMUP)
	ch <-p.up
	ch <- prometheus.MustNewConstMetric(
		p.acceptedConn,
		prometheus.GaugeValue,
		float64(metrics.AcceptedConnections),
		"php-fpm",
		)

	ch <- prometheus.MustNewConstMetric(
		p.activeProcesses,
		prometheus.GaugeValue,
		float64(metrics.ActiveProcesses),
		"php-fpm",
		)

	ch <- prometheus.MustNewConstMetric(
		p.idleProcesses,
		prometheus.GaugeValue,
		float64(metrics.IdleProcesses),
		"php-fpm",
	)

	ch <- prometheus.MustNewConstMetric(
		p.slowRequests,
		prometheus.GaugeValue,
		float64(metrics.SlowRequests),
		"php-fpm",
	)

	ch <- prometheus.MustNewConstMetric(
		p.maxActiveProcesses,
		prometheus.GaugeValue,
		float64(metrics.MaxActiveProcesses),
		"php-fpm",
	)

	ch <- prometheus.MustNewConstMetric(
		p.maxChildrenReached,
		prometheus.GaugeValue,
		float64(metrics.MaxChildrenReached),
		"php-fpm",
	)

	ch <- prometheus.MustNewConstMetric(
		p.listenQueue,
		prometheus.GaugeValue,
		float64(metrics.ListenQueue),
		"php-fpm",
	)

	ch <- prometheus.MustNewConstMetric(
		p.listenQueueLen,
		prometheus.GaugeValue,
		float64(metrics.ListenQueueLength),
		"php-fpm",
	)

	ch <- prometheus.MustNewConstMetric(
		p.totalProcesses,
		prometheus.GaugeValue,
		float64(metrics.TotalProcesses),
		"php-fpm",
	)

	ch <- prometheus.MustNewConstMetric(
		p.maxListenQueue,
		prometheus.GaugeValue,
		float64(metrics.MaxListenQueue),
		"php-fpm",
	)
}
