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

func NewHandler(cloudKMS *cloudkms.Service) sdk.Handler {
	return &Handler{
		CloudKMS: cloudKMS,
	}
}

type Handler struct {
	CloudKMS *cloudkms.Service
}

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

	parent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		cr.Spec.Provider.GoogleCloud.Project,
		cr.Spec.Provider.GoogleCloud.Location,
		cr.Spec.Provider.GoogleCloud.Keyring,
		cr.Spec.Provider.GoogleCloud.Key)

	req := &cloudkms.DecryptRequest{
		Ciphertext: cr.Spec.Provider.GoogleCloud.Data,
	}

	logrus.Debugln("Sending decrypt request")
	resp, err := h.CloudKMS.Projects.Locations.KeyRings.CryptoKeys.Decrypt(parent, req).Do()
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

	logrus.Infoln("Decrypted secret", string(b))

	secret.Data = make(map[string][]byte)
	secret.Data[cr.Spec.File] = b

	return sdk.Create(secret)
}

func (h *Handler) DeleteSecret(cr *v1alpha1.SecretKMS) error {
	return nil
}

// // newbusyBoxPod demonstrates how to create a busybox pod
// func newbusyBoxPod(cr *v1alpha1.SecretKMS) *corev1.Pod {
// 	labels := map[string]string{
// 		"app": "busy-box",
// 	}
// 	return &corev1.Pod{
// 		TypeMeta: metav1.TypeMeta{
// 			Kind:       "Pod",
// 			APIVersion: "v1",
// 		},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "busy-box",
// 			Namespace: cr.Namespace,
// 			OwnerReferences: []metav1.OwnerReference{
// 				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
// 					Group:   v1alpha1.SchemeGroupVersion.Group,
// 					Version: v1alpha1.SchemeGroupVersion.Version,
// 					Kind:    "SecretKMS",
// 				}),
// 			},
// 			Labels: labels,
// 		},
// 		Spec: corev1.PodSpec{
// 			Containers: []corev1.Container{
// 				{
// 					Name:    "busybox",
// 					Image:   "docker.io/busybox",
// 					Command: []string{"sleep", "3600"},
// 				},
// 			},
// 		},
// 	}
// }
