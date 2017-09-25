package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ContaAzul/hystrix-exporter/config"
	"github.com/ContaAzul/hystrix-exporter/hystrix"
)

func init() {
	log.SetHandler(cli.Default)
}

var (
	version    = "dev"
	app        = kingpin.New("hystrix-exporter", "exports hystrix metrics in the prometheus format")
	configFile = app.Flag("config", "config file").Short('c').Default("config.yml").ExistingFile()

	// command gauges
	latencyTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_latency_total",
		Help: "latencies total",
	}, []string{"cluster", "name", "statistic"})
	latencyExecute = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_latency_Execute",
		Help: "latencies execute",
	}, []string{"cluster", "name", "statistic"})
	rollingCountCollapsedRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountCollapsedRequests",
		Help: "rollingCountCollapsedRequests",
	}, []string{"cluster", "name"})
	rollingCountShortCircuited = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountShortCircuited",
		Help: "rollingCountShortCircuited",
	}, []string{"cluster", "name"})
	rollingCountThreadPoolRejected = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountThreadPoolRejected",
		Help: "rollingCountThreadPoolRejected",
	}, []string{"cluster", "name"})
	rollingCountFallbackEmit = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountFallbackEmit",
		Help: "rollingCountFallbackEmit",
	}, []string{"cluster", "name"})
	rollingCountSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountSuccess",
		Help: "rollingCountSuccess",
	}, []string{"cluster", "name"})
	rollingCountTimeout = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountTimeout",
		Help: "rollingCountTimeout",
	}, []string{"cluster", "name"})
	rollingCountFailure = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountFailure",
		Help: "rollingCountFailure",
	}, []string{"cluster", "name"})
	rollingCountExceptionsThrown = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_rollingCountExceptionsThrown",
		Help: "rollingCountExceptionsThrown",
	}, []string{"cluster", "name"})
	circuitOpen = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_command_circuit_open",
		Help: "circuit open, 1 means true",
	}, []string{"cluster", "name"})

	// thread pool gauges:
	threadPoolCurrentCorePoolSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentCorePoolSize",
		Help: "currentCorePoolSize",
	}, []string{"cluster", "name"})
	threadPoolCurrentLargestPoolSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentLargestPoolSize",
		Help: "currentLargestPoolSize",
	}, []string{"cluster", "name"})
	threadPoolCurrentActiveCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentActiveCount",
		Help: "currentActiveCount",
	}, []string{"cluster", "name"})
	threadPoolCurrentMaximumPoolSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentMaximumPoolSize",
		Help: "currentMaximumPoolSize",
	}, []string{"cluster", "name"})
	threadPoolCurrentQueueSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentQueueSize",
		Help: "currentQueueSize",
	}, []string{"cluster", "name"})
	threadPoolCurrentTaskCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentTaskCount",
		Help: "currentTaskCount",
	}, []string{"cluster", "name"})
	threadPoolCurrentCompletedTaskCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentCompletedTaskCount",
		Help: "currentCompletedTaskCount",
	}, []string{"cluster", "name"})
	threadPoolRollingMaxActiveThreads = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_rollingMaxActiveThreads",
		Help: "rollingMaxActiveThreads",
	}, []string{"cluster", "name"})
	threadPoolRollingCountCommandRejections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_rollingCountCommandRejections",
		Help: "rollingCountCommandRejections",
	}, []string{"cluster", "name"})
	threadPoolReportingHosts = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_reportingHosts",
		Help: "reportingHosts",
	}, []string{"cluster", "name"})
	threadPoolCurrentPoolSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_currentPoolSize",
		Help: "currentPoolSize",
	}, []string{"cluster", "name"})
	threadPoolRollingCountThreadsExecuted = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hystrix_threadpool_rollingCountThreadsExecuted",
		Help: "rollingCountThreadsExecuted",
	}, []string{"cluster", "name"})
)

func main() {
	app.Version(version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	prometheus.DefaultRegisterer.MustRegister(
		latencyTotal,
		latencyExecute,
		rollingCountCollapsedRequests,
		rollingCountShortCircuited,
		rollingCountThreadPoolRejected,
		rollingCountFallbackEmit,
		rollingCountSuccess,
		rollingCountTimeout,
		rollingCountFailure,
		rollingCountExceptionsThrown,
		circuitOpen,

		threadPoolCurrentCorePoolSize,
		threadPoolCurrentLargestPoolSize,
		threadPoolCurrentActiveCount,
		threadPoolCurrentMaximumPoolSize,
		threadPoolCurrentQueueSize,
		threadPoolCurrentTaskCount,
		threadPoolCurrentCompletedTaskCount,
		threadPoolRollingMaxActiveThreads,
		threadPoolRollingCountCommandRejections,
		threadPoolReportingHosts,
		threadPoolCurrentPoolSize,
		threadPoolRollingCountThreadsExecuted,
	)

	http.Handle("/metrics", promhttp.Handler())
	conf, err := config.Parse(*configFile)
	if err != nil {
		log.WithError(err).Fatalf("failed to parse config file: %s", *configFile)
	}

	for _, cluster := range conf.Clusters {
		go read(cluster.URL, cluster.Name)
	}

	http.ListenAndServe(":9333", http.DefaultServeMux)
}

func read(url, cluster string) {
	var log = log.WithField("url", url).WithField("cluster", cluster)
	log.Info("starting")
	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).Warn("failed to read url")
		time.Sleep(time.Second * 10)
		read(url, cluster)
		return
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !isData(line) {
			continue
		}
		doReport(cluster, line)
	}
	// stream stop reporting
	time.Sleep(time.Second * 1)
	read(url, cluster)
}

func doReport(cluster, line string) {
	data, err := hystrix.Unmarshal(line)
	if err != nil {
		return
	}
	switch data.Type {
	case "HystrixThreadPool":
		reportThreadPool(cluster, data)
	case "HystrixCommand":
		reportCommand(cluster, data)
	case "meta":
		return
	default:
		log.Warnf("don't know what to do with type '%s'", data.Type)
	}
}

func isData(line string) bool {
	return strings.HasPrefix(line, "data:")
}

func reportCommand(cluster string, data hystrix.Data) {
	latencyTotal.WithLabelValues(cluster, data.Name, "0").Set(data.LatencyTotal.L0)
	latencyTotal.WithLabelValues(cluster, data.Name, "25").Set(data.LatencyTotal.L25)
	latencyTotal.WithLabelValues(cluster, data.Name, "50").Set(data.LatencyTotal.L50)
	latencyTotal.WithLabelValues(cluster, data.Name, "75").Set(data.LatencyTotal.L75)
	latencyTotal.WithLabelValues(cluster, data.Name, "90").Set(data.LatencyTotal.L90)
	latencyTotal.WithLabelValues(cluster, data.Name, "95").Set(data.LatencyTotal.L95)
	latencyTotal.WithLabelValues(cluster, data.Name, "99").Set(data.LatencyTotal.L99)
	latencyTotal.WithLabelValues(cluster, data.Name, "99.5").Set(data.LatencyTotal.L995)
	latencyTotal.WithLabelValues(cluster, data.Name, "100").Set(data.LatencyTotal.L100)
	latencyTotal.WithLabelValues(cluster, data.Name, "mean").Set(data.LatencyTotalMean)
	latencyExecute.WithLabelValues(cluster, data.Name, "0").Set(data.LatencyExecute.L0)
	latencyExecute.WithLabelValues(cluster, data.Name, "25").Set(data.LatencyExecute.L25)
	latencyExecute.WithLabelValues(cluster, data.Name, "50").Set(data.LatencyExecute.L50)
	latencyExecute.WithLabelValues(cluster, data.Name, "75").Set(data.LatencyExecute.L75)
	latencyExecute.WithLabelValues(cluster, data.Name, "90").Set(data.LatencyExecute.L90)
	latencyExecute.WithLabelValues(cluster, data.Name, "95").Set(data.LatencyExecute.L95)
	latencyExecute.WithLabelValues(cluster, data.Name, "99").Set(data.LatencyExecute.L99)
	latencyExecute.WithLabelValues(cluster, data.Name, "99.5").Set(data.LatencyExecute.L995)
	latencyExecute.WithLabelValues(cluster, data.Name, "100").Set(data.LatencyExecute.L100)
	latencyExecute.WithLabelValues(cluster, data.Name, "mean").Set(data.LatencyExecuteMean)

	rollingCountCollapsedRequests.WithLabelValues(cluster, data.Name).Set(data.RollingCountCollapsedRequests)
	rollingCountShortCircuited.WithLabelValues(cluster, data.Name).Set(data.RollingCountShortCircuited)
	rollingCountThreadPoolRejected.WithLabelValues(cluster, data.Name).Set(data.RollingCountThreadPoolRejected)
	rollingCountFallbackEmit.WithLabelValues(cluster, data.Name).Set(data.RollingCountFallbackEmit)
	rollingCountSuccess.WithLabelValues(cluster, data.Name).Set(data.RollingCountSuccess)
	rollingCountTimeout.WithLabelValues(cluster, data.Name).Set(data.RollingCountTimeout)
	rollingCountFailure.WithLabelValues(cluster, data.Name).Set(data.RollingCountFailure)
	rollingCountExceptionsThrown.WithLabelValues(cluster, data.Name).Set(data.RollingCountExceptionsThrown)

	circuitOpen.WithLabelValues(cluster, data.Name).Set(boolToFloat64(data.Open))
}

func reportThreadPool(cluster string, data hystrix.Data) {
	threadPoolCurrentCorePoolSize.WithLabelValues(cluster, data.Name).Set(data.CurrentCorePoolSize)
	threadPoolCurrentLargestPoolSize.WithLabelValues(cluster, data.Name).Set(data.CurrentLargestPoolSize)
	threadPoolCurrentActiveCount.WithLabelValues(cluster, data.Name).Set(data.CurrentActiveCount)
	threadPoolCurrentMaximumPoolSize.WithLabelValues(cluster, data.Name).Set(data.CurrentMaximumPoolSize)
	threadPoolCurrentQueueSize.WithLabelValues(cluster, data.Name).Set(data.CurrentQueueSize)
	threadPoolCurrentTaskCount.WithLabelValues(cluster, data.Name).Set(data.CurrentTaskCount)
	threadPoolCurrentCompletedTaskCount.WithLabelValues(cluster, data.Name).Set(data.CurrentCompletedTaskCount)
	threadPoolRollingMaxActiveThreads.WithLabelValues(cluster, data.Name).Set(data.RollingMaxActiveThreads)
	threadPoolRollingCountCommandRejections.WithLabelValues(cluster, data.Name).Set(data.RollingCountCommandRejections)
	threadPoolReportingHosts.WithLabelValues(cluster, data.Name).Set(data.ReportingHosts)
	threadPoolCurrentPoolSize.WithLabelValues(cluster, data.Name).Set(data.CurrentPoolSize)
	threadPoolRollingCountThreadsExecuted.WithLabelValues(cluster, data.Name).Set(data.RollingCountThreadsExecuted)
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
