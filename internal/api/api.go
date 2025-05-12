package api

import (
	"AWS_Trafiic_Monitor/internal/model"
	"AWS_Trafiic_Monitor/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"log"
	"time"
)

var ErrAWSClientNotInitialized = errors.New("aws client not initialized")

var (
	DefaultSharedCredentialsPath = config.DefaultSharedCredentialsFilename()
	DefaultRegion                = "ap-northeast-1"
)

var lsClient *lightsail.Client

type Config struct {
	CredentialPath  string `json:"credentialPath,omitempty"`
	Region          string `json:"region,omitempty"`
	AccessKeyID     string `json:"accessKeyId,omitempty"`
	SecretAccessKey string `json:"secretAccessKey,omitempty"`
	SessionToken    string `json:"sessionToken,omitempty"`
}

// Init initializes the AWS SDK with the given credentials and region.
// sessionToken is optional, it could be empty if using long-term credentials
func Init(cfg *Config) error {
	if lsClient != nil {
		return nil
	}

	// Set default values if not provided
	if cfg.CredentialPath == "" {
		cfg.CredentialPath = DefaultSharedCredentialsPath
	}
	if cfg.Region == "" {
		cfg.Region = DefaultRegion
	}

	var awsConfig aws.Config
	// Set default credentials if not provided
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		// Configure options
		opts := []func(*config.LoadOptions) error{
			config.WithSharedConfigFiles([]string{}), // we don't want to use the default config file
			config.WithSharedCredentialsFiles([]string{cfg.CredentialPath}),
			config.WithDefaultRegion(cfg.Region),
		}

		// Load the Shared AWS Configuration
		var err error
		awsConfig, err = config.LoadDefaultConfig(context.Background(), opts...)
		if err != nil {
			return fmt.Errorf("unable to load aws config, err=%w", err)
		}
	} else {
		awsConfig = aws.Config{
			Region: DefaultRegion,
			Credentials: credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				cfg.SessionToken,
			),
		}
	}

	lsClient = lightsail.NewFromConfig(awsConfig)
	return nil
}

type API struct {
	InstanceName string
}

func NewAPI(name string) *API {
	return &API{name}
}

func (a *API) MetricDataMonth(name types.InstanceMetricName) (*model.MetricSum, error) {
	if lsClient == nil {
		return nil, ErrAWSClientNotInitialized
	}

	// Get time parameter for current month
	now := time.Now()
	start, end := utils.BeginningOfDay(now), utils.EndOfDay(now)
	period := end.Sub(start)

	// Prepare the input for GetInstanceMetricData
	input := &lightsail.GetInstanceMetricDataInput{
		InstanceName: aws.String(a.InstanceName),
		MetricName:   name,
		StartTime:    aws.Time(start),
		EndTime:      aws.Time(end),
		Unit:         types.MetricUnitBytes,
		Statistics:   []types.MetricStatistic{types.MetricStatisticSum},
		Period:       aws.Int32(multipleOfSixty(int32(period.Seconds()))),
	}
	output, err := lsClient.GetInstanceMetricData(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric for instance=%s metric=%s data: %w", a.InstanceName, name, err)
	}

	if len(output.MetricData) == 0 {
		return nil, fmt.Errorf("no metric data found for instance=%s metric=%s", a.InstanceName, name)
	}

	md := &model.MetricSum{
		Name: string(name),
		Data: []model.MetricSumData{
			{Sum: *output.MetricData[0].Sum},
		},
	}

	return md, nil
}

func (a *API) Shutdown(force bool) error {
	log.Printf("Shutting down instance=%s", a.InstanceName)

	input := &lightsail.StopInstanceInput{
		InstanceName: aws.String(a.InstanceName),
		Force:        aws.Bool(force),
	}

	_, err := lsClient.StopInstance(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to stop instance=%s: %w", a.InstanceName, err)
	}
	return nil
}

func multipleOfSixty(seconds int32) int32 {
	if seconds <= 0 {
		return 60 // Minimum valid period
	}

	// If already a multiple of 60, return as is
	if seconds%60 == 0 {
		return seconds
	}

	// Round to the nearest multiple of 60
	remainder := seconds % 60

	if remainder >= 30 {
		// Round up
		return seconds + (60 - remainder)
	} else {
		// Round down
		return seconds - remainder
	}
}
