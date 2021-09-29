// This package contains helper functions for OpenstackNodes type CR
// in kupenstack.
package osknode

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ksktypes "github.com/kupenstack/kupenstack/pkg/kupenstack/types/v1alpha1"
)

// get unstructured node list
// get unstructured node // Get
// convert to struct // osknode.AsStruct()
// conver struct to unstructed node //  osknode.AsMap()

// Returns list of all OpenstackNode resources as k8s-UnstructuredList-type
func GetList(ctx context.Context, c client.Client) (*unstructured.UnstructuredList, error) {
	oskNodeList := &unstructured.UnstructuredList{}
	oskNodeList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "cluster.kupenstack.io",
		Kind:    "OpenstackNodeList",
		Version: "v1alpha1",
	})
	err := c.List(ctx, oskNodeList)
	if err != nil {
		return nil, err
	}
	return oskNodeList, nil
}

// Returns OpenstackNode resources with name `NodeName` as k8s-Unstructured-type
func Get(ctx context.Context, c client.Client, name string) (*unstructured.Unstructured, error) {
	nodeName := types.NamespacedName{Name: name}
	oskNode := &unstructured.Unstructured{}
	oskNode.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "cluster.kupenstack.io",
		Kind:    "OpenstackNode",
		Version: "v1alpha1",
	})
	err := c.Get(ctx, nodeName, oskNode)
	if err != nil {
		return nil, err
	}
	return oskNode, nil
}

// Takes OpenstackNode as k8s-Unstructured-type and returns it as correesponding struct
func AsStruct(in *unstructured.Unstructured) (*ksktypes.OpenstackNode, error) {
	jsonString, err := json.Marshal(in.Object)
	if err != nil {
		return nil, err
	}

	out := &ksktypes.OpenstackNode{}
	err = json.Unmarshal(jsonString, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// Todo:
// Takes OpenstackNode as struct and returns it as correesponding map[string]interface{}
// Usage: the reurned map[string]interface{} can be stored back in k8s-unstructured-types
//         unstructedtypeNode.Object = osknode.AsMap(mynodestruct)
//         client.Update(ctx, &unstructedtypeNode)
//
// func AsMap(in *ksktypes.OpenstackNode) (map[string]interface{}, error) {
// }
