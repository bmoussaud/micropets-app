package internal

import (
	"os"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestLoadConfig(t *testing.T) {

	os.Setenv("SERVICE_CONFIG_DIR", "test")

	config := LoadConfiguration()
	assertEqual(t, ":9001", config.Service.Port)
	if config.Observability.Application != "" {
		assertEqual(t, "micropets", config.Observability.Application)
	}
	if config.Observability.Service != "" {
		assertEqual(t, "pets", config.Observability.Service)
	}
	assertEqual(t, "us-west", config.Observability.Cluster)
	assertEqual(t, true, config.Observability.Enable)
	assertEqual(t, "RANDOM_NUMBER", config.Service.Mode)
	assertEqual(t, 10, config.Service.FrequencyError)
	assertEqual(t, 100, config.Service.Delay.Period)
	assertEqual(t, 0.3, config.Service.Delay.Amplitude)

}
