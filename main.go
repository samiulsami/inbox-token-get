package main

import (
	"context"
	"fmt"
	"gomodules.xyz/encoding/json"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/tools/clientcmd"
	v1alpha12 "kmodules.xyz/resource-metadata/apis/identity/v1alpha1"
	"kmodules.xyz/resource-metadata/client/clientset/versioned/typed/identity/v1alpha1"
	"os"
)

func getJMAPToken() (*string, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("error building kubeconfig, %v", err)
	}

	myClient, err := v1alpha1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error building clientset, %v", err)
	}

	tokenRequest := myClient.InboxTokenRequests()
	tokenResp, err := tokenRequest.Create(context.TODO(), &v1alpha12.InboxTokenRequest{
		TypeMeta: v1.TypeMeta{
			Kind:       "InboxTokenRequest",
			APIVersion: "identity.k8s.appscode.com/v1alpha1",
		},
	}, v1.CreateOptions{})

	if err != nil || tokenResp.Response == nil {
		return nil, fmt.Errorf("error creating token %v", err)
	}
	var tokenMap map[string]string
	if jsonErr := json.Unmarshal([]byte(tokenResp.Response.AdminJWTToken), &tokenMap); jsonErr != nil {
		return nil, fmt.Errorf("error unmarshalling token %v", jsonErr)
	}

	if token, ok := tokenMap["token"]; ok {
		return &token, nil
	}

	return nil, fmt.Errorf("key \"token\" missing in token response")
}

func main() {
	jmapToken, err := getJMAPToken()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", *jmapToken)
}
