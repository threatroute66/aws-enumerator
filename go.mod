module github.com/threatroute66/aws-enumerator

go 1.22

toolchain go1.24.2

require (
	github.com/aws/aws-sdk-go v1.44.0
	github.com/aws/aws-sdk-go-v2/config v1.29.16
	github.com/aws/aws-sdk-go-v2/service/acm v1.32.2
	github.com/aws/aws-sdk-go-v2/service/amplify v1.33.2
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.31.2
	github.com/aws/aws-sdk-go-v2/service/appmesh v1.30.3
	github.com/aws/aws-sdk-go-v2/service/appsync v1.47.1
	github.com/aws/aws-sdk-go-v2/service/athena v1.51.1
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.53.2
	github.com/aws/aws-sdk-go-v2/service/backup v1.42.2
	github.com/aws/aws-sdk-go-v2/service/batch v1.52.5
	github.com/aws/aws-sdk-go-v2/service/chime v1.36.3
	github.com/aws/aws-sdk-go-v2/service/cloud9 v1.29.3
	github.com/aws/aws-sdk-go-v2/service/clouddirectory v1.25.3
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.60.2
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.46.2
	github.com/aws/aws-sdk-go-v2/service/cloudhsm v1.25.3
	github.com/aws/aws-sdk-go-v2/service/cloudhsmv2 v1.30.3
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v1.27.3
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.49.2
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.61.1
	github.com/aws/aws-sdk-go-v2/service/codecommit v1.28.3
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.30.5
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.42.1
	github.com/aws/aws-sdk-go-v2/service/codestar v1.23.4
	github.com/aws/aws-sdk-go-v2/service/comprehend v1.36.5
	github.com/aws/aws-sdk-go-v2/service/datapipeline v1.26.3
	github.com/aws/aws-sdk-go-v2/service/datasync v1.49.2
	github.com/aws/aws-sdk-go-v2/service/dax v1.24.3
	github.com/aws/aws-sdk-go-v2/service/devicefarm v1.31.1
	github.com/aws/aws-sdk-go-v2/service/directconnect v1.32.4
	github.com/aws/aws-sdk-go-v2/service/dlm v1.30.6
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.43.3
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.225.1
	github.com/aws/aws-sdk-go-v2/service/ecr v1.44.2
	github.com/aws/aws-sdk-go-v2/service/ecs v1.57.4
	github.com/aws/aws-sdk-go-v2/service/eks v1.66.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.46.2
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.29.3
	github.com/aws/aws-sdk-go-v2/service/elastictranscoder v1.28.3
	github.com/aws/aws-sdk-go-v2/service/firehose v1.37.6
	github.com/aws/aws-sdk-go-v2/service/fms v1.40.4
	github.com/aws/aws-sdk-go-v2/service/fsx v1.54.1
	github.com/aws/aws-sdk-go-v2/service/gamelift v1.41.2
	github.com/aws/aws-sdk-go-v2/service/globalaccelerator v1.30.3
	github.com/aws/aws-sdk-go-v2/service/glue v1.113.2
	github.com/aws/aws-sdk-go-v2/service/greengrass v1.28.3
	github.com/aws/aws-sdk-go-v2/service/guardduty v1.54.6
	github.com/aws/aws-sdk-go-v2/service/health v1.30.3
	github.com/aws/aws-sdk-go-v2/service/iam v1.42.1
	github.com/aws/aws-sdk-go-v2/service/inspector v1.26.3
	github.com/aws/aws-sdk-go-v2/service/iot v1.64.3
	github.com/aws/aws-sdk-go-v2/service/iotanalytics v1.27.3
	github.com/aws/aws-sdk-go-v2/service/kafka v1.39.4
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.35.2
	github.com/aws/aws-sdk-go-v2/service/kinesisanalytics v1.26.5
	github.com/aws/aws-sdk-go-v2/service/kinesisvideo v1.28.3
	github.com/aws/aws-sdk-go-v2/service/kms v1.40.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.71.4
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.43.3
	github.com/aws/aws-sdk-go-v2/service/machinelearning v1.29.3
	github.com/aws/aws-sdk-go-v2/service/macie v1.19.2
	github.com/aws/aws-sdk-go-v2/service/mediaconnect v1.40.1
	github.com/aws/aws-sdk-go-v2/service/mediaconvert v1.74.1
	github.com/aws/aws-sdk-go-v2/service/medialive v1.76.1
	github.com/aws/aws-sdk-go-v2/service/mediapackage v1.35.3
	github.com/aws/aws-sdk-go-v2/service/mediastore v1.25.3
	github.com/aws/aws-sdk-go-v2/service/mediatailor v1.48.1
	github.com/aws/aws-sdk-go-v2/service/mobile v1.21.3
	github.com/aws/aws-sdk-go-v2/service/mq v1.29.1
	github.com/aws/aws-sdk-go-v2/service/opsworks v1.27.4
	github.com/aws/aws-sdk-go-v2/service/organizations v1.38.4
	github.com/aws/aws-sdk-go-v2/service/pinpoint v1.35.3
	github.com/aws/aws-sdk-go-v2/service/polly v1.48.3
	github.com/aws/aws-sdk-go-v2/service/pricing v1.34.4
	github.com/aws/aws-sdk-go-v2/service/ram v1.30.5
	github.com/aws/aws-sdk-go-v2/service/rds v1.97.2
	github.com/aws/aws-sdk-go-v2/service/redshift v1.54.5
	github.com/aws/aws-sdk-go-v2/service/rekognition v1.47.1
	github.com/aws/aws-sdk-go-v2/service/robomaker v1.31.3
	github.com/aws/aws-sdk-go-v2/service/route53 v1.52.1
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.29.3
	github.com/aws/aws-sdk-go-v2/service/route53resolver v1.35.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.80.2
	github.com/aws/aws-sdk-go-v2/service/sagemaker v1.195.1
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.35.6
	github.com/aws/aws-sdk-go-v2/service/securityhub v1.57.5
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v1.34.1
	github.com/aws/aws-sdk-go-v2/service/shield v1.30.3
	github.com/aws/aws-sdk-go-v2/service/signer v1.27.3
	github.com/aws/aws-sdk-go-v2/service/sms v1.25.3
	github.com/aws/aws-sdk-go-v2/service/snowball v1.31.4
	github.com/aws/aws-sdk-go-v2/service/sns v1.34.6
	github.com/aws/aws-sdk-go-v2/service/sqs v1.38.7
	github.com/aws/aws-sdk-go-v2/service/ssm v1.59.2
	github.com/aws/aws-sdk-go-v2/service/storagegateway v1.37.2
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.21
	github.com/aws/aws-sdk-go-v2/service/support v1.27.3
	github.com/aws/aws-sdk-go-v2/service/transcribe v1.46.1
	github.com/aws/aws-sdk-go-v2/service/transfer v1.60.2
	github.com/aws/aws-sdk-go-v2/service/translate v1.29.3
	github.com/aws/aws-sdk-go-v2/service/waf v1.26.3
	github.com/aws/aws-sdk-go-v2/service/workdocs v1.26.3
	github.com/aws/aws-sdk-go-v2/service/worklink v1.23.2
	github.com/aws/aws-sdk-go-v2/service/workmail v1.31.3
	github.com/aws/aws-sdk-go-v2/service/workspaces v1.57.1
	github.com/aws/aws-sdk-go-v2/service/xray v1.31.6
	github.com/fatih/color v1.18.0
	github.com/joho/godotenv v1.5.1
	github.com/wayneashleyberry/terminal-dimensions v1.1.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.36.4 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.69 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.35 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.35 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.2 // indirect
	github.com/aws/smithy-go v1.22.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
