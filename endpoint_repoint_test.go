package insights

import "testing"

// TestEndpointRepointDefault: with neither an explicit Config.Endpoint nor the env
// override, the SDK targets the live Hanzo Cloud Go ingest, and the capture/flags
// paths are the cloud routes.
func TestEndpointRepointDefault(t *testing.T) {
	if DefaultEndpoint != "https://api.hanzo.ai" {
		t.Fatalf("DefaultEndpoint = %q, want https://api.hanzo.ai", DefaultEndpoint)
	}
	if capturePath != "/v1/insights/e" {
		t.Fatalf("capturePath = %q, want /v1/insights/e", capturePath)
	}
	if flagsPath != "/v1/flags" {
		t.Fatalf("flagsPath = %q, want /v1/flags", flagsPath)
	}
	if got := makeConfig(Config{}).Endpoint; got != "https://api.hanzo.ai" {
		t.Fatalf("resolved endpoint = %q, want https://api.hanzo.ai", got)
	}
}

// TestEndpointEnvOverride: INSIGHTS_ENDPOINT repoints ingest when Config.Endpoint
// is unset (and trailing slashes are trimmed so capturePath joins cleanly); an
// explicit Config.Endpoint always wins over the env.
func TestEndpointEnvOverride(t *testing.T) {
	t.Setenv(EndpointEnvVar, "https://staging.hanzo.ai/")
	if got := makeConfig(Config{}).Endpoint; got != "https://staging.hanzo.ai" {
		t.Fatalf("env-override endpoint = %q, want https://staging.hanzo.ai (slash trimmed)", got)
	}
	if got := makeConfig(Config{Endpoint: "https://explicit.example.com/"}).Endpoint; got != "https://explicit.example.com" {
		t.Fatalf("explicit endpoint = %q, want it to win over env (slash trimmed)", got)
	}
}
