// Code generated by gotestmd DO NOT EDIT.
package basic

import (
	"github.com/stretchr/testify/suite"

	"github.com/networkservicemesh/integration-tests/extensions/base"
	"github.com/networkservicemesh/integration-tests/suites/spire"
)

type Suite struct {
	base.Suite
	spireSuite spire.Suite
}

func (s *Suite) SetupSuite() {
	var base interface{} = &s.Suite
	if v, ok := base.(suite.SetupAllSuite); ok {
		v.SetupSuite()
	}
	suite.Run(s.T(), &s.spireSuite)
	r := s.Runner("../deployments-k8s/examples/basic")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns nsm-system`)
	})
	r.Run(`kubectl create ns nsm-system`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/nsm-system/sa/registry-k8s-sa \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:nsm-system \` + "\n" + `-selector k8s:sa:registry-k8s-sa`)
	r.Run(`kubectl apply -k .`)
}
func (s *Suite) TestLocalConnection() {
	r := s.Runner("../deployments-k8s/examples/LocalConnection")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODE=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}')[0])`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../apps/kernel-nsc` + "\n" + `- ../../apps/kernel-nse` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODE}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`kubectl logs -l app=nsc -n ${NAMESPACE} | grep "All client init operations are done."`)
}
func (s *Suite) TestRemoteConnection() {
	r := s.Runner("../deployments-k8s/examples/RemoteConnection")
	s.T().Cleanup(func() {
		r.Run(`kubectl delete ns ${NAMESPACE}`)
	})
	r.Run(`NAMESPACE=($(kubectl create -f namespace.yaml)[0])` + "\n" + `NAMESPACE=${NAMESPACE:10}`)
	r.Run(`kubectl exec -n spire spire-server-0 -- \` + "\n" + `/opt/spire/bin/spire-server entry create \` + "\n" + `-spiffeID spiffe://example.org/ns/${NAMESPACE}/sa/default \` + "\n" + `-parentID spiffe://example.org/ns/spire/sa/spire-agent \` + "\n" + `-selector k8s:ns:${NAMESPACE} \` + "\n" + `-selector k8s:sa:default`)
	r.Run(`NODES=($(kubectl get nodes -o go-template='{{range .items}}{{ if not .spec.taints  }}{{index .metadata.labels "kubernetes.io/hostname"}} {{end}}{{end}}'))`)
	r.Run(`cat > kustomization.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: kustomize.config.k8s.io/v1beta1` + "\n" + `kind: Kustomization` + "\n" + `` + "\n" + `namespace: ${NAMESPACE}` + "\n" + `` + "\n" + `bases:` + "\n" + `- ../../apps/kernel-nsc` + "\n" + `- ../../apps/kernel-nse` + "\n" + `` + "\n" + `patchesStrategicMerge:` + "\n" + `- patch-nsc.yaml` + "\n" + `- patch-nse.yaml` + "\n" + `EOF`)
	r.Run(`cat > patch-nsc.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nsc` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      containers:` + "\n" + `        - name: nsc` + "\n" + `          env:` + "\n" + `            - name: NSM_NETWORK_SERVICES` + "\n" + `              value: kernel://icmp-responder/nsm-1` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[0]}` + "\n" + `EOF`)
	r.Run(`cat > patch-nse.yaml <<EOF` + "\n" + `---` + "\n" + `apiVersion: apps/v1` + "\n" + `kind: Deployment` + "\n" + `metadata:` + "\n" + `  name: nse` + "\n" + `spec:` + "\n" + `  template:` + "\n" + `    spec:` + "\n" + `      nodeSelector:` + "\n" + `        kubernetes.io/hostname: ${NODES[1]}` + "\n" + `EOF`)
	r.Run(`kubectl apply -k .`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nsc -n ${NAMESPACE}`)
	r.Run(`kubectl wait --for=condition=ready --timeout=1m pod -l app=nse -n ${NAMESPACE}`)
	r.Run(`kubectl logs -l app=nsc -n ${NAMESPACE} | grep "All client init operations are done."`)
}
