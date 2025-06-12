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
	
	// TODO: Integrate with existing enumeration logic
	// The existing enumeration functions would be called here
	// with the AWS session containing the proper credentials
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

// DumpInfo placeholder for existing dump functionality
func DumpInfo(services *string, print *bool, filter *string, errors *bool) {
	// This is a placeholder for the existing DumpInfo function
	// The original implementation should be preserved here
	fmt.Printf("%s Dumping info for services: %s%s\n", 
		utils.Green("Info:"), *services, utils.Reset())
	
	// TODO: Add back the original dump logic
}