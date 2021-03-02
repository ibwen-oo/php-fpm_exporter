package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"php-fpm_exporter/logger"
	"php-fpm_exporter/phpfpm"
	"time"
)

func main() {
	// php-fpm 相关参数
	scheme := flag.String("scheme", "tcp", "协议, unix or tcp")
	address := flag.String("address", "127.0.0.1:9000", "请求fpm status地址.\n使用端口: 127.0.0.1:9000;\n使用socket文件: /tmp/php-fcgi.sock")
	path := flag.String("path", "/fpm_status", "请求fpm status路径")
	timeout := flag.Duration("timeout", 3 * time.Second, "请求超时时间")

	// prometheus 相关参数
	namespace := flag.String("namespace", "null", "exporter namespace")
	listenAddress := flag.String("web.address", ":9005", "暴露metrics的端口")
	metricsPath := flag.String("web.path", "/metrics", "暴露metrics的访问路径")

	// 日志相关参数
	logLevel := flag.String("log.level", "Error", "日志级别 [Debug Info Error Warn]")
	logPath := flag.String("log.path", "./", "日志路径,默认为当前路径")
	flag.Parse()
	// 0.解析命令行参数
	url := phpfpm.URL{
		Scheme:  *scheme,
		Address: *address,
		Path:    *path,
		Timeout: *timeout,
	}
	// 1.初始化日志
	if err := logger.InitLogger(*logLevel, *logPath); err != nil {
		log.Fatalln("Init Logger error:", err)
	}

	// 2.构造请求信息
	url.GenClient()

	// prometheus client
	registry := prometheus.NewRegistry()
	if *namespace == "null" {
		*namespace = phpfpm.DefaultNameSpace
	}

	// 构造PHPCollector实例
	fpmCollector := phpfpm.NewPHPCollector(*namespace, &url)
	// 注册自定义的Collector
	registry.MustRegister(fpmCollector)

	// 启动http服务,对外暴露metrics
	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}

	handler := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})

	http.HandleFunc(*metricsPath, func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	log.Fatalln(http.ListenAndServe(*listenAddress, nil))
}