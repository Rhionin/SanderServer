package main

import (
	"github.com/Rhionin/SanderServer/config"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	StackName          = "StormWatch-ProgressCheck"
	FuncionName        = "ProgressCheck"
	MemorySizeMB       = 128
	MaxDurationSeconds = 20
	CodePath           = "./getProgressLambda"
	Handler            = "bootstrap"
	SecretName         = "StormlightArchive"
)

type StormWatchCdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *StormWatchCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Define the Lambda function resource
	progressCheckFunction := awslambda.NewFunction(stack, jsii.String("ProgressCheck"), &awslambda.FunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Architecture: awslambda.Architecture_ARM_64(),
		MemorySize:   jsii.Number(MemorySizeMB),
		Timeout:      awscdk.Duration_Seconds(jsii.Number(MaxDurationSeconds)),
		Code:         awslambda.AssetCode_FromAsset(jsii.String(CodePath), nil),
		LogRetention: awslogs.RetentionDays_ONE_DAY,
		Handler:      jsii.String(Handler),
	})

	// Define the Lambda function URL resource
	progressCheckFunctionUrl := progressCheckFunction.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})

	// Define a CloudFormation output for your URL
	awscdk.NewCfnOutput(stack, jsii.String("progressCheckFunctionUrlOutput"), &awscdk.CfnOutputProps{
		Value: progressCheckFunctionUrl.Url(),
	})

	secret := awssecretsmanager.Secret_FromSecretNameV2(stack, jsii.String(SecretName+"SecretID"), jsii.String(SecretName))
	secret.GrantRead(progressCheckFunction, nil)

	awsevents.NewRule(stack, jsii.String("storm-check"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Expression(jsii.String("rate(5 minutes)")),
		Targets:  &[]awsevents.IRuleTarget{awseventstargets.NewLambdaFunction(progressCheckFunction, nil)},
	})

	history := awsdynamodb.NewTableV2(stack, jsii.String(config.HistoryDynamoTableName), &awsdynamodb.TablePropsV2{
		TableName: jsii.String(config.HistoryDynamoTableName),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("ID"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		SortKey: &awsdynamodb.Attribute{
			Name: jsii.String("TimestampUnixNano"),
			Type: awsdynamodb.AttributeType_NUMBER,
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})
	progressCheckFunction.Role().AttachInlinePolicy(awsiam.NewPolicy(stack, jsii.String("stormwatch-dynamo"), &awsiam.PolicyProps{
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: jsii.Strings(
					"dynamodb:Query",
					"dynamodb:PutItem",
				),
				Resources: jsii.Strings(*history.TableArn()),
			}),
		},
	}))

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "StormWatchStack", &StormWatchCdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
