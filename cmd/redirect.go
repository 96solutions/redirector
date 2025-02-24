// Package cmd contains the command-line interface implementations for the redirector service.
// It provides commands for testing redirect rules, viewing configuration, and managing the service.
package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"strings"

	"github.com/lroman242/redirector/config"
	"github.com/lroman242/redirector/domain/dto"
	"github.com/lroman242/redirector/registry"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

// redirectCmd represents the redirect command for testing tracking link redirects.
// It allows simulating redirect requests with various parameters and displays the results.
var redirectCmd = &cobra.Command{
	Use:   "redirect [slug]",
	Short: "Test redirect rules for a tracking link",
	Long: `The redirect command allows testing redirect rules for a specific tracking link.
It simulates a redirect request with customizable parameters like user agent, IP address,
protocol, referrer, and tracking parameters (p1-p4).

Examples:
  # Basic redirect test
  redirector redirect test-slug --ip=192.168.1.1 --ua="Mozilla/5.0" --protocol=https

  # Test with tracking parameters
  redirector redirect test-slug --p1=value1 --p2=value2

  # Test with custom parameters
  redirector redirect test-slug --param key1=value1 --param key2=value2`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize service
		reg := registry.NewRegistry(config.GetConfig())
		service := reg.NewService()

		// Get slug from args
		slug := args[0]

		// Get flags
		requestID, _ := cmd.Flags().GetString("request-id")
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}

		userAgent, _ := cmd.Flags().GetString("ua")
		ipAddress, _ := cmd.Flags().GetString("ip")
		protocol, _ := cmd.Flags().GetString("protocol")
		urlStr, _ := cmd.Flags().GetString("url")
		referrer, _ := cmd.Flags().GetString("referrer")

		// Initialize params map
		params := make(map[string][]string)

		// Get p1-p4 parameters
		for i := 1; i <= 4; i++ {
			if val, _ := cmd.Flags().GetString(fmt.Sprintf("p%d", i)); val != "" {
				params[fmt.Sprintf("p%d", i)] = []string{val}
			}
		}

		// Get custom parameters
		customParams, _ := cmd.Flags().GetStringArray("param")
		for _, param := range customParams {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) == 2 {
				key, value := parts[0], parts[1]
				if existing, ok := params[key]; ok {
					params[key] = append(existing, value)
				} else {
					params[key] = []string{value}
				}
			}
		}

		// Parse URL if provided
		var incomeURL *url.URL
		if urlStr != "" {
			var err error
			incomeURL, err = url.Parse(urlStr)
			if err != nil {
				slog.Error("Invalid URL", "error", err)
				return
			}
		}

		// Create request data
		requestData := &dto.RedirectRequestData{
			RequestID: requestID,
			Slug:      slug,
			Params:    params,
			Headers:   make(map[string][]string),
			UserAgent: userAgent,
			IP:        net.ParseIP(ipAddress),
			Protocol:  protocol,
			URL:       incomeURL,
			Referer:   referrer,
		}

		// Validate request data
		if err := requestData.Validate(); err != nil {
			slog.Error("Invalid request data", "error", err)
			return
		}

		// Execute redirect
		result, err := service.Redirect(context.Background(), slug, requestData)
		if err != nil {
			slog.Error("Redirect failed", "error", err)
			return
		}

		// Print result
		fmt.Printf("Redirect Result:\n")
		fmt.Printf("  Target URL: %s\n", result.TargetURL)
		fmt.Printf("  Request ID: %s\n", requestData.RequestID)
		if len(params) > 0 {
			fmt.Printf("  Parameters:\n")
			for k, v := range params {
				fmt.Printf("    %s: %s\n", k, strings.Join(v, ", "))
			}
		}

		// Wait for click processing results
		for result := range result.OutputCh {
			if result.Err != nil {
				slog.Error("Click processing failed", "error", result.Err)
			} else {
				slog.Info("Click processed successfully")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(redirectCmd)

	// Add flags
	redirectCmd.Flags().String("request-id", "", "Custom request ID (optional, UUID v4 will be generated if not provided)")
	redirectCmd.Flags().String("ua", "", "User agent string")
	redirectCmd.Flags().String("ip", "127.0.0.1", "IP address (default: 127.0.0.1)")
	redirectCmd.Flags().String("protocol", "https", "Protocol (http/https)")
	redirectCmd.Flags().String("url", "", "Full URL of the request")
	redirectCmd.Flags().String("referrer", "", "Referrer URL")

	// Add p1-p4 parameter flags
	redirectCmd.Flags().String("p1", "", "Value for p1 parameter")
	redirectCmd.Flags().String("p2", "", "Value for p2 parameter")
	redirectCmd.Flags().String("p3", "", "Value for p3 parameter")
	redirectCmd.Flags().String("p4", "", "Value for p4 parameter")

	// Add custom parameter flag
	redirectCmd.Flags().StringArray("param", []string{}, "Custom parameters in key=value format (can be used multiple times)")
}
