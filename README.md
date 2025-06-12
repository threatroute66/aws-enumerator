# AWS Service Enumeration

Forked from Pavel Shabarkin's https://github.com/shabarkin/aws-enumerator

# Disclaimer

The tool is in beta stage (testing in progress), no destructive API Calls used ( read only actions ).
I hope, there will be no issues with the tool. If any issues encountered, please submit the ticket. 

# Description

The AWS Enumerator was created for service enumeration and info dumping for investigations of penetration testers during Black-Box  testing. The tool is intended to speed up the process of Cloud review in case the security researcher compromised AWS Account Credentials. 

AWS Enumerator supports more than 600 API Calls ( reading actions `Get`,  `List`, `Describe` etc... ), and will be extended. 

The tool provides interface for result analysis. All results are saved in json files (one time "Database").

# Install

If you have Go installed and configured (i.e. with `$GOPATH/bin` in your `$PATH`):

```
go get -u github.com/threatroute66/aws-enumerator
```
```
go install -v github.com/threatroute66/aws-enumerator@latest
```

# Basic Usage

## Credentials setup

To setup credentials, you should use cred subcommand and supply credentials: 

```bash
./aws-enumerator cred -aws_access_key_id AKIA***********XKU -aws_region us-west-2 -aws_secret_access_key kIm6m********************5JPF
```

![_img/Screenshot_2021-04-10_at_14.43.51.png](_img/Screenshot_2021-04-10_at_14.43.51.png)

![_img/Screenshot_2021-04-10_at_14.45.51.png](_img/Screenshot_2021-04-10_at_14.45.51.png)

It creates `.env` file, which is loaded to global variables each time you call `enum` subcommand.

**WARNING:** If you set these values `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN` in global variables manually before running the tool, it will not be able to load AWS Credentials specified in `.env` file ( It can't overwrite global variables ).

## Enumeration

To enumerate all services, you should use enum subcommand and supply all value or iam,s3,sts,rds ( no spaces between commas ), etc. ...

```bash
./aws-enumerator enum -services all
```

 If you want to check specific services (up to 10 ):

```bash
./aws-enumerator enum -services iam,sts,rds
```

![_img/Screenshot_2021-04-10_at_13.36.56.png](_img/Screenshot_2021-04-10_at_13.36.56.png)

(`-speed` flag is optional, the default value is `normal` ) There are 3 options `slow`, `normal`, `fast` 

```bash
./aws-enumerator enum -services all -speed slow
```

## Analysis

To analyse the collected information, you should use `dump` subcommand: ( Use `all` for quick overview of available API calls )

```bash
./aws-enumerator dump -services all
```

![_img/Screenshot_2021-04-10_at_13.56.12.png](_img/Screenshot_2021-04-10_at_13.56.12.png)

Analyze specific services (up to 10) `iam,s3,sts`, etc ...

```bash
./aws-enumerator dump -services iam,s3,sts
```

![_img/Screenshot_2021-04-10_at_14.03.16.png](_img/Screenshot_2021-04-10_at_14.03.16.png)

To filter API calls, you should use `-filter` option, start typing the name of API call (`GetA` ...):

```bash
./aws-enumerator dump -services iam -filter GetA
```

![_img/Screenshot_2021-04-10_at_14.06.18.png](_img/Screenshot_2021-04-10_at_14.06.18.png)

To retrieve the result of API call, you should use `-print` option

```bash
./aws-enumerator dump -services iam -filter ListS -print
```

![_img/Screenshot_2021-04-10_at_14.08.01.png](_img/Screenshot_2021-04-10_at_14.08.01.png)

## Demo Video

[Pavel Shabarkin LinkedIn](https://www.linkedin.com/posts/pavelshabarkin_cybersecurity-hacking-awssecurity-activity-6785479892881416192-O29U/)

# AWS Enumerator Profile Enhancement Summary

## Overview

This is an enhanced version of the original AWS Enumerator code that adds AWS profile support while maintaining full backward compatibility. Here's a comprehensive summary of the improvements:

## Key Improvements

### 1. **AWS Profile Support** 
- **Added profile flag**: New `-profile` option for the `enum` command
- **Credential file parsing**: Full support for `~/.aws/credentials` format
- **Config file integration**: Automatic region detection from `~/.aws/config`
- **Named profiles**: Support for multiple profiles (dev, staging, production, etc.)

### 2. **Enhanced Credential Management** 
- **Priority-based loading**: Environment variables → Profiles → .env file
- **Credential source tracking**: Clear indication of which source is being used
- **Session token support**: Full support for temporary credentials
- **Error handling**: Comprehensive error messages and validation

### 3. **New Commands** 
- **`profiles` command**: Lists all available AWS profiles
- **Profile validation**: Checks for credential completeness
- **User guidance**: Helpful messages for profile setup

### 4. **Security Enhancements** 🛡
- **File permissions**: Guidance on proper credential file security
- **Temporary credentials**: Better support for session tokens
- **Separation of concerns**: Cleaner credential management architecture

## Code Structure Changes

### New Files Added:
1. **Enhanced `main.go`**: Added profile command handling
2. **Updated `helper/flags.go`**: Added profile flag definition
3. **New `utils/credentials.go`**: Comprehensive credential management
4. **Enhanced `helper/enumeration.go`**: Profile-aware enumeration

### Key Functions Added:
- `LoadCredentials(profile string)`: Multi-source credential loading
- `parseAWSCredentialsFile()`: AWS credentials file parser
- `getRegionFromConfig()`: Region extraction from config files
- `ListAvailableProfiles()`: Profile discovery
- `CreateAWSSession()`: Session creation with credentials
- `HandleProfilesCommand()`: Profile listing command

## Usage Examples

### Before (Original):
```bash
# Only .env file support
./aws-enumerator cred -aws_access_key_id AKIA... -aws_secret_access_key ...
./aws-enumerator enum -services all
```

### After (Enhanced):
```bash
# Multiple credential sources
./aws-enumerator profiles                              # List profiles
./aws-enumerator enum -services all -profile dev       # Use profile
./aws-enumerator enum -services all                    # Use env vars or .env
./aws-enumerator cred -aws_access_key_id AKIA...       # Still works
```

## Backward Compatibility

**Fully maintained**:
- All existing commands work unchanged
- `.env` file approach still supported
- Original flags and options preserved
- No breaking changes to existing workflows

## Technical Implementation Details

### 1. **Credential Loading Logic**
```go
func LoadCredentials(profile string) (*AWSCredentials, error) {
    if profile != "" {
        return loadFromProfile(profile)
    }
    if creds := loadFromEnvironment(); creds != nil {
        return creds, nil
    }
    return loadFromEnvFile()
}
```

### 2. **AWS File Format Support**
- **Credentials file**: Standard INI format parsing
- **Config file**: Profile-specific region extraction
- **Section handling**: Support for `[default]` and `[profile name]` formats

### 3. **Session Management**
```go
func CreateAWSSession(creds *AWSCredentials) (*session.Session, error) {
    config := &aws.Config{}
    if creds.Source == "profile" {
        config.Credentials = credentials.NewSharedCredentials("", profileName)
    } else {
        config.Credentials = credentials.NewStaticCredentials(...)
    }
    return session.NewSession(config)
}
```

## Benefits for Users

### 1. **Security**
- No need to store credentials in project directories
- Leverages AWS CLI credential management
- Supports IAM role assumption and temporary credentials

### 2. **Convenience**
- Multiple environment support (dev/staging/prod)
- Quick profile switching
- Standard AWS tooling integration

### 3. **Professional Workflows**
- Aligns with AWS best practices
- Enterprise-ready credential management
- Audit-friendly credential handling

## Migration Path

### For existing users:
1. **Continue using current workflow** (no changes needed)
2. **Gradually migrate to profiles** for better security
3. **Use new profile features** for multi-environment testing

### For new users:
1. **Set up AWS profiles** using `aws configure`
2. **Use profile-based enumeration** from day one
3. **Leverage modern AWS credential practices**

## Testing Checklist

**Backward compatibility**:
- `.env` file approach works
- All original commands function
- Environment variable support maintained

**New functionality**:
- Profile listing works
- Profile-based enumeration functions
- Credential source indication correct
- Region detection from config files

**Error handling**:
- Missing profile errors
- Invalid credential errors
- File permission issues
- Helpful user guidance

## Integration Notes

The enhanced version integrates seamlessly with:
- **AWS CLI**: Uses same credential files
- **AWS SDK**: Standard session management
- **IAM roles**: Supports role-based access
- **CI/CD pipelines**: Environment variable support
- **Docker containers**: Flexible credential mounting


