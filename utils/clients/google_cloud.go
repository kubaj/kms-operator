package clients

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
	"io/ioutil"
)

// GetGoogleCloudKMS returns initialized Google Cloud KMS client based on provided flags
func GetGoogleCloudKMS(enabled bool, serviceAccount string) (*cloudkms.Service, error) {
	if !enabled {
		return nil, nil
	}

	ctx := context.Background()

	// Use default service account
	if serviceAccount == "" {
		client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
		if err != nil {
			return nil, err
		}

		return cloudkms.New(client)
	}

	// Parse service account provided by path
	b, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return nil, err
	}

	credentials, err := google.CredentialsFromJSON(ctx, b, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	return cloudkms.New(oauth2.NewClient(ctx, credentials.TokenSource))
}
