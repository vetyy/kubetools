package genericflags

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"

	"github.com/vetyy/kubetools/internal/util"
)

type GenericFlags struct {
	Namespace string
	Cluster   string

	LabelSelector        string
	labelSelectorsParsed map[string]string

	KindSelector       string
	kindSelectorParsed string

	Context *util.Context
	Config  *util.Config

	Mapper          meta.RESTMapper
	DiscoveryClient discovery.CachedDiscoveryInterface
}

func (o *GenericFlags) Complete() (err error) {
	if o.labelSelectorsParsed == nil {
		o.labelSelectorsParsed = map[string]string{}
	}

	if o.LabelSelector != "" {
		o.labelSelectorsParsed, err = util.ParseLabelSelectors(o.LabelSelector)
		if err != nil {
			return err
		}
	}

	if o.KindSelector != "" {
		singular, err := o.Mapper.ResourceSingularizer(o.KindSelector)
		if err != nil {
			return fmt.Errorf("invalid kind selector provided: %v", err)
		}
		o.kindSelectorParsed = singular
	}

	o.Context, err = util.NewContext(o.Cluster, o.Namespace, "environment")
	if err != nil {
		return fmt.Errorf("failed to get context: %v", err)
	}
	o.Cluster = o.Context.Cluster
	o.Namespace = o.Context.Namespace

	o.Config, err = util.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}
	return err
}

func (o *GenericFlags) FilterByLabelSelector(labels map[string]string) bool {
	for key, value := range o.labelSelectorsParsed {
		if labels[key] != value {
			return false
		}
	}
	return true
}

func (o *GenericFlags) FilterByKindSelector(kind string) bool {
	if o.kindSelectorParsed == "" {
		return true
	}

	singular, err := o.Mapper.ResourceSingularizer(kind)
	if err != nil {
		return false
	}
	return singular == o.kindSelectorParsed
}

func NewGenericFlags() (*GenericFlags, error) {
	configFlags := genericclioptions.NewConfigFlags(true)
	mapper, err := configFlags.ToRESTMapper()
	if err != nil {
		return nil, fmt.Errorf("failed to create mapper: %v", err)
	}

	discoveryClient, err := configFlags.ToDiscoveryClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %v", err)
	}

	return &GenericFlags{
		Mapper:          mapper,
		DiscoveryClient: discoveryClient,
	}, nil
}
