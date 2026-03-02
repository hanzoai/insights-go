package posthog

// Insights is an alias for Client - use Insights for new code.
// The underlying type is the posthog Client interface.
type Insights = Client

// InsightsConfig is an alias for Config - use InsightsConfig for new code.
type InsightsConfig = Config

// NewInsights creates a new Insights client. Equivalent to New.
func NewInsights(apiKey string) Insights {
	return New(apiKey)
}

// NewInsightsWithConfig creates a new Insights client with config. Equivalent to NewWithConfig.
func NewInsightsWithConfig(apiKey string, config InsightsConfig) (Insights, error) {
	return NewWithConfig(apiKey, config)
}
