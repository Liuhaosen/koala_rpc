package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"modtest/gostudy/lesson2/ibinarytree/koala/logs"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	"google.golang.org/grpc/metadata"
)

const (
	bindHdrSuffix = "-bin"
)

type metadataTextMap metadata.MD

func (m metadataTextMap) Set(key, val string) {
	encodeKey, encodeVal := encodeKeyValue(key, val)
	m[encodeKey] = []string{encodeVal}
}

func (m metadataTextMap) ForeachKey(callback func(key, val string) error) error {
	for k, vv := range m {
		for _, v := range vv {
			if decodeKey, decodeVal, err := metadata.DecodeKeyValue(k, v); err == nil {
				if err = callback(decodeKey, decodeVal); err != nil {
					return err
				}

			} else {
				return fmt.Errorf("failed decoding opentracing from gRPC metadata : %v", err)
			}
		}
	}
	return nil
}

func encodeKeyValue(k, v string) (string, string) {
	k = strings.ToLower(k)
	if strings.HasSuffix(k, bindHdrSuffix) {
		val := base64.StdEncoding.EncodeToString([]byte(v))
		v = string(val)
	}
	return k, v
}

func InitTrace(serviceName, reportAddr, sampleType string, sampleRate float64) (err error) {
	transport, err := zipkin.NewHTTPTransport(
		reportAddr,
		zipkin.HTTPBatchSize(16),
		zipkin.HTTPLogger(jaeger.StdLogger),
	)
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init zipkin, err: %v", err)
		return
	}
	//配置
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  sampleType,
			Param: sampleRate,
		},
	}

	//上报实例
	r := jaeger.NewRemoteReporter(transport)
	tracer, closer, err := cfg.New(serviceName, config.Logger(jaeger.StdLogger), config.Reporter(r))

	// tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger), config.Reporter(r))
	cfg.ServiceName = serviceName
	if err != nil {
		logs.Error(context.TODO(), "ERROR: cannot init jaeger : %v", err)
		return
	}

	_ = closer
	opentracing.SetGlobalTracer(tracer)
	return
}
