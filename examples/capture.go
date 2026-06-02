package main

import (
	"fmt"
	"time"

	"github.com/hanzoai/insights-go"
)

func TestCapture(projectAPIKey, endpoint string) {
	fmt.Println("📊 Capturing basic events...")

	client, _ := insights.NewWithConfig(projectAPIKey, insights.Config{
		Interval:  30 * time.Second,
		BatchSize: 100,
		Verbose:   true,
		Endpoint:  endpoint,
	})
	defer client.Close()

	// Send a few different types of events
	fmt.Println("→ Sending 'Download' event...")
	if err := client.Enqueue(insights.Capture{
		Event:      "Download",
		DistinctId: "user_123456",
		Properties: map[string]interface{}{
			"application": "Insights Go",
			"version":     "1.0.0",
			"file_size":   "2.5MB",
		},
	}); err != nil {
		fmt.Println("❌ Error sending Download event:", err)
		return
	}

	fmt.Println("→ Sending 'Page View' event...")
	if err := client.Enqueue(insights.Capture{
		Event:      "$pageview",
		DistinctId: "user_123456",
		Properties: map[string]interface{}{
			"$current_url": "https://example.com/dashboard",
			"$title":       "Dashboard - Insights",
			"$referrer":    "https://google.com",
		},
	}); err != nil {
		fmt.Println("❌ Error sending Page View event:", err)
		return
	}

	fmt.Println("→ Sending 'Button Clicked' event...")
	if err := client.Enqueue(insights.Capture{
		Event:      "Button Clicked",
		DistinctId: "user_123456",
		Properties: map[string]interface{}{
			"button_text": "Sign Up",
			"page":        "/landing",
			"experiment":  "homepage_test_v2",
		},
	}); err != nil {
		fmt.Println("❌ Error sending Button Clicked event:", err)
		return
	}

	// Give the client time to send events
	time.Sleep(1 * time.Second)
	fmt.Println("✅ Basic events sent successfully!")
}

func TestCaptureWithSendFeatureFlagOption(projectAPIKey, personalAPIKey, endpoint string) {
	fmt.Println("🏁 Capturing events with feature flags...")
	fmt.Println("   This demonstrates how to automatically include feature flag states with events")

	client, _ := insights.NewWithConfig(projectAPIKey, insights.Config{
		Interval:       30 * time.Second,
		BatchSize:      100,
		Verbose:        true,
		PersonalApiKey: personalAPIKey,
		Endpoint:       endpoint,
	})
	defer client.Close()

	fmt.Println("→ Sending event with SendFeatureFlags enabled...")
	if err := client.Enqueue(insights.Capture{
		Event:      "Purchase",
		DistinctId: "user_123456",
		Properties: map[string]interface{}{
			"amount":   99.99,
			"currency": "USD",
			"product":  "Premium Plan",
		},
		SendFeatureFlags: insights.SendFeatureFlags(true),
	}); err != nil {
		fmt.Println("❌ Error sending Purchase event:", err)
		return
	}

	fmt.Println("→ Sending event without feature flags for comparison...")
	if err := client.Enqueue(insights.Capture{
		Event:      "Login",
		DistinctId: "user_123456",
		Properties: map[string]interface{}{
			"method":     "google",
			"first_time": false,
		},
		// SendFeatureFlags not specified (defaults to false)
	}); err != nil {
		fmt.Println("❌ Error sending Login event:", err)
		return
	}

	// Give the client time to send events
	time.Sleep(1 * time.Second)
	fmt.Println("✅ Events with feature flag states sent successfully!")
	fmt.Println("   ℹ️ The first event will include all active feature flag states for the user")
	fmt.Println("   ℹ️ The second event will not include feature flag information")
}

func TestCaptureWithSendFeatureFlagsOptions(projectAPIKey, personalAPIKey, endpoint string) {
	fmt.Println("🚀 Advanced feature flags with SendFeatureFlagsOptions...")
	fmt.Println("   This demonstrates advanced feature flag evaluation with custom properties")

	client, _ := insights.NewWithConfig(projectAPIKey, insights.Config{
		Interval:       30 * time.Second,
		BatchSize:      100,
		Verbose:        true,
		PersonalApiKey: personalAPIKey,
		Endpoint:       endpoint,
	})
	defer client.Close()

	fmt.Println("→ Sending event with custom person properties for flag evaluation...")
	if err := client.Enqueue(insights.Capture{
		Event:      "Feature Used",
		DistinctId: "premium_user_456",
		Properties: map[string]interface{}{
			"feature_name": "advanced_analytics",
			"usage_count":  1,
		},
		SendFeatureFlags: &insights.SendFeatureFlagsOptions{
			PersonProperties: insights.NewProperties().Set("plan", "premium").Set("beta_user", true),
		},
	}); err != nil {
		fmt.Println("❌ Error sending feature usage event:", err)
		return
	}

	fmt.Println("→ Sending event with local-only evaluation and group properties...")
	if err := client.Enqueue(insights.Capture{
		Event:      "Team Action",
		DistinctId: "enterprise_user_789",
		Properties: map[string]interface{}{
			"action_type": "export_data",
			"data_size":   "50MB",
		},
		SendFeatureFlags: &insights.SendFeatureFlagsOptions{
			OnlyEvaluateLocally: true,
			PersonProperties:    insights.NewProperties().Set("plan", "enterprise").Set("role", "admin"),
			GroupProperties: map[string]insights.Properties{
				"company": insights.NewProperties().Set("name", "Hanzo AI").Set("plan", "enterprise").Set("employees", 100),
			},
		},
	}); err != nil {
		fmt.Println("❌ Error sending team action event:", err)
		return
	}

	fmt.Println("→ Sending event with minimal local evaluation...")
	if err := client.Enqueue(insights.Capture{
		Event:      "Quick Action",
		DistinctId: "basic_user_321",
		Properties: map[string]interface{}{
			"action": "button_click",
			"page":   "homepage",
		},
		SendFeatureFlags: &insights.SendFeatureFlagsOptions{
			OnlyEvaluateLocally: true,
			PersonProperties:    insights.NewProperties().Set("plan", "free"),
		},
	}); err != nil {
		fmt.Println("❌ Error sending quick action event:", err)
		return
	}

	// Give the client time to send events
	time.Sleep(1 * time.Second)
	fmt.Println("✅ Advanced feature flag events sent successfully!")
	fmt.Println("   ℹ️ First event: Custom person properties used for flag evaluation")
	fmt.Println("   ℹ️ Second event: Local-only evaluation with group properties")
	fmt.Println("   ℹ️ Third event: Minimal local evaluation for performance")
}
