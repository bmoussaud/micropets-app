package internal

import (	
	"log"
	"net/http"
)

func NewGlobalTracer() {

	log.Printf("* NewGlobalTracer not implemented TODO move to opentelemetry")
}

func NewServerSpan(req *http.Request, spanName string) {
	log.Printf("* NewGlobalTracer not implemented TODO move to opentelemetry")
}
