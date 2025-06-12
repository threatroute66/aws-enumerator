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