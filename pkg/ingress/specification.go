package ingress

import (
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

const (
	// IngressAnnotationLoadBalancerInternal is an annotation for
	// specifying that a load balancer should be internal.
	IngressAnnotationLoadBalancerInternal = "ingress.beta.kubernetes.io/oci-load-balancer-internal"
	// IngressAnnotationLoadBalancerShape is an annotation for
	// specifying the shape of a load balancer. Available shapes include
	// "100Mbps", "400Mbps", and "8000Mbps".
	IngressAnnotationLoadBalancerShape = "ingress.beta.kubernetes.io/oci-load-balancer-shape"
	// IngressAnnotationLoadBalancerCompartment allows for load balancers to be created in a compartment
	// different to that specified in config.
	IngressAnnotationLoadBalancerCompartment = "ingress.beta.kubernetes.io/oci-load-balancer-compartment"
)

// Specification describes the desired state of the OCI load balancer.
// It provides a mapping bridge between K8s and OCI LB.
type Specification struct {
	Name    string
	Ingress *v1beta1.Ingress
	Nodes   []*core_v1.Node
}

// NewSpecification creates a new load balancer specification for a
// given Ingress
func NewSpecification(name string, ingress *v1beta1.Ingress, nodes []*core_v1.Node) Specification {
	return Specification{
		Name:    name,
		Ingress: ingress,
		Nodes:   nodes,
	}
}

// GetLoadBalancerShape will return the load balancer shape required.
// The shape can be controlled by setting ingress object annotations.
func (spec Specification) GetLoadBalancerShape() string {
	return getIngressAnnotationOrDefault(spec.Ingress, IngressAnnotationLoadBalancerShape, "100Mbps")
}

// GetLoadBalancerSubnets will return a list of load balancer subnets based on configuration.
func (spec Specification) GetLoadBalancerSubnets() []string {
	// TODO pull from config
	return []string{
		"ocid1.subnet.oc1.uk-london-1.aaaaaaaaqalydfvmgw7pdw3tittizpoyondib7hedwayyswrrfcrsmc4j7dq",
		"ocid1.subnet.oc1.uk-london-1.aaaaaaaa2tqtopdpynhbjglh3szj2j6h6pwwwohrcanbeyj6dpbiboyuvrza",
	}
}

// GetLoadBalancerCompartment will return the compartment in which a load balancer should exist
// based on either configuration or annotations.
func (spec Specification) GetLoadBalancerCompartment() string {
	// TODO check annotation here and pull from config
	return "ocid1.compartment.oc1..aaaaaaaaob4ckouj3cjmf36ifjkff33wvln5fnnarumafqzpqq7tmbig2n5q"
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
