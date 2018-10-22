package clients

import (
	"context"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"google.golang.org/api/option"
)

// GetGoogleCloudKMS returns initialized Google Cloud KMS client based on provided flags
func GetGoogleCloudKMS(enabled bool, serviceAccount string) (*cloudkms.KeyManagementClient, error) {
	if !enabled {
		return nil, nil
	}

	ctx := context.Background()

	// Use default service account
	if serviceAccount == "" {
		return cloudkms.NewKeyManagementClient(ctx)
	}

	return cloudkms.NewKeyManagementClient(ctx, option.WithCredentialsFile(serviceAccount))
}
