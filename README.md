# etcd-key-value-add-update-delete-operator

This operator is implemented for managing an external etcd cluster with ArgoCD. It is written using the Operator SDK framework.

# Description

Argo watches the Git repository (specifically resources of kind EtcdConfig) and keeps the Git status synchronized with the Kubernetes cluster. 

    apiVersion: etcd.dilshan.com/v1
    kind: EtcdConfig

# How the ETCD Operator Works
    
### Git Repository:

A YAML file (with kind: EtcdConfig) is stored and managed in a Git repository.

### ArgoCD Watches YAML for Diffs:

ArgoCD continuously monitors the Git repository for changes to the YAML file.
When a difference (diff) is detected in the YAML file compared to what is deployed in the Kubernetes cluster, ArgoCD triggers an action.

### If a Diff is Detected:

ArgoCD Applies the YAML to the Kubernetes Cluster:
    The updated YAML file is applied to the Kubernetes cluster.

Operator Watches applied YAML file and YAML-last-synced:
    The operator monitors both the newly applied YAML file and a YAML-last-synced file (which is a copy of the last applied YAML).

### If the YAML File Exists:

The operator checks for differences between the newly applied YAML and the YAML-last-synced file.

### If Differences Are Found Between YAML and YAML-last-synced:

1. Update the etcd Cluster:

    The operator updates the external etcd cluster based on the new configuration in the YAML file.

2. Delete YAML-last-synced:

    The YAML-last-synced file is deleted as it is no longer up-to-date.

3. Copy Newly Applied YAML to YAML-last-synced:

    The newly applied YAML file is copied to create a new YAML-last-synced file, representing the current state.

### If the ArgoCD Applied YAML File Does Not Exist in the Cluster:

1. Update the External etcd Cluster:

    The external etcd cluster is updated directly based on the configuration that was intended to be applied by ArgoCD.

2. Copy YAML File to YAML-last-synced:

    A copy of this YAML file is saved as YAML-last-synced to ensure the operator has the most recent configuration for future comparisons.

## Getting Started

### Prerequisites
- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### Create Secrets
Create a secret named etcd-operator-etcd-auth-secrets for authenticating with the etcd cluster.

Update the etcd-secret.yaml file located in config/manager with the following content:
        
    data:
        etcdEndpoints: ""
        etcdPassword:  ""
        etcdUsername: ""


### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/etcd-operator:tag
```

or
 
    change the Makefile

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/etcd-operator:tag
```

**If Changes Were Made in the Makefile**

```sh
make deploy 
``` 

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

