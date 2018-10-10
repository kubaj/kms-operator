package stub

import (
	"context"
	"fmt"

	"github.com/kubaj/kms-operator/pkg/apis/kubaj/v1alpha1"
	cloudkms "google.golang.org/api/cloudkms/v1"

	"encoding/base64"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewHandler constructs Handler
func NewHandler(cloudKMS *cloudkms.Service) sdk.Handler {
	return &Handler{
		CloudKMS: cloudKMS,
	}
}

type Handler struct {
	CloudKMS *cloudkms.Service
}

// Handle is method for handling all watched events
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.SecretKMS:
		if event.Deleted {
			err := h.DeleteSecret(o)
			if err != nil {
				logrus.Warn("Failed to delete", err)
				return err
			}
		} else {
			err := h.CreateSecret(o)
			if err != nil {
				logrus.Warn("Failed to create", err)
				return err
			}
		}
	}
	return nil
}

// CreateSecret is method that creates v1/Secret according to SecretKMS resource
func (h *Handler) CreateSecret(cr *v1alpha1.SecretKMS) error {
	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Secret,
			Namespace: cr.Namespace,
		},
	}

	err := sdk.Get(secret)
	if err == nil {
		return nil
	}

	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	logrus.Debugf("Creating Secret from SecretKMS %s", cr.Name)

	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		cr.Spec.Provider.GoogleCloud.Project,
		cr.Spec.Provider.GoogleCloud.Location,
		cr.Spec.Provider.GoogleCloud.Keyring,
		cr.Spec.Provider.GoogleCloud.Key)

	req := &cloudkms.DecryptRequest{
		Ciphertext: cr.Spec.Provider.GoogleCloud.Data,
	}

	logrus.Debugln("Sending decrypt request")
	reqCall := h.CloudKMS.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parent, req)
	resp, err := reqCall.Do()
	if err != nil {
		return err
	}

	// Base64 decode after KMS call
	b, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return err
	}

	b, err = base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return err
	}

	secret.Data = make(map[string][]byte)
	secret.Data[cr.Spec.File] = b

	return sdk.Create(secret)
}

// DeleteSecret is method for handling Delete events of SecretKMS resource
func (h *Handler) DeleteSecret(cr *v1alpha1.SecretKMS) error {

	logrus.Debugf("Deleting Secret from SecretKMS %s", cr.Name)

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.Secret,
			Namespace: cr.Namespace,
		},
	}

	err := sdk.Delete(secret)
	if errors.IsNotFound(err) {
		return nil
	}

	return err
}
