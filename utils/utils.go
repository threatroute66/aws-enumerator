package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Global variables for file paths
var (
	FILEPATH       = "enum-results/"
	ERROR_FILEPATH = "enum-results/errors/"
)

// Color functions for terminal output
func Red(text interface{}) string {
	return fmt.Sprintf("\033[31m%v\033[0m", text)
}

func Green(text interface{}) string {
	return fmt.Sprintf("\033[32m%v\033[0m", text)
}

func Yellow(text interface{}) string {
	return fmt.Sprintf("\033[33m%v\033[0m", text)
}

func Reset() string {
	return "\033[0m"
}

// AWSCredentials represents AWS credential information
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	Source          string // "profile", "env", or "env_file"
}

// PackResponse packs response data for JSON output (returns string for servicemaster)
func PackResponse(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return "{}"  // Return empty JSON string on error
	}
	return string(jsonData)
}

// CreateAWScredentialsFile creates the .env file (backward compatibility)
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

// LoadCredentials loads AWS credentials with profile support
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

// loadFromProfile loads credentials from AWS profile
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

// loadFromEnvironment loads credentials from environment variables
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

// loadFromEnvFile loads credentials from .env file (backward compatibility)
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

// parseAWSCredentialsFile parses AWS credentials file
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

// getRegionFromConfig gets region from AWS config file
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

// CreateAWSSession creates an AWS session with the loaded credentials
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

// ListAvailableProfiles lists available AWS profiles
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

// EnsureDirectories creates necessary directories for output files
func EnsureDirectories() error {
	// Create main results directory
	if err := os.MkdirAll(FILEPATH, 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %v", err)
	}

	// Create errors directory
	if err := os.MkdirAll(ERROR_FILEPATH, 0755); err != nil {
		return fmt.Errorf("failed to create errors directory: %v", err)
	}

	return nil
}

// WriteJSONFile writes data to a JSON file
func WriteJSONFile(filename string, data interface{}) error {
	jsonString := PackResponse(data)
	jsonData := []byte(jsonString)

	fullPath := filepath.Join(FILEPATH, filename)
	if err := ioutil.WriteFile(fullPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", fullPath, err)
	}

	return nil
}

// WriteErrorFile writes error data to a file
func WriteErrorFile(filename string, errorData interface{}) error {
	jsonString := PackResponse(errorData)
	jsonData := []byte(jsonString)

	fullPath := filepath.Join(ERROR_FILEPATH, filename)
	if err := ioutil.WriteFile(fullPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write error file %s: %v", fullPath, err)
	}

	return nil
}

// StringToInt converts string to integer
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// IntToString converts integer to string
func IntToString(i int) string {
	return strconv.Itoa(i)
}

// CheckEnvFileExistance checks if .env file exists  
func CheckEnvFileExistance() bool {
	_, err := os.Stat(".env")
	return !os.IsNotExist(err)
}

// Find searches for a value in a slice
func Find(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}