// Insights Go library examples
//
// This script demonstrates various Insights Go SDK capabilities including:
// - Basic event capture and user identification
// - Feature flag local evaluation
// - Feature flag payloads
// - Context management
//
// Setup:
// 1. Copy .env.example to .env and fill in your Insights credentials
// 2. Run this script and choose from the interactive menu

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/hanzoai/insights-go" // Used by other files in this package
)

var (
	projectAPIKey  string
	personalAPIKey string
	endpoint       string
)

func init() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Get configuration from environment variables
	projectAPIKey = os.Getenv("INSIGHTS_PROJECT_API_KEY")
	personalAPIKey = os.Getenv("INSIGHTS_PERSONAL_API_KEY")
	endpoint = os.Getenv("INSIGHTS_ENDPOINT")

	if endpoint == "" {
		endpoint = "https://insights.hanzo.ai"
	}
}

func promptForInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		// If we can't read from stdin (e.g., running in CI), return "6" to run all examples
		fmt.Println("6 (auto-selected for non-interactive environment)")
		return "6"
	}
	input = strings.TrimSpace(input)
	if input == "" {
		// If input is empty (e.g., just pressed enter), return "6" to run all examples
		fmt.Println("6 (auto-selected for empty input)")
		return "6"
	}
	return input
}

func checkCredentials() {
	// Check if credentials are provided
	if projectAPIKey == "" || personalAPIKey == "" {
		fmt.Println("Missing Insights credentials!")
		fmt.Println("   Please set INSIGHTS_PROJECT_API_KEY and INSIGHTS_PERSONAL_API_KEY environment variables")
		fmt.Println("   or copy .env.example to .env and fill in your values")
		fmt.Println()

		if projectAPIKey == "" {
			projectAPIKey = promptForInput("Enter your Insights project API key: ")
		}
		if personalAPIKey == "" {
			personalAPIKey = promptForInput("Enter your Insights personal API key: ")
		}
	} else {
		fmt.Println("Insights credentials loaded successfully!")
		fmt.Println("   Project API Key: [REDACTED]")
		fmt.Println("   Personal API Key: [REDACTED]")
		fmt.Printf("   Endpoint: %s\n\n", endpoint)
	}
}

func showMenu() {
	fmt.Println("Insights Go SDK Demo - Choose an example to run:")
	fmt.Println()
	fmt.Println("1. Basic capture examples")
	fmt.Println("2. Capture with feature flags examples")
	fmt.Println("3. Feature flag evaluation examples")
	fmt.Println("4. Feature flag with SendFeatureFlagsOptions examples")
	fmt.Println("5. Flag dependencies examples")
	fmt.Println("6. ETag polling test (continuous, Ctrl+C to stop)")
	fmt.Println("7. Run all examples (except ETag polling)")
	fmt.Println("8. Exit")
}

func runBasicCaptureExamples() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("BASIC CAPTURE EXAMPLES")
	fmt.Println(strings.Repeat("=", 60))
	TestCapture(projectAPIKey, endpoint)
}

func runCaptureWithFeatureFlagsExamples() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("CAPTURE WITH FEATURE FLAGS EXAMPLES")
	fmt.Println(strings.Repeat("=", 60))
	TestCaptureWithSendFeatureFlagOption(projectAPIKey, personalAPIKey, endpoint)
}

func runFeatureFlagEvaluationExamples() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("FEATURE FLAG EVALUATION EXAMPLES")
	fmt.Println(strings.Repeat("=", 60))
	TestIsFeatureEnabled(projectAPIKey, personalAPIKey, endpoint)
}

func runAdvancedFeatureFlagsExamples() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("ADVANCED FEATURE FLAGS (SendFeatureFlagsOptions) EXAMPLES")
	fmt.Println(strings.Repeat("=", 60))
	TestCaptureWithSendFeatureFlagsOptions(projectAPIKey, personalAPIKey, endpoint)
}

func runFlagDependenciesExamples() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("FLAG DEPENDENCIES EXAMPLES")
	fmt.Println(strings.Repeat("=", 60))
	TestFlagDependencies(projectAPIKey, personalAPIKey, endpoint)
}

func runETagPollingExample() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("ETAG POLLING TEST")
	fmt.Println(strings.Repeat("=", 60))
	TestETagPolling(projectAPIKey, personalAPIKey, endpoint)
}

func runAllExamples() {
	fmt.Println("\nRunning all examples...")

	fmt.Println("\n--- BASIC CAPTURE ---")
	TestCapture(projectAPIKey, endpoint)

	fmt.Println("\n--- CAPTURE WITH FEATURE FLAGS ---")
	TestCaptureWithSendFeatureFlagOption(projectAPIKey, personalAPIKey, endpoint)

	fmt.Println("\n--- FEATURE FLAG EVALUATION ---")
	TestIsFeatureEnabled(projectAPIKey, personalAPIKey, endpoint)
	TestErrorTrackingThroughEnqueueing(projectAPIKey, endpoint)
	TestErrorTrackingThroughLogHandler(projectAPIKey, endpoint)

	fmt.Println("\n--- ADVANCED FEATURE FLAGS ---")
	TestCaptureWithSendFeatureFlagsOptions(projectAPIKey, personalAPIKey, endpoint)

	fmt.Println("\n--- FLAG DEPENDENCIES ---")
	TestFlagDependencies(projectAPIKey, personalAPIKey, endpoint)
}

func isInteractive() bool {
	// Check if we're running in an interactive terminal
	fileInfo, _ := os.Stdin.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

func main() {
	checkCredentials()

	// If not interactive, just run all examples
	if !isInteractive() {
		fmt.Println("Non-interactive mode detected. Running all examples...")
		runAllExamples()
		return
	}

	for {
		showMenu()
		choice := promptForInput("\nEnter your choice (1-8): ")

		switch choice {
		case "1":
			runBasicCaptureExamples()
		case "2":
			runCaptureWithFeatureFlagsExamples()
		case "3":
			runFeatureFlagEvaluationExamples()
		case "4":
			runAdvancedFeatureFlagsExamples()
		case "5":
			runFlagDependenciesExamples()
		case "6":
			runETagPollingExample()
			// ETag polling runs continuously, so exit after it returns
			return
		case "7":
			runAllExamples()
		case "8":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select 1-8.")
			continue
		}

		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("Example completed!")
		fmt.Println(strings.Repeat("=", 60))

		// Ask if user wants to run another example
		again := promptForInput("\nWould you like to run another example? (y/N): ")
		if strings.ToLower(again) != "y" && strings.ToLower(again) != "yes" {
			fmt.Println("Goodbye!")
			break
		}
		fmt.Println()
	}
}
