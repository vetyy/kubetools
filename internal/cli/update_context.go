package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/vetyy/kubetools/internal/cli/genericflags"
	cmdutil "github.com/vetyy/kubetools/internal/cli/util"
	log "github.com/vetyy/kubetools/internal/logging"
)

type UpdateContextFlags struct {
	*genericflags.GenericFlags
	DeleteCluster string
}

func NewCmdUpdateContext(flags *genericflags.GenericFlags) *cobra.Command {
	o := &UpdateContextFlags{
		GenericFlags: flags,
	}
	cmd := &cobra.Command{
		Use:   "update-context",
		Short: "Update contexts",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.PersistentFlags().StringVarP(&o.DeleteCluster, "delete-cluster", "", "", "Cluster")
	return cmd
}

func (o *UpdateContextFlags) Complete() (err error) {
	return o.GenericFlags.Complete()
}

func (o *UpdateContextFlags) Run() (err error) {
	if o.DeleteCluster != "" {
		contexts := o.Config.List(o.DeleteCluster)
		if len(contexts) == 0 {
			log.Println("Nothing to delete.")
			return nil
		}
		return o.deleteContexts(contexts)
	}

	currentNamespaces, err := o.listNamespaces()
	if err != nil {
		return err
	}

	currentContexts := o.Config.List(o.Context.Cluster)
	var createdContexts []string
	for _, namespace := range currentNamespaces {
		contextName := o.Config.ContextName(o.Context.Cluster, namespace)
		if _, ok := currentContexts[contextName]; !ok {
			o.Config.AddContext(o.Context.Cluster, namespace)
			createdContexts = append(createdContexts, contextName)
		}
		delete(currentContexts, contextName) // Non-existent context will be  in allContexts after all iterations
	}

	err = o.Config.Update()
	if err != nil {
		return fmt.Errorf("could not update config: %v", err)
	}

	if len(createdContexts) > 0 {
		log.Success("Contexts created:")
		log.Println(strings.Join(createdContexts, "\n"))
	}

	if len(currentContexts) > 0 {
		return o.deleteContexts(currentContexts)
	}

	return nil
}

func (o *UpdateContextFlags) deleteContexts(contexts map[string]*api.Context) error {
	contextsNames := ""
	for name := range contexts {
		contextsNames += name + "\n"
	}

	log.Warningln("Following contexts are not used anymore:")
	log.Println(contextsNames)
	log.Important("Do you want to delete them? [Y/n]: ")
	confirmed, err := cmdutil.AskForConfirmation()
	if err != nil {
		return err
	}

	if confirmed {
		err := o.Config.DeleteContexts(contexts)
		if err != nil {
			return err
		}
		log.Success("Contexts deleted successfully")
	}
	return nil
}

func (o *UpdateContextFlags) listNamespaces() ([]string, error) {
	client, err := o.Context.GetKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("could get kubernetes client: %v", err)
	}

	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get list of namespaces: %v", err)
	}

	var namespaceList []string
	for _, ns := range namespaces.Items {
		namespaceList = append(namespaceList, ns.Name)
	}

	return namespaceList, nil
}
