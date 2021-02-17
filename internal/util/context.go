package util

import (
	"fmt"
	"strings"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Context struct {
	Cluster   string `yaml:"cluster"`
	Namespace string `yaml:"namespace"`

	restConfig *restclient.Config
}

func (c *Context) GetKubernetesClient() (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(c.restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return clientSet, nil
}

func (c *Context) GetKubernetesDynamicClient() (dynamic.Interface, error) {
	client, err := dynamic.NewForConfig(c.restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return client, nil
}

func NewContext(cluster string, namespace string, dir string) (context *Context, err error) {
	configAccess := clientcmd.NewDefaultPathOptions()
	cmdConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	if cluster == "" {
		if cmdConfig.CurrentContext == "" {
			return nil, fmt.Errorf("current-context is not set")
		}

		currContext := cmdConfig.Contexts[cmdConfig.CurrentContext]
		cluster = currContext.Cluster
	}

	if namespace == "" {
		if cmdConfig.CurrentContext == "" {
			return nil, fmt.Errorf("current-context is not set")
		}

		currContext := cmdConfig.Contexts[cmdConfig.CurrentContext]
		namespace = currContext.Namespace
	}

	if _, ok := cmdConfig.Clusters[cluster]; !ok {
		return nil, fmt.Errorf("cluster not found")
	}

	restConfig, err := GetKubernetesConfig(true, cmdConfig.Clusters[cluster].Server)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster config: %v", err)
	}
	return &Context{
		Cluster:    cluster,
		Namespace:  namespace,
		restConfig: restConfig,
	}, nil
}

type Config struct {
	cmdConfig    *clientcmdapi.Config
	configAccess *clientcmd.PathOptions
}

func (c *Config) ContextName(cluster string, namespace string) string {
	clusterParts := strings.Split(cluster, "/")
	cluster = clusterParts[len(clusterParts)-1]
	return fmt.Sprintf("%s/%s", namespace, cluster)
}

func (c *Config) SetCurrentContext(cluster string, namespace string) {
	c.cmdConfig.CurrentContext = c.ContextName(cluster, namespace)
}

func (c *Config) Update() (err error) {
	return clientcmd.ModifyConfig(c.configAccess, *c.cmdConfig, true)
}

func (c *Config) List(cluster string) map[string]*api.Context {
	contexts := map[string]*api.Context{}
	for name, ctx := range c.cmdConfig.Contexts {
		if ctx.Cluster == cluster {
			contexts[name] = ctx
		}
	}
	return contexts
}

func (c *Config) DeleteContexts(contexts map[string]*api.Context) error {
	for name := range contexts {
		delete(c.cmdConfig.Contexts, name)
	}
	return c.Update()
}

func (c *Config) AddContext(cluster string, namespace string) {
	newContext := *api.NewContext()
	newContext.Cluster = cluster
	newContext.Namespace = namespace
	newContext.AuthInfo = cluster

	contextName := c.ContextName(cluster, namespace)
	c.cmdConfig.Contexts[contextName] = &newContext
	c.cmdConfig.CurrentContext = contextName
}

func NewConfig() (config *Config, err error) {
	configAccess := clientcmd.NewDefaultPathOptions()
	cmdConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		cmdConfig:    cmdConfig,
		configAccess: configAccess,
	}, nil
}
