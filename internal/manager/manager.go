package manager

// IngressManager is an interface for managing state between K8s Ingress and OCI LBs.
type IngressManager interface {
	// EnsureIngress will ensure that an OCI load balancer exists and is
	// configured correctly for the provided ingress object.
	EnsureIngress()
	// EnsureIngressDeleted ensures that all OCI resources associated with
	// an Ingress are removed.
	EnsureIngressDeleted()
}

type OCIIngressManager struct {
}

func NewOCIIngressManager() *OCIIngressManager {
	return nil
}

func (m *OCIIngressManager) EnsureIngress() {

}

func (m *OCIIngressManager) EnsureIngressDeleted() {

}