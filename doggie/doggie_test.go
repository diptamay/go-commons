package doggie

import (
	"fmt"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	Chance "github.com/ZeFort/chance"
	"github.com/diptamay/go-commons/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFormatTags(t *testing.T) {
	initOnce.Do(Init)

	chance := Chance.New()
	tagsMap := map[string]string{
		"a": chance.Word(),
		"b": chance.Word(),
	}
	result := formatTags(tagsMap)
	assert.Condition(t, func() bool {
		hasTags := (result[0] == fmt.Sprintf("a:%s", tagsMap["a"]) || result[0] == fmt.Sprintf("b:%s", tagsMap["b"]))
		if hasTags {
			switch result[0] {
			case fmt.Sprintf("a:%s", tagsMap["a"]):
				return result[1] == fmt.Sprintf("b:%s", tagsMap["b"])
			case fmt.Sprintf("b:%s", tagsMap["b"]):
				return result[1] == fmt.Sprintf("a:%s", tagsMap["a"])
			default:
				return false
			}
		}
		return false
	}, "should generate correct tag from tagsMap")
}

func TestIncrement(t *testing.T) {
	initOnce.Do(Init)

	chance := Chance.New()
	client := &mocks.MockStatsDClient{}
	metric := chance.String()
	tagsMap := map[string]string{
		"a": chance.Word(),
		"b": chance.Word(),
	}
	client.
		On("Incr", metric, mock.Anything, defaultSampleRate).
		Return(nil)
	doggieClient := &dataDogClientImpl{
		client,
		chance.Word(),
	}
	doggieClient.Increment(metric, tagsMap)
	client.AssertExpectations(t)
}

func TestHistogram(t *testing.T) {
	initOnce.Do(Init)

	chance := Chance.New()
	client := &mocks.MockStatsDClient{}
	metric := chance.String()
	value := chance.Float()
	tagsMap := map[string]string{
		"a": chance.Word(),
		"b": chance.Word(),
	}
	client.
		On("Histogram", metric, value, mock.Anything, defaultSampleRate).
		Return(nil)
	doggieClient := &dataDogClientImpl{
		client,
		chance.Word(),
	}
	doggieClient.Histogram(metric, value, tagsMap)
	client.AssertExpectations(t)
}

func TestGauge(t *testing.T) {
	initOnce.Do(Init)

	chance := Chance.New()
	client := &mocks.MockStatsDClient{}
	metric := chance.String()
	value := chance.Float()
	tagsMap := map[string]string{
		"a": chance.Word(),
		"b": chance.Word(),
	}
	client.
		On("Gauge", metric, value, mock.Anything, defaultSampleRate).
		Return(nil)
	doggieClient := &dataDogClientImpl{
		client,
		chance.Word(),
	}
	doggieClient.Gauge(metric, value, tagsMap)
	client.AssertExpectations(t)
}

func TestTiming(t *testing.T) {
	initOnce.Do(Init)

	chance := Chance.New()
	client := &mocks.MockStatsDClient{}
	metric := chance.String()
	value := time.Duration(chance.Int())
	tagsMap := map[string]string{
		"a": chance.Word(),
		"b": chance.Word(),
	}
	client.
		On("Timing", metric, value, mock.Anything, defaultSampleRate).
		Return(nil)
	doggieClient := &dataDogClientImpl{
		client,
		chance.Word(),
	}
	doggieClient.Timing(metric, value, tagsMap)
	client.AssertExpectations(t)
}

type mockStatsD struct {
	mock.Mock
}

func (s *mockStatsD) New(addr string, options ...statsd.Option) (*statsd.Client, error) {
	s.Called(addr, options)
	return &statsd.Client{}, nil
}

//func TestMakeClient(t *testing.T) {
//	initOnce.Do(Init)
//
//	chance := Chance.New()
//	statsdMock := &mockStatsD{}
//	if packageDependencies == nil {
//		packageDependencies = &dependencies{}
//	}
//	packageDependencies.New = statsdMock.New
//	DataDogConfig.Host = "127.0.0.1"
//	DataDogConfig.Port = chance.IntBtw(1000, 9999)
//
//	statsdMock.
//		On("New", fmt.Sprintf("%s:%d", DataDogConfig.Host, DataDogConfig.Port), nil).
//		Return(mock.Anything, nil)
//	namespace := chance.Word()
//	dataDogClient, _ := MakeClient(namespace)
//	assert.Equal(t, namespace, dataDogClient.Namespace(), "should initialize with the defined namespace")
//	statsdMock.AssertExpectations(t)
//}
