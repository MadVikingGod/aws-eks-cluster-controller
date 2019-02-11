package eks

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"

	componentsv1alpha1 "github.com/awslabs/aws-eks-cluster-controller/pkg/apis/components/v1alpha1"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/runtime"

)

var componentsToDelete = []runtime.Object{
	&componentsv1alpha1.ConfigMapList{},
	&componentsv1alpha1.DeploymentList{},
	&componentsv1alpha1.IngressList{},
	&componentsv1alpha1.SecretList{},
	&componentsv1alpha1.ServiceList{},
	&componentsv1alpha1.ServiceAccountList{},
	&componentsv1alpha1.ClusterRoleList{},
	&componentsv1alpha1.ClusterRoleBindingList{},
}
var resources = []schema.GroupKind{
	componentsv1alpha1.SchemeGroupVersion.WithKind("ConfigMapsList"),
	componentsv1alpha1.SchemeGroupVersion.WithKind("DeploymentList"),
	componentsv1alpha1.SchemeGroupVersion.WithKind("IngressList"),
	componentsv1alpha1.SchemeGroupVersion.WithKind("SecretList"),
	componentsv1alpha1.SchemeGroupVersion.WithKind("ServiceList"),
	componentsv1alpha1.SchemeGroupVersion.WithKind("ServiceAccountList"),
	componentsv1alpha1.SchemeGroupVersion.WithKind("ClusterRoleList"),
	componentsv1alpha1.SchemeGroupVersion.WithKind("ClusterRolebindingList"),
}

func deleteDynamicComponents(ownerName, ownerNamespace string, c client.Client, logger *zap.Logger) (count int, err error) {

	listOptions := client.MatchingLabels(map[string]string{
		"eks.owner.name":      ownerName,
		"eks.owner.namespace": ownerNamespace,
		"eks.needsdeleting":   "true",
	})
	for _, resource := range resources {
		list := &unstructured.UnstructuredList{}
		list.SetGroupVersionKind(resource)
		
		c.List(context.TODO(), listOptions, list)
		

		
		if len(list.Items) > 0 {
			count += len(list.Items)
			err = client.DeleteCollection(&metav1.DeleteOptions{}, listOptions)
			if err != nil {
				logger.Error("error deleting resource", zap.String("resource", resource.String()), zap.Error(err))
				return count, err
			}
		}
	}
	return count, nil
}

func deleteComponents(ownerName, ownerNamespace string, c client.Client, kc kclient.Interface, logger *zap.Logger) (int, error) {



	delete := []runtime.Object{}
	for _, componentList := range componentsToDelete {
		list := componentList.DeepCopyObject()
		err := c.List(context.TODO(), client.MatchingLabels(map[string]string{
			"eks.owner.name":      ownerName,
			"eks.owner.namespace": ownerNamespace,
			"eks.needsdeleting":   "true",
		}), list)

		
		
		if err != nil {
			return 0, err
		}
		switch l := list.(type) {
		case *componentsv1alpha1.ConfigMapList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		case *componentsv1alpha1.DeploymentList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		case *componentsv1alpha1.IngressList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		case *componentsv1alpha1.SecretList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		case *componentsv1alpha1.ServiceList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		case *componentsv1alpha1.ClusterRoleList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		case *componentsv1alpha1.ClusterRoleBindingList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		case *componentsv1alpha1.ServiceAccountList:
			for _, obj := range l.Items {
				item := obj
				delete = append(delete, &item)
			}
		default:
			logger.Error("Got object type we didn't understand", zap.Any("object", l))
			return 0, fmt.Errorf("unknown type error")
		}

	}

	errs := []error{}
	for _, obj := range delete {
		err := c.Delete(context.TODO(), obj)
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return len(delete), fmt.Errorf("error deleting objects %v", errs)
	}

	return len(delete), nil
}
