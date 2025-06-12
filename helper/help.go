package helper

// CLI usage examples and help text
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

const Cloudrider_dump_help = `
Usage: aws-enumerator dump [options]

Options:
  -services string
        Services to dump
  -print bool
        Print results
  -filter string
        Filter API calls
  -errors bool
        Show errors

Examples:
  ./aws-enumerator dump -services all
  ./aws-enumerator dump -services iam -filter GetUser -print
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