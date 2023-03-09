package internal

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/wavefronthq/wavefront-sdk-go/application"
	"github.com/wavefronthq/wavefront-sdk-go/senders"

	"github.com/wavefronthq/wavefront-opentracing-sdk-go/reporter"
	wfTracer "github.com/wavefronthq/wavefront-opentracing-sdk-go/tracer"

	"github.com/opentracing/opentracing-go"
	otrext "github.com/opentracing/opentracing-go/ext"
)

func NewGlobalTracer() io.Closer {

	var sender senders.Sender
	var err error

	ObservabilityConfig := LoadConfiguration().Observability
	if ObservabilityConfig.Enable {

		if ObservabilityConfig.Server != "" && ObservabilityConfig.Token != "" {
			config := &senders.DirectConfiguration{
				Server: ObservabilityConfig.Server,
				Token:  ObservabilityConfig.Token,
			}
			sender, err = senders.NewDirectSender(config)
			if err != nil {
				log.Fatalf("error creating wavefront sender: %q", err)
			}

			log.Printf("* Enabled Observability on %s \n", ObservabilityConfig.Server)

		} else {
			log.Fatalf("Not enough configuration parameter has been specified for sender.")
		}

		appTags := application.New(ObservabilityConfig.Application, ObservabilityConfig.Service)
		log.Printf("* Enabled Observability new Application  %s/%s \n", ObservabilityConfig.Application, ObservabilityConfig.Service)
		log.Printf("* Enabled Observability on Cluster  %s/%s \n", ObservabilityConfig.Cluster, ObservabilityConfig.Shard)

		appTags.Cluster = ObservabilityConfig.Cluster
		appTags.Shard = ObservabilityConfig.Shard

		var spanReporter reporter.WavefrontSpanReporter
		if ObservabilityConfig.Source != "" {
			spanReporter = reporter.New(sender, appTags, reporter.Source(ObservabilityConfig.Source))
		} else {
			spanReporter = reporter.New(sender, appTags)
			log.Printf("* Enabled Observability with default span reporter")
		}

		consoleReporter := reporter.NewConsoleSpanReporter(ObservabilityConfig.Service)
		compositeReporter := reporter.NewCompositeSpanReporter(spanReporter, consoleReporter)
		wavefrontTracer := wfTracer.New(compositeReporter)
		opentracing.SetGlobalTracer(wavefrontTracer)

	}
	return ioutil.NopCloser(nil)

}

func NewServerSpan(req *http.Request, spanName string) opentracing.Span {
	tracer := opentracing.GlobalTracer()
	parentCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	var span opentracing.Span
	if err == nil { // has parent context
		span = tracer.StartSpan(spanName, opentracing.ChildOf(parentCtx))
	} else if err == opentracing.ErrSpanContextNotFound { // no parent
		span = tracer.StartSpan(spanName)
	} else {
		log.Printf("Error in extracting tracer context: %s", err.Error())
	}
	otrext.SpanKindRPCServer.Set(span)
	return span
}
