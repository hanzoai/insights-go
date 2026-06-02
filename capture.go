package insights

import (
	"time"
)

// SendFeatureFlagsValue defines the legacy Capture.SendFeatureFlags configuration interface.
type SendFeatureFlagsValue interface {
	// ShouldSend reports whether Enqueue should attach feature flag values to the Capture event.
	ShouldSend() bool
	// GetOptions returns optional feature flag evaluation parameters, or nil for boolean-only values.
	GetOptions() *SendFeatureFlagsOptions
}

// SendFeatureFlagsOptions allows granular control over legacy feature flag evaluation for Capture.
// Deprecated: Prefer Client.EvaluateFlags and Capture.Flags so captures use the exact
// flag snapshot your application branched on without another /flags request.
type SendFeatureFlagsOptions struct {
	// OnlyEvaluateLocally forces evaluation to use locally loaded flag definitions only.
	OnlyEvaluateLocally bool
	// DeviceId provides a device_id for remote flag evaluation requests.
	DeviceId *string
	// PersonProperties provides explicit person properties for flag evaluation.
	PersonProperties Properties
	// GroupProperties provides explicit group properties for flag evaluation.
	GroupProperties map[string]Properties
}

// ShouldSend reports true when opts is non-nil.
func (opts *SendFeatureFlagsOptions) ShouldSend() bool {
	return opts != nil
}

// GetOptions returns opts so Capture can pass the options to feature flag evaluation.
func (opts *SendFeatureFlagsOptions) GetOptions() *SendFeatureFlagsOptions {
	return opts
}

// SendFeatureFlagsBool wraps a boolean value to implement SendFeatureFlagsValue.
type SendFeatureFlagsBool bool

// ShouldSend reports the wrapped boolean value.
func (b SendFeatureFlagsBool) ShouldSend() bool {
	return bool(b)
}

// GetOptions returns nil because SendFeatureFlagsBool has no additional options.
func (b SendFeatureFlagsBool) GetOptions() *SendFeatureFlagsOptions {
	return nil
}

// SendFeatureFlags converts enabled into a legacy Capture.SendFeatureFlags value.
// Deprecated: Prefer Client.EvaluateFlags and Capture.Flags.
func SendFeatureFlags(enabled bool) SendFeatureFlagsValue {
	return SendFeatureFlagsBool(enabled)
}

// SendFeatureFlagsWithOptions wraps legacy feature flag capture options.
// Deprecated: Prefer Client.EvaluateFlags and Capture.Flags.
func SendFeatureFlagsWithOptions(opts *SendFeatureFlagsOptions) SendFeatureFlagsValue {
	return opts
}

// Helper functions to work with SendFeatureFlags interface
func (c *Capture) shouldSendFeatureFlags() bool {
	if c.SendFeatureFlags == nil {
		return false
	}
	return c.SendFeatureFlags.ShouldSend()
}

func (c *Capture) getFeatureFlagsOptions() *SendFeatureFlagsOptions {
	if c.SendFeatureFlags == nil {
		return nil
	}
	return c.SendFeatureFlags.GetOptions()
}

var _ Message = (*Capture)(nil)

// Capture represents a custom event to send to PostHog.
// Enqueue validates Event and DistinctId (or obtains DistinctId from request context),
// fills Type, Uuid, and Timestamp, merges Config.DefaultEventProperties, and queues
// the event for a future batch upload.
type Capture struct {
	// Type is reserved for SDK serialization and is overwritten by Enqueue.
	Type string
	// You don't usually need to specify this field - Insights will generate it automatically.
	// Use it only when necessary - for example, to prevent duplicate events.
	Uuid             string
	DistinctId       string
	Event            string
	Timestamp        time.Time
	Properties       Properties
	Groups           Groups
	SendFeatureFlags SendFeatureFlagsValue
	// Flags, when set, attaches $feature/<key> and $active_feature_flags
	// properties from a snapshot returned by Client.EvaluateFlags. It is
	// preferred over SendFeatureFlags: the snapshot guarantees the event
	// carries the exact values the application branched on and avoids a
	// hidden /flags request on every capture. Flags takes precedence when
	// both are set.
	Flags *FeatureFlagEvaluations
}

func (msg Capture) internal() {
	panic(unimplementedError)
}

// Validate checks that the capture message has a non-empty Event and DistinctId.
func (msg Capture) Validate() error {
	if len(msg.Event) == 0 {
		return FieldError{
			Type:  "insights.Capture",
			Name:  "Event",
			Value: msg.Event,
		}
	}

	if len(msg.DistinctId) == 0 {
		return FieldError{
			Type:  "insights.Capture",
			Name:  "DistinctId",
			Value: msg.DistinctId,
		}
	}

	return nil
}

func validateCaptureEvent(msg Capture) error {
	if len(msg.Event) == 0 {
		return FieldError{
			Type:  "posthog.Capture",
			Name:  "Event",
			Value: msg.Event,
		}
	}

	return nil
}

// CaptureInApi is the wire-format payload produced from a Capture message.
type CaptureInApi struct {
	// Type is the message type sent to the batch API.
	Type string `json:"type"`
	// Uuid is the event UUID sent to the batch API.
	Uuid string `json:"uuid"`
	// Library is the SDK name sent to the batch API.
	Library string `json:"library"`
	// LibraryVersion is the SDK version sent to the batch API.
	LibraryVersion string `json:"library_version"`
	// Timestamp is the event timestamp sent to the batch API.
	Timestamp time.Time `json:"timestamp"`

	// DistinctId identifies the user or entity for the event.
	DistinctId string `json:"distinct_id"`
	// Event is the captured event name.
	Event string `json:"event"`
	// Properties contains event, SDK, system, group, and feature flag properties.
	Properties Properties `json:"properties"`
	// SendFeatureFlags carries the legacy send_feature_flags value for compatibility.
	SendFeatureFlags SendFeatureFlagsValue `json:"send_feature_flags"`
}

// APIfy converts a Capture message into the PostHog batch API representation.
func (msg Capture) APIfy() APIMessage {
	libraryVersion := getVersion()

	myProperties := Properties{}.
		Merge(msg.Properties).
		Set("$lib", SDKName).
		Set("$lib_version", libraryVersion).
		Merge(getSystemContext().ToProperties())

	if msg.Groups != nil {
		myProperties.Set("$groups", msg.Groups)
	}

	apified := CaptureInApi{
		Type:             msg.Type,
		Uuid:             msg.Uuid,
		Library:          SDKName,
		LibraryVersion:   libraryVersion,
		Timestamp:        msg.Timestamp,
		DistinctId:       msg.DistinctId,
		Event:            msg.Event,
		Properties:       myProperties,
		SendFeatureFlags: msg.SendFeatureFlags,
	}

	return apified
}
