package helper

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/threatroute66/aws-enumerator/utils"
	"github.com/threatroute66/aws-enumerator/servicemaster"
	"github.com/threatroute66/aws-enumerator/servicestructs"
)

// SetEnumerationPipeline sets up credentials and runs servicemaster enumeration
func SetEnumerationPipeline(services, speed, profile *string) {
	var profileName string
	if profile != nil {
		profileName = *profile
	}

	// Load credentials using new credential management
	if profileName != "" {
		creds, err := utils.LoadCredentials(profileName)
		if err != nil {
			log.Fatalf("Failed to load credentials: %v", err)
		}

		// Set environment variables for servicemaster
		os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID)
		os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey)
		if creds.SessionToken != "" {
			os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken)
		}
		if creds.Region != "" {
			os.Setenv("AWS_REGION", creds.Region)
		}

		fmt.Printf("%s Using credentials from: profile%s\n", 
			utils.Green("Info:"), utils.Reset())
		if creds.Region != "" {
			fmt.Printf("%s Region: %s%s\n", 
				utils.Green("Info:"), utils.Yellow(creds.Region), utils.Reset())
		}
	} else {
		// Load from env/env_file for backward compatibility
		creds, err := utils.LoadCredentials("")
		if err == nil && creds != nil {
			fmt.Printf("%s Using credentials from: %s%s\n", 
				utils.Green("Info:"), utils.Yellow(creds.Source), utils.Reset())
			if creds.Region != "" {
				fmt.Printf("%s Region: %s%s\n", 
					utils.Green("Info:"), utils.Yellow(creds.Region), utils.Reset())
			}
		}
	}

	// Check credentials are available
	if !servicemaster.CheckAWSCredentials() {
		log.Fatalf("AWS credentials not found or invalid")
	}

	// Get all AWS services from servicestructs
	allServices := servicestructs.GetServices()

	// Parse services - convert "all" or "iam,s3,sts" to string slice
	var wantedServices []string
	if *services == "all" {
		wantedServices = []string{"all"}
	} else {
		wantedServices = strings.Split(*services, ",")
		// Trim whitespace from each service
		for i := range wantedServices {
			wantedServices[i] = strings.TrimSpace(wantedServices[i])
		}
	}

	// Convert speed string to int
	speedInt := convertSpeedToInt(*speed)

	fmt.Printf("%s Starting enumeration with services: %s, speed: %s%s\n",
		utils.Green("Info:"), *services, *speed, utils.Reset())

	// Call the actual servicemaster enumeration - THIS IS THE KEY LINE
	servicemaster.ServiceCall(allServices, wantedServices, speedInt)
}

// convertSpeedToInt converts speed string to int (based on original logic)
func convertSpeedToInt(speed string) int {
	switch speed {
	case "slow":
		return 1
	case "normal":
		return 2
	case "fast":
		return 3
	default:
		return 2 // default to normal
	}
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

// DumpInfo - preserve original dump functionality
func DumpInfo(services *string, print *bool, filter *string, errors *bool) {
	// TODO: Copy the original DumpInfo implementation here
	fmt.Printf("%s Dumping info for services: %s%s\n", 
		utils.Green("Info:"), *services, utils.Reset())
}