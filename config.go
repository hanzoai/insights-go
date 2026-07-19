package insights

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// CompressionMode specifies the compression algorithm for batch payloads.
type CompressionMode uint8

const (
	// CompressionNone disables compression (default).
	CompressionNone CompressionMode = 0
	// CompressionGzip enables GZIP compression for batch payloads.
	CompressionGzip CompressionMode = 1
)

// Config carries configuration options used when constructing a Client with NewWithConfig.
//
// Each exported field's zero value is either meaningful or replaced by the
// default documented on that field.
type Config struct {

	// Endpoint is the PostHog API host used for event ingestion and feature flag requests.
	// If empty, it defaults to DefaultEndpoint.
	Endpoint string

	// Specifying a Personal API key will make feature flag evaluation more performant,
	// but it's not required for feature flags.  If you don't have a personal API key,
	// you can leave this field empty, and all of the relevant feature flag evaluation
	// methods will still work.
	// Information on how to get a personal API key: https://insights.hanzo.ai/docs/api/overview
	PersonalApiKey string

	// DisableGeoIP controls whether event and feature flag requests include
	// $geoip_disable/geoip_disable. Nil defaults to true because this SDK usually
	// runs server-side; set Ptr(false) to allow GeoIP lookup.
	DisableGeoIP *bool

	// Interval is the flush interval for queued messages. Messages are sent when
	// BatchSize is reached or when this timer fires. If zero, it defaults to
	// DefaultInterval.
	Interval time.Duration

	// DefaultFeatureFlagsPollingInterval is the interval for reloading local feature
	// flag definitions when PersonalApiKey is configured. If zero, it defaults to
	// DefaultFeatureFlagsPollingInterval.
	DefaultFeatureFlagsPollingInterval time.Duration

	// FeatureFlagRequestTimeout is the timeout for feature flag and remote config
	// HTTP requests. If zero, it defaults to DefaultFeatureFlagRequestTimeout.
	// Use time.Duration values such as 3 * time.Second.
	FeatureFlagRequestTimeout time.Duration

	// NextFeatureFlagsPollingTick optionally calculates the next local feature flag
	// polling delay. When set, it overrides DefaultFeatureFlagsPollingInterval.
	NextFeatureFlagsPollingTick func() time.Duration

	// Flag to enable historical migration
	// See more in our migration docs: https://insights.hanzo.ai/docs/migrate
	HistoricalMigration bool

	// Transport is the HTTP transport used by the client. Set it to customize
	// low-level request behavior such as connection pooling or proxies. If nil,
	// the client uses a clone of http.DefaultTransport with SDK defaults.
	Transport http.RoundTripper

	// Logger receives informational, warning, and error messages from background
	// operations. If nil, the client logs to os.Stderr with the standard logger.
	Logger Logger

	// DefaultEventProperties are merged into every Capture event before sending.
	// They are useful for common metadata like service name or app version. On key
	// conflicts, values from DefaultEventProperties overwrite event properties.
	DefaultEventProperties Properties

	// Callback receives success or failure notifications for messages sent to the
	// PostHog batch API.
	Callback Callback

	// BatchSize is the maximum number of messages sent in one batch API call.
	// Messages are sent when BatchSize is reached or when Interval fires. If zero,
	// it defaults to DefaultBatchSize. The API still enforces a 500KB request limit.
	BatchSize int

	// Verbose enables more frequent and detailed debug logging through Logger.
	Verbose bool

	// RetryAfter returns the delay before retrying a failed batch upload. The int
	// argument is the retry attempt number. If nil, DefaultBackoff().Duration is used.
	RetryAfter func(int) time.Duration

	// MaxRetries is the maximum number of retries after the first send attempt.
	// It must be in [0,9]. If nil, it defaults to 9 retries (10 total attempts).
	MaxRetries *int

	// ShutdownTimeout is the maximum time Close waits for in-flight messages to be
	// sent. If zero or negative, Close waits indefinitely for backward compatibility.
	ShutdownTimeout time.Duration

	// BatchUploadTimeout is the timeout for uploading one batch to the /batch/
	// endpoint. If zero, it defaults to DefaultBatchUploadTimeout.
	BatchUploadTimeout time.Duration

	// BatchSubmitTimeout is the maximum time to wait when submitting a batch to the
	// worker pool while its queue is full. If zero, it defaults to
	// DefaultBatchSubmitTimeout. Set a negative duration for non-blocking behavior
	// that drops immediately when the queue is full.
	BatchSubmitTimeout time.Duration

	// MaxEnqueuedRequests is the maximum number of batches waiting for upload.
	// When the queue is full, new batches are dropped and the failure callback is
	// invoked. If zero, it defaults to DefaultMaxEnqueuedRequests.
	MaxEnqueuedRequests int

	// Compression selects the compression mode for batch payloads. CompressionGzip
	// compresses payloads and adds the appropriate headers/query params. If zero,
	// it defaults to CompressionNone.
	Compression CompressionMode

	// A function called by the client to get the current time, `time.Now` is
	// used by default.
	// This field is not exported and only exposed internally to control concurrency.
	now func() time.Time

	// maxAttempts is a maximum numbers we try to send data to capture endpoint, must be in range [1,10].
	maxAttempts int
}

// GetDisableGeoIP instructs the client to set $geoip_disable on event properties or feature flag requests.
// It is on by default as Go is mainly used on server side.
func (c Config) GetDisableGeoIP() bool {
	return c.DisableGeoIP == nil || *c.DisableGeoIP
}

const (
	SDKName = "insights-go"

	// DefaultEndpoint is the default HOST events and flag requests are sent to when
	// none is set explicitly (Config.Endpoint) or via EndpointEnvVar. It is the live
	// Hanzo Cloud Go ingest (api.hanzo.ai), NOT the retired PostHog-OSS capture host
	// (insights.hanzo.ai, replicas:0) — events posted there were black-holed.
	// capturePath / flagsPath below are appended to it.
	DefaultEndpoint = "https://api.hanzo.ai"

	// EndpointEnvVar overrides DefaultEndpoint at runtime when Config.Endpoint is
	// empty, so an operator can repoint ingest without a rebuild (regional/staging
	// host). An explicit Config.Endpoint always wins over the env.
	EndpointEnvVar = "INSIGHTS_ENDPOINT"

	// capturePath is the cloud native ingest route the batch uploader POSTs to —
	// the PostHog-compatible front door (cloud clients/analytics insightsIngest),
	// which reads {api_key, batch:[...]}. Replaces PostHog's "/batch/".
	capturePath = "/v1/insights/e"

	// flagsPath is the cloud native flags engine route (cloud clients/flags:
	// POST /v1/flags, /v1/flags/decide alias). Replaces PostHog's "/flags/".
	flagsPath = "/v1/flags"

	// DefaultInterval is the default flush interval used when Config.Interval is zero.
	DefaultInterval = 5 * time.Second

	// DefaultFeatureFlagsPollingInterval is the default local feature flag reload interval.
	DefaultFeatureFlagsPollingInterval = 5 * time.Minute

	// DefaultFeatureFlagRequestTimeout is the default timeout for feature flag requests.
	DefaultFeatureFlagRequestTimeout = 3 * time.Second

	// DefaultBatchSize is the default batch size used when Config.BatchSize is zero.
	DefaultBatchSize = 250

	// DefaultBatchUploadTimeout is the default timeout for uploading batched
	// events to the /batch/ endpoint.
	DefaultBatchUploadTimeout = 10 * time.Second

	// DefaultBatchSubmitTimeout is the default timeout for submitting batches
	// to the worker pool when the queue is full. This allows workers time to
	// complete during transient latency spikes, reducing unnecessary data loss.
	DefaultBatchSubmitTimeout = 100 * time.Millisecond

	// DefaultMaxEnqueuedRequests is the default maximum number of batches that
	// can be queued for sending.
	DefaultMaxEnqueuedRequests = 1000
)

func (c *Config) normalize() {
	// TrimRight "/" so capturePath / flagsPath (both leading-slash) join without a
	// double slash — a bare host never legitimately ends in "/", and Traefik matches
	// the ingest path exactly, so "//v1/insights/e" would miss the route.
	c.Endpoint = strings.TrimRight(strings.TrimSpace(c.Endpoint), "/")
	c.PersonalApiKey = strings.TrimSpace(c.PersonalApiKey)
}

// endpointDefault resolves the ingest host when Config.Endpoint is unset: the
// EndpointEnvVar override if present, else DefaultEndpoint. Trailing slashes are
// trimmed for the same clean-join reason as normalize.
func endpointDefault() string {
	if v := strings.TrimRight(strings.TrimSpace(os.Getenv(EndpointEnvVar)), "/"); v != "" {
		return v
	}
	return DefaultEndpoint
}

// Validate verifies that fields that don't have zero-values are set to valid values,
// returns an error describing the problem if a field was invalid.
func (c *Config) Validate() error {
	c.normalize()

	if c.Interval < 0 {
		return ConfigError{
			Reason: "negative time intervals are not supported",
			Field:  "Interval",
			Value:  c.Interval,
		}
	}

	if c.BatchSize < 0 {
		return ConfigError{
			Reason: "negative batch sizes are not supported",
			Field:  "BatchSize",
			Value:  c.BatchSize,
		}
	}

	if _, err := url.Parse(c.Endpoint); err != nil {
		return ConfigError{
			Reason: "invalid endpoint",
			Field:  "Endpoint",
			Value:  c.Endpoint,
		}
	}

	if c.MaxRetries != nil && (*c.MaxRetries < 0 || *c.MaxRetries > 9) {
		return ConfigError{
			Reason: "max retries out of range [0,9]",
			Field:  "MaxRetries",
			Value:  *c.MaxRetries,
		}
	}

	if c.Compression > CompressionGzip {
		return ConfigError{
			Reason: "invalid compression mode",
			Field:  "Compression",
			Value:  c.Compression,
		}
	}

	return nil
}

// Given a config object as argument the function will set all zero-values to
// their defaults and return the modified object.
func makeConfig(c Config) Config {
	c.normalize()

	if len(c.Endpoint) == 0 {
		c.Endpoint = endpointDefault()
	}

	if c.Interval == 0 {
		c.Interval = DefaultInterval
	}

	if c.DefaultFeatureFlagsPollingInterval == 0 {
		c.DefaultFeatureFlagsPollingInterval = DefaultFeatureFlagsPollingInterval
	}

	if c.FeatureFlagRequestTimeout == 0 {
		c.FeatureFlagRequestTimeout = DefaultFeatureFlagRequestTimeout
	}

	// Note: c.Transport == nil is handled by makeHttpClient() which clones
	// DefaultTransport with tuned connection pool settings

	if c.Logger == nil {
		c.Logger = newDefaultLogger(c.Verbose)
	}

	if c.BatchSize == 0 {
		c.BatchSize = DefaultBatchSize
	}

	if c.RetryAfter == nil {
		c.RetryAfter = DefaultBackoff().Duration
	}

	if c.now == nil {
		c.now = time.Now
	}

	if c.MaxRetries != nil && 0 <= *c.MaxRetries && *c.MaxRetries <= 9 {
		c.maxAttempts = 1 + *c.MaxRetries
	} else {
		c.maxAttempts = 10
	}

	// Note: ShutdownTimeout == 0 means wait indefinitely (backward compatible).
	// Users opt-in to timeout by setting a positive duration.

	if c.BatchUploadTimeout == 0 {
		c.BatchUploadTimeout = DefaultBatchUploadTimeout
	}

	if c.BatchSubmitTimeout == 0 {
		c.BatchSubmitTimeout = DefaultBatchSubmitTimeout
	}

	if c.MaxEnqueuedRequests == 0 {
		c.MaxEnqueuedRequests = DefaultMaxEnqueuedRequests
	}

	if c.GetDisableGeoIP() {
		if c.DefaultEventProperties == nil {
			c.DefaultEventProperties = NewProperties()
		}
		c.DefaultEventProperties.Set(propertyGeoipDisable, true)
	}

	return c
}

// Ptr returns a pointer to v.
func Ptr[T any](v T) *T {
	return &v
}
