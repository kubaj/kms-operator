# KMS operator for Kubernetes

Operator for decryption of data for Kubernetes. This allows storing encrypted credentials or sensitive data in the repository. These are decrypted on the fly when data are created in a Kubernetes.

Currently supported KMS provider is Google Cloud. 

Resource for encrypted resource is similar to v1/Secret resource:

```yaml
apiVersion: kubaj.kms/v1alpha1
kind: SecretKMS
metadata:
  name: example-service-account
spec:
  secret: example-service-account
  file: credentials.json
  provider:
    google-cloud:
      project: gcp-project  # Google Cloud project
      location: global      # KMS Location
      keyring: testring     # Name of the keyring
      key: test             # Name of the key
      data: CiQAFRg31wZQ1pHlR4bBAU8O7nrlz/QEkeKUyrLRsgD92CzIWxkSaQAJc5gIwtzhUZXW9vt1d3+oVl2i+l+tPrUMCN59zybemHro2Y6Gyzrgn0YQ2r3QDR1V+nFMcAvnsCgbInEELhJdXwH/SIRDIHCVVyQqlLr2xEmVXsZVdd3XVH2ivNFEP54XihkRBBaCCg==
```

After creating this resource, operator decrypts data using specified provider and creates v1/Secret with decrypted data:

```yaml
apiVersion: v1
kind: Secret
metadata: 
  name: example-service-account
data:
  credentials.json: dGhlIG1vc3Qgc2VjcmV0IHNlY3JldCBpbiB0aGUgd29ybGQgb2Ygc2VjcmV0cw==
```

[![asciicast](https://asciinema.org/a/Bo54tevSIrogyIZG7rrtzrVeR.png)](https://asciinema.org/a/Bo54tevSIrogyIZG7rrtzrVeR)

# Installation

Cluster with RBAC enabled:
```
$ git clone https://github.com/kubaj/kms-operator.git
$ kubectl apply -n kube-system -f kms-operator/deploy/operator_rbac.yaml
```

Cluster with RBAC disabled (not recommended):
```
$ git clone https://github.com/kubaj/kms-operator.git
$ kubectl apply -n kube-system -f kms-operator/deploy/operator.yaml
```

By default, Google Cloud provider is enabled, you have to create secret that contains Google Cloud service account with KMS decryption scope. To create secret from file:

```
$ kubectl create secret generic google-sa-kms -n kube-system --from-file=credentials.json=/path/to/service-account.json
```


# SecretKMS resource creation

## Google Cloud

Unencrypted data are in file `plaintext.txt`. Encrypt using gcloud sdk:
```
$ gcloud kms encrypt --location global --keyring testring --key test --plaintext-file=plaintext.txt --ciphertext-file=encrypted.bin
```

Encode encrypted data in Base64 and put them to resource:

```
$ cat encrypted.bin | base64

CiQAFRg31wZQ1pHlR4bBAU8O7nrlz/QEkeKUyrLRsgD92CzIWxkSaQAJc5gIwtzhUZXW9vt1d3+oVl2i+l+tPrUMCN59zybemHro2Y6Gyzrgn0YQ2r3QDR1V+nFMcAvnsCgbInEELhJdXwH/SIRDIHCVVyQqlLr2xEmVXsZVdd3XVH2ivNFEP54XihkRBBaCCg==
```

Final resource will look like this:
```
apiVersion: kubaj.kms/v1alpha1
kind: SecretKMS
metadata:
  name: example-service-account
spec:
  secret: example-service-account
  file: credentials.json
  provider:
    google-cloud:
      project: gcp-project  # Google Cloud project
      location: global      # KMS Location
      keyring: testring     # Name of the keyring
      key: test             # Name of the key
      data: CiQAFRg31wZQ1pHlR4bBAU8O7nrlz/QEkeKUyrLRsgD92CzIWxkSaQAJc5gIwtzhUZXW9vt1d3+oVl2i+l+tPrUMCN59zybemHro2Y6Gyzrgn0YQ2r3QDR1V+nFMcAvnsCgbInEELhJdXwH/SIRDIHCVVyQqlLr2xEmVXsZVdd3XVH2ivNFEP54XihkRBBaCCg==
```
