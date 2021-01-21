package util

import (
	"fmt"
	"strings"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesConfig(local bool, masterURL string) (config *restclient.Config, err error) {
	if local {
		configAccess := clientcmd.NewDefaultPathOptions()
		config, err = clientcmd.BuildConfigFromKubeconfigGetter(masterURL, configAccess.GetStartingConfig)
		if err != nil {
			return nil, err
		}
	} else {
		// creates the in-cluster config
		config, err = restclient.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func AppendUniqueElements(left, right []string) []string {
	set := map[string]struct{}{}
	for _, item := range left {
		set[item] = struct{}{}
	}
	for _, item := range right {
		set[item] = struct{}{}
	}

	var deduplicatedSlice []string
	for item := range set {
		deduplicatedSlice = append(deduplicatedSlice, item)
	}

	return deduplicatedSlice
}

func ParseLabelSelectors(labelSelector string) (map[string]string, error) {
	labels := map[string]string{}

	selectors := strings.Split(labelSelector, ",")
	for _, s := range selectors {
		keyValue := strings.Split(s, "=")
		if len(keyValue) != 2 {
			return nil, fmt.Errorf("invalid label selector, expected: key1=value1,key2=value2")
		}
		labels[keyValue[0]] = keyValue[1]
	}
	return labels, nil
}
