# MetaKube Crossplane Provider Examples

This directory contains examples for using the SysEleven MetaKube Crossplane provider.

## Prerequisites

1. A Kubernetes cluster with Crossplane installed
2. The SysEleven MetaKube provider installed
3. A MetaKube API token
4. A MetaKube project ID

## Setup

### 1. Configure Provider Credentials

First, create a secret with your MetaKube API token:

```bash
# Edit the secret file and replace YOUR_METAKUBE_TOKEN_HERE with your actual token
vi examples/setup/secret.yaml

# Apply the secret
kubectl apply -f examples/setup/secret.yaml
```

### 2. Create ProviderConfig

```bash
kubectl apply -f examples/setup/providerconfig.yaml
```

Verify the ProviderConfig:

```bash
kubectl get providerconfigs
```

## Using Compositions

### 1. Install the XRD (CompositeResourceDefinition)

```bash
kubectl apply -f examples/xrd/xrd-metakube-cluster.yaml
```

### 2. Install the Composition

```bash
kubectl apply -f examples/composition/composition-metakube-cluster.yaml
```

### 3. Create a Cluster Claim

Edit the claim file and replace `YOUR_PROJECT_ID_HERE` with your MetaKube project ID:

```bash
vi examples/claim/metakube-cluster-claim.yaml
```

Create the cluster:

```bash
kubectl apply -f examples/claim/metakube-cluster-claim.yaml
```

### 4. Monitor the Cluster Creation

```bash
# Watch the claim
kubectl get metakubeclusters -n default

# Watch the composite resource
kubectl get xmetakubeclusters

# Watch the managed resources
kubectl get clusters.metakube.syseleven-metakube
kubectl get deployments.node.syseleven-metakube

# Check events
kubectl describe metakubecluster my-test-cluster -n default
```

## Direct Resource Creation (Without Compositions)

If you want to create resources directly without using compositions:

### 1. Create a Cluster

```yaml
apiVersion: metakube.syseleven-metakube/v1alpha1
kind: Cluster
metadata:
  name: my-cluster
spec:
  forProvider:
    projectId: "your-project-id"
    dcName: europe-west3-c
    spec:
    - cloud:
      - openstack: []
      version: "1.34.1"
  providerConfigRef:
    name: default
```

### 2. Create Node Deployment

After the cluster is created, get its ID and create nodes:

```yaml
apiVersion: node.syseleven-metakube/v1alpha1
kind: Deployment
metadata:
  name: my-cluster-nodes
spec:
  forProvider:
    clusterId: "cluster-id-from-status"
    spec:
    - replicas: 2
      template:
      - cloud:
        - openstack:
          - flavor: m1.small
            image: "Ubuntu 22.04"
        operatingSystem:
        - ubuntu:
          - distUpgradeOnBoot: false
  providerConfigRef:
    name: default
```

## Available Parameters

### XRD Parameters

- `projectId` (required): Your MetaKube project ID
- `dcName`: Datacenter name (default: europe-west3-c)
- `kubernetesVersion`: Kubernetes version (default: 1.28.0)
- `cloudProvider`: Cloud provider - openstack, aws, or azure (default: openstack)
- `nodeCount`: Number of worker nodes (default: 2)
- `nodeSize`: Node flavor/size (default: m1.small)
- `nodeImage`: Node OS image (default: Ubuntu 22.04)

## Cleanup

```bash
# Delete the claim (this will delete the cluster and nodes)
kubectl delete -f examples/claim/metakube-cluster-claim.yaml

# Or delete directly
kubectl delete clusters.metakube.syseleven-metakube my-cluster
kubectl delete deployments.node.syseleven-metakube my-cluster-nodes
```

## Troubleshooting

### Check Provider Logs

```bash
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-syseleven-metakube -f
```

### Check Provider Status

```bash
kubectl get providers.pkg.crossplane.io
```

### Check ProviderConfig

```bash
kubectl describe providerconfig default
```

### Common Issues

1. **Invalid credentials**: Ensure your MetaKube token is correct in the secret
2. **Project ID not found**: Verify your project ID is correct
3. **Insufficient permissions**: Ensure your token has the necessary permissions in MetaKube
4. **Node deployment fails**: Check that the cluster ID is correctly propagated from the cluster resource
