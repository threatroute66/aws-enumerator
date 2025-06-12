package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Color functions for terminal output (keep original implementations)
func Red(text string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", text)
}

func Green(text string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", text)
}

func Yellow(text string) string {
	return fmt.Sprintf("\033[33m%s\033[0m", text)
}

func Reset() string {
	return "\033[0m"
}

// AWSCredentials represents AWS credential information (NEW)
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	Source          string // "profile", "env", or "env_file"
}

// CreateAWScredentialsFile creates the .env file (ENHANCED - backward compatibility)
func CreateAWScredentialsFile(region, accessKeyID, secretAccessKey, sessionToken *string) {
	file, err := os.Create(".env")
	if err != nil {
		fmt.Println(Red("Error: "), Yellow("Failed to create .env file"))
		return
	}
	defer file.Close()

	if region != nil && *region != "" {
		file.WriteString(fmt.Sprintf("AWS_REGION=%s\n", *region))
	}
	if accessKeyID != nil && *accessKeyID != "" {
		file.WriteString(fmt.Sprintf("AWS_ACCESS_KEY_ID=%s\n", *accessKeyID))
	}
	if secretAccessKey != nil && *secretAccessKey != "" {
		file.WriteString(fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s\n", *secretAccessKey))
	}
	if sessionToken != nil && *sessionToken != "" {
		file.WriteString(fmt.Sprintf("AWS_SESSION_TOKEN=%s\n", *sessionToken))
	}
}

// LoadCredentials loads AWS credentials with profile support (NEW)
func LoadCredentials(profile string) (*AWSCredentials, error) {
	// Priority order:
	// 1. AWS profile from ~/.aws/credentials (if profile specified)
	// 2. Environment variables
	// 3. .env file (backward compatibility)

	// If profile is specified, load from AWS credentials file
	if profile != "" {
		return loadFromProfile(profile)
	}

	// Check environment variables first
	if creds := loadFromEnvironment(); creds != nil {
		return creds, nil
	}

	// Fall back to .env file for backward compatibility
	return loadFromEnvFile()
}

// loadFromProfile loads credentials from AWS profile (NEW)
func loadFromProfile(profileName string) (*AWSCredentials, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}

	credFile := filepath.Join(homeDir, ".aws", "credentials")
	configFile := filepath.Join(homeDir, ".aws", "config")

	// Load credentials from profile
	creds, err := parseAWSCredentialsFile(credFile, profileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load profile %s: %v", profileName, err)
	}

	// Load region from config file if not in credentials
	if creds.Region == "" {
		if region := getRegionFromConfig(configFile, profileName); region != "" {
			creds.Region = region
		}
	}

	creds.Source = "profile"
	return creds, nil
}

// loadFromEnvironment loads credentials from environment variables (NEW)
func loadFromEnvironment() *AWSCredentials {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" {
		return nil
	}

	return &AWSCredentials{
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		SessionToken:    sessionToken,
		Region:          region,
		Source:          "env",
	}
}

// loadFromEnvFile loads credentials from .env file (ENHANCED - backward compatibility)
func loadFromEnvFile() (*AWSCredentials, error) {
	file, err := os.Open(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to open .env file: %v", err)
	}
	defer file.Close()

	creds := &AWSCredentials{Source: "env_file"}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "AWS_ACCESS_KEY_ID":
			creds.AccessKeyID = value
		case "AWS_SECRET_ACCESS_KEY":
			creds.SecretAccessKey = value
		case "AWS_SESSION_TOKEN":
			creds.SessionToken = value
		case "AWS_REGION":
			creds.Region = value
		}
	}

	if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
		return nil, fmt.Errorf("incomplete credentials in .env file")
	}

	return creds, nil
}

// parseAWSCredentialsFile parses AWS credentials file (NEW)
func parseAWSCredentialsFile(filename, profileName string) (*AWSCredentials, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var currentProfile string
	creds := &AWSCredentials{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for profile section
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentProfile = strings.Trim(line, "[]")
			continue
		}

		// Skip if not in target profile
		if currentProfile != profileName {
			continue
		}

		// Parse key-value pairs
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "aws_access_key_id":
			creds.AccessKeyID = value
		case "aws_secret_access_key":
			creds.SecretAccessKey = value
		case "aws_session_token":
			creds.SessionToken = value
		case "region":
			creds.Region = value
		}
	}

	if creds.AccessKeyID == "" || creds.SecretAccessKey == "" {
		return nil, fmt.Errorf("profile %s not found or incomplete", profileName)
	}

	return creds, nil
}

// getRegionFromConfig gets region from AWS config file (NEW)
func getRegionFromConfig(filename, profileName string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	var currentProfile string
	scanner := bufio.NewScanner(file)

	// For default profile, look for [default], for others look for [profile name]
	targetSection := "default"
	if profileName != "default" {
		targetSection = "profile " + profileName
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentProfile = strings.Trim(line, "[]")
			continue
		}

		if currentProfile != targetSection {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[0]) == "region" {
			return strings.TrimSpace(parts[1])
		}
	}

	return ""
}

// CreateAWSSession creates an AWS session with the loaded credentials (NEW)
func CreateAWSSession(creds *AWSCredentials) (*session.Session, error) {
	config := &aws.Config{}

	// Set region if available
	if creds.Region != "" {
		config.Region = aws.String(creds.Region)
	}

	// Set credentials based on source
	config.Credentials = credentials.NewStaticCredentials(
		creds.AccessKeyID,
		creds.SecretAccessKey,
		creds.SessionToken,
	)

	return session.NewSession(config)
}

// ListAvailableProfiles lists available AWS profiles (NEW)
func ListAvailableProfiles() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	credFile := filepath.Join(homeDir, ".aws", "credentials")
	file, err := os.Open(credFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profiles []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			profile := strings.Trim(line, "[]")
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}