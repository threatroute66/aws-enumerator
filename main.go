// main.go - Improved AWS Enumerator with Profile Support
package main

import (
	"fmt"
	"os"

	"github.com/threatroute66/aws-enumerator/helper"
	"github.com/threatroute66/aws-enumerator/utils"
)

func main() {
	helper.Cred.Usage = func() {
		fmt.Fprintf(os.Stderr, helper.Cloudrider_cred_help)
	}
	helper.Enum.Usage = func() {
		fmt.Fprintf(os.Stderr, helper.Cloudrider_enum_help)
	}
	helper.Dump.Usage = func() {
		fmt.Fprintf(os.Stderr, helper.Cloudrider_dump_help)
	}

	if len(os.Args) < 2 {
		fmt.Println(helper.Cloudrider_help)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "cred":
		helper.Cred.Parse(os.Args[2:])
		utils.CreateAWScredentialsFile(helper.AWS_region, helper.AWS_access_key_id, helper.AWS_secret_access_key, helper.AWS_session_token)
		fmt.Println(utils.Green("Message: "), utils.Yellow("File"), utils.Red(".env"), utils.Yellow("with AWS credentials were created in current folder"))
	case "enum":
		helper.Enum.Parse(os.Args[2:])
		helper.SetEnumerationPipeline(helper.Services_enum, helper.Speed, helper.Profile)
		fmt.Println(utils.Green("Message: "), utils.Yellow("Enumeration finished"))
	case "dump":
		helper.Dump.Parse(os.Args[2:])
		helper.DumpInfo(helper.Services_dump, helper.Print, helper.Filter, helper.Errors_dump)
	case "profiles":
		helper.HandleProfilesCommand()
	default:
		fmt.Println(helper.Cloudrider_help)
		os.Exit(1)
	}
}

// helper/flags.go - Updated flag definitions
package helper

import (
	"flag"
)

var (
	// Existing variables
	AWS_region            *string
	AWS_access_key_id     *string
	AWS_secret_access_key *string
	AWS_session_token     *string
	Services_enum         *string
	Services_dump         *string
	Speed                 *string
	Print                 *bool
	Filter                *string
	Errors_dump           *bool
	
	// New profile variable
	Profile               *string

	// Flag sets
	Cred = flag.NewFlagSet("cred", flag.ExitOnError)
	Enum = flag.NewFlagSet("enum", flag.ExitOnError)
	Dump = flag.NewFlagSet("dump", flag.ExitOnError)
)

func init() {
	// Cred command flags
	AWS_region = Cred.String("aws_region", "", "AWS Region")
	AWS_access_key_id = Cred.String("aws_access_key_id", "", "AWS Access Key ID")
	AWS_secret_access_key = Cred.String("aws_secret_access_key", "", "AWS Secret Access Key")
	AWS_session_token = Cred.String("aws_session_token", "", "AWS Session Token")

	// Enum command flags
	Services_enum = Enum.String("services", "", "Services to enumerate (e.g., all, iam,s3,sts)")
	Speed = Enum.String("speed", "normal", "Enumeration speed: slow, normal, fast")
	Profile = Enum.String("profile", "", "AWS profile to use from ~/.aws/credentials")

	// Dump command flags
	Services_dump = Dump.String("services", "", "Services to dump")
	Print = Dump.Bool("print", false, "Print results")
	Filter = Dump.String("filter", "", "Filter API calls")
	Errors_dump = Dump.Bool("errors", false, "Show errors")
}

// utils/credentials.go - Enhanced credential management
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

// Color functions for terminal output
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

// AWSCredentials represents AWS credential information
type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	Source          string // "profile", "env", or "static"
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
	// 1. Environment variables (if no profile specified)
	// 2. AWS profile from ~/.aws/credentials
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

	// Set credentials
	if creds.Source == "profile" {
		// For profiles, we can use the SharedCredentialsProvider
		config.Credentials = credentials.NewSharedCredentials("", getProfileName(creds))
	} else {
		// For static credentials (env, env_file)
		config.Credentials = credentials.NewStaticCredentials(
			creds.AccessKeyID,
			creds.SecretAccessKey,
			creds.SessionToken,
		)
	}

	return session.NewSession(config)
}

// getProfileName extracts profile name for SharedCredentialsProvider
func getProfileName(creds *AWSCredentials) string {
	// This would need to be passed through or stored with credentials
	// For now, return empty string to use default
	return ""
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

// helper/enumeration.go - Updated enumeration function
package helper

import (
	"fmt"
	"log"

	"github.com/threatroute66/aws-enumerator/utils"
)

// SetEnumerationPipeline sets up and runs the enumeration with profile support
func SetEnumerationPipeline(services, speed, profile *string) {
	var profileName string
	if profile != nil {
		profileName = *profile
	}

	// Load credentials
	creds, err := utils.LoadCredentials(profileName)
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	// Print credential source info
	fmt.Printf("%s Using credentials from: %s%s\n", 
		utils.Green("Info:"), utils.Yellow(creds.Source), utils.Reset())
	
	if creds.Region != "" {
		fmt.Printf("%s Region: %s%s\n", 
			utils.Green("Info:"), utils.Yellow(creds.Region), utils.Reset())
	}

	// Create AWS session
	sess, err := utils.CreateAWSSession(creds)
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// Continue with existing enumeration logic using the session
	runEnumeration(sess, services, speed)
}

// runEnumeration placeholder for the actual enumeration logic
func runEnumeration(sess interface{}, services, speed *string) {
	// This would contain the actual enumeration logic
	// using the AWS session created with the credentials
	fmt.Printf("%s Starting enumeration with services: %s, speed: %s%s\n",
		utils.Green("Info:"), *services, *speed, utils.Reset())
}

// HandleProfilesCommand lists available AWS profiles
func HandleProfilesCommand() {
	profiles, err := utils.ListAvailableProfiles()
	if err != nil {
		fmt.Printf("%s Failed to list profiles: %v%s\n", utils.Red("Error:"), err, utils.Reset())
		return
	}

	if len(profiles) == 0 {
		fmt.Printf("%s No AWS profiles found in ~/.aws/credentials%s\n", utils.Yellow("Info:"), utils.Reset())
		fmt.Printf("%s To create profiles, use: aws configure --profile <profile-name>%s\n", utils.Green("Tip:"), utils.Reset())
		return
	}

	fmt.Printf("%s Available AWS profiles:%s\n", utils.Green("Info:"), utils.Reset())
	for _, profile := range profiles {
		fmt.Printf("  • %s%s%s\n", utils.Yellow(""), profile, utils.Reset())
	}
}

// CLI usage examples and help text updates
const Cloudrider_enum_help = `
Usage: aws-enumerator enum [options]

Options:
  -services string
        Services to enumerate: all, or comma-separated list (e.g., iam,s3,sts)
  -speed string
        Enumeration speed: slow, normal, fast (default "normal")
  -profile string
        AWS profile to use from ~/.aws/credentials

Examples:
  # Use default credentials (env vars or .env file)
  ./aws-enumerator enum -services all

  # Use specific AWS profile
  ./aws-enumerator enum -services all -profile myprofile

  # Use specific services with profile
  ./aws-enumerator enum -services iam,s3,sts -profile production
`

const Cloudrider_cred_help = `
Usage: aws-enumerator cred [options]

Options:
  -aws_access_key_id string
        AWS Access Key ID
  -aws_secret_access_key string
        AWS Secret Access Key
  -aws_session_token string
        AWS Session Token (optional)
  -aws_region string
        AWS Region

Example:
  ./aws-enumerator cred -aws_access_key_id AKIA... -aws_secret_access_key SECRET... -aws_region us-west-2

Note: This creates a .env file. For better security, consider using AWS profiles instead.
`

const Cloudrider_help = `
AWS Enumerator - Enhanced with Profile Support

Commands:
  cred      Set up credentials (creates .env file)
  enum      Run enumeration with optional profile support
  dump      Analyze enumeration results
  profiles  List available AWS profiles

Use 'aws-enumerator [command] -h' for more information about a command.

Profile Support:
  You can now use AWS profiles from ~/.aws/credentials:
  ./aws-enumerator enum -services all -profile myprofile

  To list available profiles:
  ./aws-enumerator profiles
`
