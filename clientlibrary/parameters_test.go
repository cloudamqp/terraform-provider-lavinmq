package clientlibrary

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParameterResponseSanitizedRedactsURICredentials(t *testing.T) {
	body := []byte(`{"name":"s1","vhost":"/","component":"shovel","value":{
		"src-uri":"amqp://user:hunter2@src:5672/%2f",
		"dest-uri":"amqps://guest:guest@dst",
		"uri":"amqp://fed:s3cret@upstream/%2f",
		"src-queue":"q1"}}`)
	var result *ParameterResponse
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatal(err)
	}

	out := (&service{name: "parameters"}).DataLog("Get", "api/parameters/shovel/%2f/s1", result.Sanitized())

	for _, secret := range []string{"hunter2", "guest:guest", "s3cret"} {
		if strings.Contains(out, secret) {
			t.Errorf("sanitized log leaks %q: %s", secret, out)
		}
	}
	for _, kept := range []string{"src:5672", "dst", "upstream", "src-queue:q1", "Component:shovel"} {
		if !strings.Contains(out, kept) {
			t.Errorf("sanitized log lost non-secret %q: %s", kept, out)
		}
	}
	if !strings.Contains(result.Value.(map[string]any)["src-uri"].(string), "hunter2") {
		t.Error("Sanitized must not mutate the receiver")
	}
}

func TestParameterResponseSanitizedHandlesNilAndNonMap(t *testing.T) {
	var nilResp *ParameterResponse
	if nilResp.Sanitized() != nil {
		t.Error("nil receiver should stay nil")
	}
	plain := &ParameterResponse{Name: "n", Value: "not-a-map"}
	if plain.Sanitized().Value != "not-a-map" {
		t.Error("non-map value should pass through unchanged")
	}
}
