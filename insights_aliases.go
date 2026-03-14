package insights

// Insights is an alias for Client, provided for convenience.
// Both insights.Insights and insights.Client refer to the same interface.
type Insights = Client

// InsightsConfig is an alias for Config, provided for convenience.
type InsightsConfig = Config

// NewInsights creates a new Insights client. Equivalent to New.
func NewInsights(apiKey string) Insights {
	return New(apiKey)
}

// NewInsightsWithConfig creates a new Insights client with config. Equivalent to NewWithConfig.
func NewInsightsWithConfig(apiKey string, config InsightsConfig) (Insights, error) {
	return NewWithConfig(apiKey, config)
}
