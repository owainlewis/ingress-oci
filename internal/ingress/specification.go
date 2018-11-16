package ingress

import (
	"github.com/oracle/oci-go-sdk/loadbalancer"
	"github.com/owainlewis/oci-kubernetes-ingress/internal/config"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

const (
	// IngressAnnotationLoadBalancerVisibility is an annotation for
	// specifying that a load balancer should be public or private.
	// By default all load balancers will be public.
	// Values an be one of "public" or "private"
	IngressAnnotationLoadBalancerVisibility = "ingress.beta.kubernetes.io/oci-load-balancer-visibility"
	// IngressAnnotationLoadBalancerShape is an annotation for
	// specifying the shape of a load balancer. Available shapes include
	// "100Mbps", "400Mbps", and "8000Mbps".
	IngressAnnotationLoadBalancerShape = "ingress.beta.kubernetes.io/oci-load-balancer-shape"
	// IngressAnnotationLoadBalancerCompartment allows for load balancers to be created in a compartment
	// different to that specified in config.
	IngressAnnotationLoadBalancerCompartment = "ingress.beta.kubernetes.io/oci-load-balancer-compartment"
)

const defaultLoadBalancerShape = "100Mbps"

// Specification describes the desired state of the OCI load balancer.
// It provides a mapping bridge between K8s and OCI LB.
type Specification struct {
	Config  config.Config
	Ingress *v1beta1.Ingress
	Nodes   []*core_v1.Node
}

// NewSpecification creates a new load balancer specification for a
// given Ingress
func NewSpecification(configuration config.Config, ingress *v1beta1.Ingress, nodes []*core_v1.Node) Specification {
	return Specification{
		Config:  configuration,
		Ingress: ingress,
		Nodes:   nodes,
	}
}

// GetLoadBalancerShape will return the load balancer shape required.
// The shape can be controlled by setting ingress object annotations.
func (spec Specification) GetLoadBalancerShape() string {
	return getIngressAnnotationOrDefault(spec.Ingress, IngressAnnotationLoadBalancerShape, defaultLoadBalancerShape)
}

// GetLoadBalancerSubnets will return a list of load balancer subnets based on configuration.
func (spec Specification) GetLoadBalancerSubnets() []string {
	return spec.Config.Loadbalancer.Subnets
}

// GetLoadBalancerCompartment will return the compartment in which a load balancer should exist
// based on either configuration or (TODO) annotations.
func (spec Specification) GetLoadBalancerCompartment() string {
	return spec.Config.Loadbalancer.Compartment
}

// LoadBalancerIsPrivate checks if a load balancer should be declared private.
// Visibility can be controlled by annotations on the ingress object.
func (spec Specification) LoadBalancerIsPrivate() bool {
	return getIngressAnnotationOrDefault(spec.Ingress, IngressAnnotationLoadBalancerVisibility, "public") == "private"
}

// GetListeners returns a list of Listeners to create for this specification.
func (spec Specification) GetListeners() map[string]loadbalancer.ListenerDetails {
	return map[string]loadbalancer.ListenerDetails{}
}

// GetBackendSets returns a list of the Backends we need to create for this specification.
func (spec Specification) GetBackendSets() map[string]loadbalancer.BackendSetDetails {
	return map[string]loadbalancer.BackendSetDetails{}
}

// GetPathRouteSets returns a list of the PathRouteSets we need to create for this specification.
func (spec Specification) GetPathRouteSets() map[string]loadbalancer.PathRouteSetDetails {
	return map[string]loadbalancer.PathRouteSetDetails{}
}

// GetCertificates returns a list of the Certificates we need to create for this specification.
func (spec Specification) GetCertificates() map[string]loadbalancer.CertificateDetails {
	return map[string]loadbalancer.CertificateDetails{}
}

// GetLoadBalancerTags returns a map of freeform tags for an ingress load balancer.
func (spec Specification) GetLoadBalancerFreeFormTags() map[string]string {
	return map[string]string{
		"ingress.name": spec.Ingress.Name,
	}
}

func getIngressAnnotationOrDefault(ingress *v1beta1.Ingress, k, defaultValue string) string {
	if value, ok := ingress.Annotations[k]; ok {
		return value
	}
	return defaultValue
}

// getNodeInternalIPAddress will extract the OCI internal node IP address
// for a given node. Since it is impossible to launch an instance without
// an internal (private) IP, we can be sure that one exists.
func getNodeInternalIPAddress(node *core_v1.Node) string {
	for _, addr := range node.Status.Addresses {
		if addr.Type == core_v1.NodeInternalIP {
			return addr.Address
		}
	}
	return ""
}
