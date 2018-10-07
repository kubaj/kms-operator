package main

import (
	"context"
	"flag"
	"runtime"
	"time"

	stub "github.com/kubaj/kms-operator/pkg/stub"
	"github.com/kubaj/kms-operator/utils/clients"
	sdk "github.com/operator-framework/operator-sdk/pkg/sdk"
	k8sutil "github.com/operator-framework/operator-sdk/pkg/util/k8sutil"
	sdkVersion "github.com/operator-framework/operator-sdk/version"

	"github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	providerGoogle       = flag.Bool("google-provider", false, "Enable Google Cloud provider")
	serviceAccountGoogle = flag.String("google-service-account", "", "Path to Google Cloud service acccount to override default service account")
)

func printVersion() {
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	logrus.Infof("operator-sdk Version: %v", sdkVersion.Version)
}

func main() {
	printVersion()
	flag.Parse()

	googleCloudKMS, err := clients.GetGoogleCloudKMS(*providerGoogle, *serviceAccountGoogle)

	sdk.ExposeMetricsPort()

	resource := "kubaj.kms/v1alpha1"
	kind := "SecretKMS"
	namespace, err := k8sutil.GetWatchNamespace()

	logrus.Infof("Watching namespace %s\n", namespace)
	if err != nil {
		logrus.Fatalf("failed to get watch namespace: %v", err)
	}
	resyncPeriod := time.Duration(60) * time.Second
	logrus.Infof("Watching %s, %s, %s, %d", resource, kind, namespace, resyncPeriod)
	sdk.Watch(resource, kind, namespace, resyncPeriod)
	sdk.Handle(stub.NewHandler(googleCloudKMS))
	sdk.Run(context.TODO())
}
