package checks

import (
	"errors"
	"k8s.io/client-go/kubernetes"
	"monitoring/spec"
)

// TODO: Implement
func checkK8S(checkSpec *spec.CheckSpec, c *CheckHandler) (*map[string]string, error) {
	if nil == checkSpec.DSNParams {
		return nil, errors.New("DSNParams token must be set for K8S check!")
	}
	return nil, spec.NotImplemented
}

func handleK8SCheckMethods(checkSpec *spec.CheckSpec, client *kubernetes.Clientset) error {
	switch checkSpec.Method {
	case spec.CHANGE:

		return spec.NotImplemented
	case spec.QUERY:
		return spec.NotImplemented
	case spec.CONTAINS:
		return spec.NotImplemented
	case spec.STATUS:
		return spec.NotImplemented
	case spec.EXPIRES:
		return spec.NotImplemented
	}

	return spec.NoCheckPerformed
}

func (c *CheckHandler) CheckK8S(spec *spec.CheckSpec) (*map[string]string, error) {
	return checkK8S(spec, c)
}
