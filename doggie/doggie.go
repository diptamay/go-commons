package doggie

import (
	"fmt"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

const noSampling = 1
const defaultSampleRate = 0.1

var (
	DataDogConfig      *dataDogConfig
	DataDogDefaultTags *[]string
	initOnce           sync.Once
)

func Init() {
	configChan := make(chan *dataDogConfig)
	tagsChan := make(chan *[]string)
	var channelGroup sync.WaitGroup
	channelGroup.Add(2)
	go func() {
		configChan <- getDataDogConfig(&channelGroup)
	}()
	go func() {
		tagsChan <- getDefaultTags(&channelGroup)
	}()
	channelGroup.Wait()
	DataDogConfig = <-configChan
	DataDogDefaultTags = <-tagsChan
}

func formatTags(tagsMap map[string]string) []string {
	result := []string{}
	for key, val := range tagsMap {
		result = append(result, fmt.Sprintf(key+":%s", val))
	}
	return result
}

type statsdClient interface {
	Incr(string, []string, float64) error
	Histogram(string, float64, []string, float64) error
	Gauge(string, float64, []string, float64) error
	Timing(string, time.Duration, []string, float64) error
}

type DataDogClient interface {
	Increment(metric string, tags map[string]string) error
	IncrementNoSampling(metric string, tags map[string]string) error
	Histogram(metric string, value float64, tags map[string]string) error
	HistogramNoSampling(metric string, value float64, tags map[string]string) error
	Gauge(metric string, value float64, tags map[string]string) error
	Timing(metric string, value time.Duration, tags map[string]string) error
	TimingNoSampling(metric string, value time.Duration, tags map[string]string) error
	Namespace() string
}

type dataDogClientImpl struct {
	client    statsdClient
	namespace string
}

func (d *dataDogClientImpl) Increment(metric string, tags map[string]string) error {
	formattedTags := formatTags(tags)
	return d.client.Incr(metric, formattedTags, defaultSampleRate)
}

func (d *dataDogClientImpl) IncrementNoSampling(metric string, tags map[string]string) error {
	formattedTags := formatTags(tags)
	return d.client.Incr(metric, formattedTags, noSampling)
}

func (d *dataDogClientImpl) Histogram(metric string, value float64, tags map[string]string) error {
	formattedTags := formatTags(tags)
	return d.client.Histogram(metric, value, formattedTags, defaultSampleRate)
}

func (d *dataDogClientImpl) HistogramNoSampling(metric string, value float64, tags map[string]string) error {
	formattedTags := formatTags(tags)
	return d.client.Histogram(metric, value, formattedTags, noSampling)
}

func (d *dataDogClientImpl) Gauge(metric string, value float64, tags map[string]string) error {
	formattedTags := formatTags(tags)
	return d.client.Gauge(metric, value, formattedTags, defaultSampleRate)
}

func (d *dataDogClientImpl) Timing(metric string, value time.Duration, tags map[string]string) error {
	formattedTags := formatTags(tags)
	return d.client.Timing(metric, value, formattedTags, defaultSampleRate)
}

func (d *dataDogClientImpl) TimingNoSampling(metric string, value time.Duration, tags map[string]string) error {
	formattedTags := formatTags(tags)
	return d.client.Timing(metric, value, formattedTags, noSampling)
}

func (d *dataDogClientImpl) Namespace() string {
	return d.namespace
}

type dependencies struct {
	New func(string, ...statsd.Option) (*statsd.Client, error)
}

var packageDependencies *dependencies

func MakeClient(namespace string) (DataDogClient, error) {
	initOnce.Do(Init)

	if packageDependencies == nil {
		packageDependencies = &dependencies{}
		packageDependencies.New = statsd.New
	}
	client, err := packageDependencies.New(fmt.Sprintf("%s:%d", DataDogConfig.Host, DataDogConfig.Port))
	if err != nil {
		return nil, err
	}
	client.Namespace = namespace
	client.Tags = append(client.Tags, *DataDogDefaultTags...)
	return &dataDogClientImpl{
		client,
		namespace,
	}, nil
}
