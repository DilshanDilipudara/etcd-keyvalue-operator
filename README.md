# etcd-operator

This operator is implemented for managing an external etcd cluster with ArgoCD. It is written using the Operator SDK framework.

## Description
    
Argo watches the Git repository (specifically resources of kind EtcdConfig) and keeps the Git status synchronized with the Kubernetes cluster. The operator monitors the name-last-synced YAML file for existence or modification.

    1. If the YAML file doesn't exist:

        * Update the etcd cluster.
        * Copy the name YAML file as name-last-synced.
    
    2. If the YAML file exists and there is a difference:

        * Update the etcd cluster.
        * Delete the name-last-synced YAML file.
        * Copy the name YAML file as name-last-synced.
    
        ![ETCD Operator Diagram](images/diagram.png)

## Getting Started

### Prerequisites
- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### Create Secrets
Create a secret named etcd-operator-etcd-auth-secrets for authenticating with the etcd cluster.

Update the etcd-secret.yaml file located in config/samples with the following content:
        data:
            etcdEndpoints: 
            etcdPassword: 
            etcdUsername:


### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/etcd-operator:tag
```

or change the Makefile

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

