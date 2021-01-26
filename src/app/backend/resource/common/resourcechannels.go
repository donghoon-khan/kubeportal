package common

import (
	"context"

	apps "k8s.io/api/apps/v1"
	autoscaling "k8s.io/api/autoscaling/v1"
	batch "k8s.io/api/batch/v1"
	batch2 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	rbac "k8s.io/api/rbac/v1"
	storage "k8s.io/api/storage/v1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"

	api "github.com/donghoon-khan/kubeportal/src/app/backend/api"
)

type ResourceChannels struct {
	ReplicationControllerList   ReplicationControllerListChannel
	ReplicaSetList              ReplicaSetListChannel
	DeploymentList              DeploymentListChannel
	DaemonSetList               DaemonSetListChannel
	JobList                     JobListChannel
	CronJobList                 CronJobListChannel
	ServiceList                 ServiceListChannel
	EndpointList                EndpointListChannel
	IngressList                 IngressListChannel
	PodList                     PodListChannel
	EventList                   EventListChannel
	LimitRangeList              LimitRangeListChannel
	NodeList                    NodeListChannel
	NamespaceList               NamespaceListChannel
	StatefulSetList             StatefulSetListChannel
	ConfigMapList               ConfigMapListChannel
	SecretList                  SecretListChannel
	PersistentVolumeList        PersistentVolumeListChannel
	PersistentVolumeClaimList   PersistentVolumeClaimListChannel
	ResourceQuotaList           ResourceQuotaListChannel
	HorizontalPodAutoscalerList HorizontalPodAutoscalerListChannel
	StorageClassList            StorageClassListChannel
	RoleList                    RoleListChannel
	ClusterRoleList             ClusterRoleListChannel
	RoleBindingList             RoleBindingListChannel
	ClusterRoleBindingList      ClusterRoleBindingListChannel
}

type ServiceListChannel struct {
	List  chan *v1.ServiceList
	Error chan error
}

func GetServiceListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) ServiceListChannel {

	channel := ServiceListChannel{
		List:  make(chan *v1.ServiceList, numReads),
		Error: make(chan error, numReads),
	}
	go func() {
		list, err := client.CoreV1().Services(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []v1.Service
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type IngressListChannel struct {
	List  chan *extensions.IngressList
	Error chan error
}

func GetIngressListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) IngressListChannel {

	channel := IngressListChannel{
		List:  make(chan *extensions.IngressList, numReads),
		Error: make(chan error, numReads),
	}
	go func() {
		list, err := client.ExtensionsV1beta1().Ingresses(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []extensions.Ingress
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type LimitRangeListChannel struct {
	List  chan *v1.LimitRangeList
	Error chan error
}

func GetLimitRangeListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) LimitRangeListChannel {

	channel := LimitRangeListChannel{
		List:  make(chan *v1.LimitRangeList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().LimitRanges(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type NodeListChannel struct {
	List  chan *v1.NodeList
	Error chan error
}

func GetNodeListChannel(client client.Interface, numReads int) NodeListChannel {
	channel := NodeListChannel{
		List:  make(chan *v1.NodeList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().Nodes().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type NamespaceListChannel struct {
	List  chan *v1.NamespaceList
	Error chan error
}

func GetNamespaceListChannel(client client.Interface, numReads int) NamespaceListChannel {
	channel := NamespaceListChannel{
		List:  make(chan *v1.NamespaceList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().Namespaces().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type EventListChannel struct {
	List  chan *v1.EventList
	Error chan error
}

func GetEventListChannel(client client.Interface,
	nsQuery *NamespaceQuery, numReads int) EventListChannel {
	return GetEventListChannelWithOptions(client, nsQuery, api.ListEverything, numReads)
}

func GetEventListChannelWithOptions(client client.Interface,
	nsQuery *NamespaceQuery, options metaV1.ListOptions, numReads int) EventListChannel {
	channel := EventListChannel{
		List:  make(chan *v1.EventList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().Events(nsQuery.ToRequestParam()).List(context.TODO(), options)
		var filteredItems []v1.Event
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type EndpointListChannel struct {
	List  chan *v1.EndpointsList
	Error chan error
}

func GetEndpointListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) EndpointListChannel {
	return GetEndpointListChannelWithOptions(client, nsQuery, api.ListEverything, numReads)
}

func GetEndpointListChannelWithOptions(client client.Interface,
	nsQuery *NamespaceQuery, opt metaV1.ListOptions, numReads int) EndpointListChannel {
	channel := EndpointListChannel{
		List:  make(chan *v1.EndpointsList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().Endpoints(nsQuery.ToRequestParam()).List(context.TODO(), opt)

		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type PodListChannel struct {
	List  chan *v1.PodList
	Error chan error
}

func GetPodListChannel(client client.Interface,
	nsQuery *NamespaceQuery, numReads int) PodListChannel {
	return GetPodListChannelWithOptions(client, nsQuery, api.ListEverything, numReads)
}

func GetPodListChannelWithOptions(client client.Interface, nsQuery *NamespaceQuery,
	options metaV1.ListOptions, numReads int) PodListChannel {

	channel := PodListChannel{
		List:  make(chan *v1.PodList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().Pods(nsQuery.ToRequestParam()).List(context.TODO(), options)
		var filteredItems []v1.Pod
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type ReplicationControllerListChannel struct {
	List  chan *v1.ReplicationControllerList
	Error chan error
}

func GetReplicationControllerListChannel(client client.Interface,
	nsQuery *NamespaceQuery, numReads int) ReplicationControllerListChannel {

	channel := ReplicationControllerListChannel{
		List:  make(chan *v1.ReplicationControllerList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().ReplicationControllers(nsQuery.ToRequestParam()).
			List(context.TODO(), api.ListEverything)
		var filteredItems []v1.ReplicationController
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type DeploymentListChannel struct {
	List  chan *apps.DeploymentList
	Error chan error
}

func GetDeploymentListChannel(client client.Interface,
	nsQuery *NamespaceQuery, numReads int) DeploymentListChannel {

	channel := DeploymentListChannel{
		List:  make(chan *apps.DeploymentList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.AppsV1().Deployments(nsQuery.ToRequestParam()).
			List(context.TODO(), api.ListEverything)
		var filteredItems []apps.Deployment
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type ReplicaSetListChannel struct {
	List  chan *apps.ReplicaSetList
	Error chan error
}

func GetReplicaSetListChannel(client client.Interface,
	nsQuery *NamespaceQuery, numReads int) ReplicaSetListChannel {
	return GetReplicaSetListChannelWithOptions(client, nsQuery, api.ListEverything, numReads)
}

func GetReplicaSetListChannelWithOptions(client client.Interface, nsQuery *NamespaceQuery,
	options metaV1.ListOptions, numReads int) ReplicaSetListChannel {
	channel := ReplicaSetListChannel{
		List:  make(chan *apps.ReplicaSetList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.AppsV1().ReplicaSets(nsQuery.ToRequestParam()).
			List(context.TODO(), options)
		var filteredItems []apps.ReplicaSet
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type DaemonSetListChannel struct {
	List  chan *apps.DaemonSetList
	Error chan error
}

func GetDaemonSetListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) DaemonSetListChannel {
	channel := DaemonSetListChannel{
		List:  make(chan *apps.DaemonSetList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.AppsV1().DaemonSets(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []apps.DaemonSet
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type JobListChannel struct {
	List  chan *batch.JobList
	Error chan error
}

func GetJobListChannel(client client.Interface,
	nsQuery *NamespaceQuery, numReads int) JobListChannel {
	channel := JobListChannel{
		List:  make(chan *batch.JobList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.BatchV1().Jobs(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []batch.Job
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type CronJobListChannel struct {
	List  chan *batch2.CronJobList
	Error chan error
}

func GetCronJobListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) CronJobListChannel {
	channel := CronJobListChannel{
		List:  make(chan *batch2.CronJobList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.BatchV1beta1().CronJobs(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []batch2.CronJob
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type StatefulSetListChannel struct {
	List  chan *apps.StatefulSetList
	Error chan error
}

func GetStatefulSetListChannel(client client.Interface,
	nsQuery *NamespaceQuery, numReads int) StatefulSetListChannel {
	channel := StatefulSetListChannel{
		List:  make(chan *apps.StatefulSetList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		statefulSets, err := client.AppsV1().StatefulSets(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []apps.StatefulSet
		for _, item := range statefulSets.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		statefulSets.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- statefulSets
			channel.Error <- err
		}
	}()

	return channel
}

type ConfigMapListChannel struct {
	List  chan *v1.ConfigMapList
	Error chan error
}

func GetConfigMapListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) ConfigMapListChannel {

	channel := ConfigMapListChannel{
		List:  make(chan *v1.ConfigMapList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().ConfigMaps(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []v1.ConfigMap
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type SecretListChannel struct {
	List  chan *v1.SecretList
	Error chan error
}

func GetSecretListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) SecretListChannel {

	channel := SecretListChannel{
		List:  make(chan *v1.SecretList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().Secrets(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		var filteredItems []v1.Secret
		for _, item := range list.Items {
			if nsQuery.Matches(item.ObjectMeta.Namespace) {
				filteredItems = append(filteredItems, item)
			}
		}
		list.Items = filteredItems
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type RoleListChannel struct {
	List  chan *rbac.RoleList
	Error chan error
}

func GetRoleListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) RoleListChannel {
	channel := RoleListChannel{
		List:  make(chan *rbac.RoleList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.RbacV1().Roles(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type ClusterRoleListChannel struct {
	List  chan *rbac.ClusterRoleList
	Error chan error
}

func GetClusterRoleListChannel(client client.Interface, numReads int) ClusterRoleListChannel {
	channel := ClusterRoleListChannel{
		List:  make(chan *rbac.ClusterRoleList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.RbacV1().ClusterRoles().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type RoleBindingListChannel struct {
	List  chan *rbac.RoleBindingList
	Error chan error
}

func GetRoleBindingListChannel(client client.Interface, nsQuery *NamespaceQuery, numReads int) RoleBindingListChannel {
	channel := RoleBindingListChannel{
		List:  make(chan *rbac.RoleBindingList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.RbacV1().RoleBindings(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type ClusterRoleBindingListChannel struct {
	List  chan *rbac.ClusterRoleBindingList
	Error chan error
}

func GetClusterRoleBindingListChannel(client client.Interface,
	numReads int) ClusterRoleBindingListChannel {
	channel := ClusterRoleBindingListChannel{
		List:  make(chan *rbac.ClusterRoleBindingList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.RbacV1().ClusterRoleBindings().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type PersistentVolumeListChannel struct {
	List  chan *v1.PersistentVolumeList
	Error chan error
}

func GetPersistentVolumeListChannel(client client.Interface,
	numReads int) PersistentVolumeListChannel {
	channel := PersistentVolumeListChannel{
		List:  make(chan *v1.PersistentVolumeList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().PersistentVolumes().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type PersistentVolumeClaimListChannel struct {
	List  chan *v1.PersistentVolumeClaimList
	Error chan error
}

func GetPersistentVolumeClaimListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) PersistentVolumeClaimListChannel {

	channel := PersistentVolumeClaimListChannel{
		List:  make(chan *v1.PersistentVolumeClaimList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().PersistentVolumeClaims(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type CustomResourceDefinitionChannelV1 struct {
	List  chan *apiextensions.CustomResourceDefinitionList
	Error chan error
}

func GetCustomResourceDefinitionChannelV1(client apiextensionsclientset.Interface, numReads int) CustomResourceDefinitionChannelV1 {
	channel := CustomResourceDefinitionChannelV1{
		List:  make(chan *apiextensions.CustomResourceDefinitionList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type CustomResourceDefinitionChannelV1beta1 struct {
	List  chan *apiextensionsv1beta1.CustomResourceDefinitionList
	Error chan error
}

func GetCustomResourceDefinitionChannelV1beta1(client apiextensionsclientset.Interface, numReads int) CustomResourceDefinitionChannelV1beta1 {
	channel := CustomResourceDefinitionChannelV1beta1{
		List:  make(chan *apiextensionsv1beta1.CustomResourceDefinitionList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.ApiextensionsV1beta1().CustomResourceDefinitions().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type ResourceQuotaListChannel struct {
	List  chan *v1.ResourceQuotaList
	Error chan error
}

func GetResourceQuotaListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) ResourceQuotaListChannel {

	channel := ResourceQuotaListChannel{
		List:  make(chan *v1.ResourceQuotaList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.CoreV1().ResourceQuotas(nsQuery.ToRequestParam()).List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type HorizontalPodAutoscalerListChannel struct {
	List  chan *autoscaling.HorizontalPodAutoscalerList
	Error chan error
}

func GetHorizontalPodAutoscalerListChannel(client client.Interface, nsQuery *NamespaceQuery,
	numReads int) HorizontalPodAutoscalerListChannel {
	channel := HorizontalPodAutoscalerListChannel{
		List:  make(chan *autoscaling.HorizontalPodAutoscalerList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.AutoscalingV1().HorizontalPodAutoscalers(nsQuery.ToRequestParam()).
			List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}

type StorageClassListChannel struct {
	List  chan *storage.StorageClassList
	Error chan error
}

func GetStorageClassListChannel(client client.Interface, numReads int) StorageClassListChannel {
	channel := StorageClassListChannel{
		List:  make(chan *storage.StorageClassList, numReads),
		Error: make(chan error, numReads),
	}

	go func() {
		list, err := client.StorageV1().StorageClasses().List(context.TODO(), api.ListEverything)
		for i := 0; i < numReads; i++ {
			channel.List <- list
			channel.Error <- err
		}
	}()

	return channel
}
