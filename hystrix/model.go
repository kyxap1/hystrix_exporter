package hystrix

import (
	"encoding/json"
	"strings"
)

// Unmarshal the given string into a struct
func Unmarshal(line string) (data Data, err error) {
	line = strings.TrimPrefix(line, "data:")
	err = json.Unmarshal([]byte(line), &data)
	return
}

// Latencies of a circuit
type Latencies struct {
	L0   float64 `json:"0,omitempty"`
	L25  float64 `json:"25,omitempty"`
	L50  float64 `json:"50,omitempty"`
	L75  float64 `json:"75,omitempty"`
	L90  float64 `json:"90,omitempty"`
	L95  float64 `json:"95,omitempty"`
	L99  float64 `json:"99,omitempty"`
	L995 float64 `json:"99.5,omitempty"`
	L100 float64 `json:"100,omitempty"`
}

// Data Hystrix main data type
type Data struct {
	Type string `json:"Type,omitempty"`

	Group string `json:"group,omitempty"`
	Name  string `json:"name,omitempty"`

	Open           bool    `json:"isCircuitBreakerOpen,omitempty"`
	ReportingHosts float64 `json:"reportingHosts,omitempty"`

	LatencyExecuteMean float64   `json:"latencyExecute_mean,omitempty"`
	LatencyTotalMean   float64   `json:"latencyTotal_mean,omitempty"`
	LatencyTotal       Latencies `json:"latencyTotal,omitempty"`
	LatencyExecute     Latencies `json:"latencyExecute,omitempty"`

	RequestCount float64 `json:"requestCount,omitempty"`
	ErrorCount   float64 `json:"errorCount,omitempty"`

	CurrentCorePoolSize       float64 `json:"currentCorePoolSize,omitempty"`
	CurrentLargestPoolSize    float64 `json:"currentLargestPoolSize,omitempty"`
	CurrentActiveCount        float64 `json:"currentActiveCount,omitempty"`
	CurrentMaximumPoolSize    float64 `json:"currentMaximumPoolSize,omitempty"`
	CurrentQueueSize          float64 `json:"currentQueueSize,omitempty"`
	CurrentPoolSize           float64 `json:"currentPoolSize,omitempty"`
	CurrentTaskCount          float64 `json:"currentTaskCount,omitempty"`
	CurrentCompletedTaskCount float64 `json:"currentCompletedTaskCount,omitempty"`

	RollingCountCommandRejections float64 `json:"rollingCountCommandRejections,omitempty"`
	RollingCountThreadsExecuted   float64 `json:"rollingCountThreadsExecuted,omitempty"`
	RollingMaxActiveThreads       float64 `json:"rollingMaxActiveThreads,omitempty"`

	RollingCountCollapsedRequests  float64 `json:"rollingCountCollapsedRequests,omitempty"`
	RollingCountShortCircuited     float64 `json:"rollingCountShortCircuited,omitempty"`
	RollingCountThreadPoolRejected float64 `json:"rollingCountThreadPoolRejected,omitempty"`
	RollingCountFallbackEmit       float64 `json:"rollingCountFallbackEmit,omitempty"`
	RollingCountSuccess            float64 `json:"rollingCountSuccess,omitempty"`
	RollingCountTimeout            float64 `json:"rollingCountTimeout,omitempty"`
	RollingCountFailure            float64 `json:"rollingCountFailure,omitempty"`
	RollingCountExceptionsThrown   float64 `json:"rollingCountExceptionsThrown,omitempty"`
}
