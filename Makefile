# Image URL to use all building/pushing image targets
OPERATOR_IMG ?= kubeagentic/operator:latest
AGENT_IMG ?= kubeagentic/agent:latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/manager main.go

.PHONY: run
run: fmt vet ## Run a controller from your host.
	go run ./main.go

.PHONY: docker-build-operator
docker-build-operator: ## Build docker image for operator.
	docker build -f Dockerfile.operator -t ${OPERATOR_IMG} .

.PHONY: docker-build-agent
docker-build-agent: ## Build docker image for agent.
	docker build -f Dockerfile.agent -t ${AGENT_IMG} .

.PHONY: docker-build-all
docker-build-all: docker-build-operator docker-build-agent ## Build all docker images.

.PHONY: docker-push-operator
docker-push-operator: ## Push operator docker image.
	docker push ${OPERATOR_IMG}

.PHONY: docker-push-agent
docker-push-agent: ## Push agent docker image.
	docker push ${AGENT_IMG}

.PHONY: docker-push-all
docker-push-all: docker-push-operator docker-push-agent ## Push all docker images.

##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: install-crd
install-crd: ## Install CRDs into the K8s cluster specified in ~/.kube/config.
	kubectl apply -f crd/agent-crd.yaml

.PHONY: uninstall-crd
uninstall-crd: ## Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
	kubectl delete -f crd/agent-crd.yaml --ignore-not-found=$(ignore-not-found)

.PHONY: deploy-namespace
deploy-namespace: ## Create namespace for the operator.
	kubectl apply -f deploy/namespace.yaml

.PHONY: deploy-rbac
deploy-rbac: ## Deploy RBAC manifests to the K8s cluster.
	kubectl apply -f deploy/rbac.yaml

.PHONY: deploy-operator
deploy-operator: ## Deploy operator to the K8s cluster.
	kubectl apply -f deploy/operator.yaml

.PHONY: deploy-all
deploy-all: deploy-namespace install-crd deploy-rbac deploy-operator ## Deploy everything to the K8s cluster.

.PHONY: deploy-complete
deploy-complete: ## Deploy everything using the all-in-one manifest.
	kubectl apply -f deploy/all.yaml

.PHONY: undeploy
undeploy: ## Undeploy controller from the K8s cluster.
	kubectl delete -f deploy/operator.yaml --ignore-not-found=$(ignore-not-found)
	kubectl delete -f deploy/rbac.yaml --ignore-not-found=$(ignore-not-found)
	kubectl delete -f deploy/namespace.yaml --ignore-not-found=$(ignore-not-found)

##@ Examples

.PHONY: deploy-examples
deploy-examples: ## Deploy example agents.
	kubectl apply -f examples/

.PHONY: undeploy-examples
undeploy-examples: ## Remove example agents.
	kubectl delete -f examples/ --ignore-not-found=$(ignore-not-found)

##@ Complete Deployment

.PHONY: complete-deploy
complete-deploy: docker-build-all docker-push-all deploy-all ## Build, push, and deploy everything.

.PHONY: dev-deploy
dev-deploy: docker-build-all deploy-all ## Build and deploy for development (without push).

##@ Testing

.PHONY: test-standalone
test-standalone: ## Test agent in standalone Python mode.
	./local-testing/test-local.sh standalone

.PHONY: test-docker
test-docker: ## Test with Docker Compose (all providers).
	./local-testing/test-local.sh docker

.PHONY: test-kubernetes
test-kubernetes: ## Deploy and test on local Kubernetes.
	./local-testing/test-local.sh kubernetes

.PHONY: test-basic
test-basic: ## Run basic functionality tests.
	./local-testing/test-local.sh basic

.PHONY: test-all
test-all: test-standalone test-docker test-kubernetes ## Run all test suites.

.PHONY: test-agent-openai
test-agent-openai: ## Test OpenAI agent example.
	kubectl apply -f examples/openai-agent.yaml
	@echo "Waiting for agent to be ready..."
	kubectl wait --for=condition=Ready agent/customer-support-agent --timeout=300s
	@echo "Agent is ready! You can now test it with:"
	@echo "kubectl port-forward service/customer-support-agent-service 8080:80"

.PHONY: test-cleanup
test-cleanup: ## Clean up test resources.
	./local-testing/test-local.sh clean

##@ Utilities

.PHONY: logs-operator
logs-operator: ## Show operator logs.
	kubectl logs -n kubeagentic-system deployment/kubeagentic-operator -f

.PHONY: status
status: ## Show status of all components.
	@echo "=== Operator Status ==="
	kubectl get deployment -n kubeagentic-system kubeagentic-operator
	@echo
	@echo "=== Agents ==="
	kubectl get agents -A
	@echo
	@echo "=== Agent Pods ==="
	kubectl get pods -l kubeagentic.ai/agent

.PHONY: clean
clean: ## Clean up build artifacts.
	rm -rf bin/
	rm -f cover.out
