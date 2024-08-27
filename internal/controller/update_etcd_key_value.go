package controller

import (
	"context"
	"fmt"
	"time"

	"go.etcd.io/etcd/client/v3"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)


var etcdEndpoints []string
var etcdUsername string
var etcdPassword string
var	secretName = "etcd-operator-etcd-auth-secrets"
var	namespace = "etcd-operator-system"


// auth etcd cluster
func getEtcdConfigFromSecret(ctx context.Context, c client.Client) error{

	// Define the Secret object
	secret := &corev1.Secret{}

	namespacedName := client.ObjectKey{
        Namespace: namespace,
        Name:      secretName,
    }
    
   // Fetch the secret from the cluster
   err := c.Get(ctx, namespacedName, secret)
   if err != nil {
	   if client.IgnoreNotFound(err) != nil {
		   return fmt.Errorf("failed to get secret: %v", err)
	   }
	   // Secret not found, handle it here if necessary
	   fmt.Println("Secret not found")
	   return nil
   }
	// Decode and assign the values from the secret
	if endpointBytes, exists := secret.Data["etcdEndpoints"]; exists {
		etcdEndpoints = []string{string(endpointBytes)}
	} else {
		return fmt.Errorf("etcdEndpoints not found in secret")
	}

	if usernameBytes, exists := secret.Data["etcdUsername"]; exists {
		etcdUsername = string(usernameBytes)
	} else {
		return fmt.Errorf("etcdUsername not found in secret")
	}

	if passwordBytes, exists := secret.Data["etcdPassword"]; exists {
		etcdPassword = string(passwordBytes)
	} else {
		return fmt.Errorf("etcdPassword not found in secret")
	}

	return nil

}

func updateEtcdCluster(key string, value interface{}) error {

	// fmt.Println("endpoint : ",etcdEndpoints)
	// fmt.Println("username : ", etcdUsername)
	// fmt.Println("password : ", etcdPassword)
	// Create a new etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		Username:    etcdUsername,
		Password:    etcdPassword,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("failed to create etcd client: %v", err)
	}
	defer client.Close()

	// Convert value to string
	valueStr := fmt.Sprintf("%v", value)

	// Put the key-value pair into etcd
	ctx := context.Background()
	_, err = client.Put(ctx, key, valueStr)
	if err != nil {
		return fmt.Errorf("failed to update etcd cluster: %v", err)
	}

	return nil
}