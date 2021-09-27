package oskops

import (
	"context"
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kupenstack/kupenstack/apis/cluster/v1alpha1"
	ksk "github.com/kupenstack/kupenstack/oskops/apis/v1alpha1"
)

var (
	log = ctrl.Log.WithName("kupenstack").WithName("sync").WithName("osknode")
)

var client k8sclient.Client

// TODO: Improve execution using context.Context properly
func ManageOskNodes(c k8sclient.Client, kupenstackConfig string) {
	client = c

	for {
		time.Sleep(10 * time.Second)

		k8snodes, err := getK8sNodes()
		if err != nil {
			log.Error(err, "Failed to get list of k8s nodes.")
			continue
		}

		osknodes, err := getOsksNodes()
		if err != nil {
			log.Error(err, "Failed to get list of osknodes.")
			continue
		}

		err = removeExtraOskNodes(k8snodes, osknodes)
		if err != nil {
			log.Error(err, "Failed to remove extra osknodes.")
			continue
		}

		err = syncOskNodes(k8snodes, osknodes, kupenstackConfig)
		if err != nil {
			log.Error(err, "Failed to sync osknodes objects with kubernetes nodes.")
			continue
		}

		time.Sleep(30 * time.Second)
	}
}

func getK8sNodes() ([]core.Node, error) {
	var nodeList core.NodeList
	err := client.List(context.Background(), &nodeList)
	if err != nil {
		return nil, err
	}
	return nodeList.Items, nil
}

func getOsksNodes() ([]v1alpha1.OpenstackNode, error) {
	var oskNodeList v1alpha1.OpenstackNodeList
	err := client.List(context.Background(), &oskNodeList)
	if err != nil {
		return nil, err
	}
	return oskNodeList.Items, nil
}

func removeExtraOskNodes(k8snodes []core.Node, osknodes []v1alpha1.OpenstackNode) error {

	for _, osknode := range osknodes {

		isExtra := true
		for _, k8snode := range k8snodes {
			if osknode.Name == k8snode.Name {
				isExtra = false
				break
			}
		}

		if isExtra {
			err := client.Delete(context.Background(), &osknode)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func syncOskNodes(k8snodes []core.Node, osknodes []v1alpha1.OpenstackNode, kupenstackConfig string) error {

	for _, k8snode := range k8snodes {

		foundAt := -1
		for i, osknode := range osknodes {
			if osknode.Name == k8snode.Name {
				foundAt = i
				break
			}
		}

		if foundAt == -1 {
			// create
			err := createOskNode(k8snode, kupenstackConfig)
			if err != nil {
				return err
			}
		} else {
			// check for update
			err := syncOskNode(k8snode, osknodes[foundAt], kupenstackConfig)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func createOskNode(node core.Node, kupenstackConfig string) error {

	cfg, err := ReadKupenStackConfiguration(kupenstackConfig)
	if err != nil {
		return err
	}

	nodeRole := desiredNodeRole(node, cfg)

	newNode := v1alpha1.OpenstackNode{
		ObjectMeta: metav1.ObjectMeta{
			Name: node.Name,
			Annotations: map[string]string{
				"node-role": nodeRole,
			},
		},
		Spec: v1alpha1.OpenstackNodeSpec{
			Occp: cfg.Spec.DefaultProfile,
		},
	}

	return client.Create(context.Background(), &newNode)
}

func syncOskNode(node core.Node, osknode v1alpha1.OpenstackNode, kupenstackConfig string) error {

	cfg, err := ReadKupenStackConfiguration(kupenstackConfig)
	if err != nil {
		return err
	}

	nodeRole := desiredNodeRole(node, cfg)

	if osknode.Annotations["node-role"] != nodeRole ||
		osknode.Spec.Occp.Name != cfg.Spec.DefaultProfile.Name ||
		osknode.Spec.Occp.Namespace != cfg.Spec.DefaultProfile.Namespace {

		if osknode.Annotations == nil {
			osknode.Annotations = make(map[string]string)
		}
		osknode.Annotations["node-role"] = nodeRole
		osknode.Spec.Occp.Name = cfg.Spec.DefaultProfile.Name
		osknode.Spec.Occp.Namespace = cfg.Spec.DefaultProfile.Namespace

		return client.Update(context.Background(), &osknode)
	}

	return nil
}

func desiredNodeRole(node core.Node, cfg ksk.KupenstackConfiguration) string {
	nodeRole := "compute"
	_, controlplane := node.Labels["node-role.kubernetes.io/control-plane"]
	_, master := node.Labels["node-role.kubernetes.io/master"]
	if controlplane || master {
		nodeRole = "control"
	}

	for _, n := range cfg.Spec.Nodes {
		if n.Name == node.Name {
			if n.Type != "" {
				nodeRole = n.Type
			}
		}
	}

	return nodeRole
}
