# Using Compositions to Simplify MetaKube Cluster Creation

This guide explains how to use Crossplane Compositions to abstract away repetitive configuration details when creating MetaKube clusters.

## The Problem

When creating clusters directly with the MetaKube provider, you need to repeat the project ID in multiple places:

```yaml
apiVersion: metakube.metakube.syseleven.io/v1alpha1
kind: Cluster
spec:
  forProvider:
    projectId: "YOUR_PROJECT_ID_HERE"  # 1st place
    dcName: syseleven-dus2
    spec:
    - cloud:
      - openstack:
        - applicationCredentials:
          - id: "s11auth:YOUR_PROJECT_ID_HERE"  # 2nd place (must match!)
            secretSecretRef:
              key: token
              name: metakube-creds
              namespace: crossplane-system
```

## The Solution: Compositions

With a Composition, users only specify the project ID **once**, and the Composition automatically:
1. Sets `spec.forProvider.projectId`
2. Constructs `applicationCredentials.id` as `"s11auth:PROJECT_ID"`
3. References the correct secret

## Setup

### 1. Install the XRD (Composite Resource Definition)

```bash
kubectl apply -f xrd/xrd-metakube-cluster.yaml
```

This defines a simplified API (`MetaKubeCluster`) that users can create.

### 2. Install the Composition

```bash
kubectl apply -f composition/composition-metakube-cluster.yaml
```

This tells Crossplane how to transform a `MetaKubeCluster` into actual MetaKube resources.

### 3. Install the Crossplane function-patch-and-transform

```bash
kubectl apply -f setup/function-patch-and-transform.yaml
```

This is the Crossplane function that performs the transformations.

## Usage

Now users can create clusters with a much simpler manifest:

```bash
kubectl apply -f claim/metakube-cluster-claim.yaml
```

```yaml
apiVersion: platform.example.com/v1alpha1
kind: MetaKubeCluster
metadata:
  name: my-test-cluster
  namespace: default
spec:
  parameters:
    projectId: "YOUR_PROJECT_ID_HERE"  # Only specify once!
    dcName: syseleven-dus2
    kubernetesVersion: "1.34.1"
    nodeCount: 2
```

## What Happens Behind the Scenes

The Composition automatically creates:

1. **A Cluster resource** with:
   - `projectId` set from the claim
   - `applicationCredentials.id` constructed as `"s11auth:YOUR_PROJECT_ID"`
   - Secret reference to `metakube-creds`

2. **A NodeDeployment resource** for worker nodes

## Benefits

✅ **No repetition**: Project ID specified only once
✅ **No manual string construction**: The `"s11auth:PROJECT_ID"` format is automatic
✅ **Simpler for users**: Hide provider-specific details
✅ **Consistent**: All clusters use the same pattern
✅ **Composable**: Can add more resources (monitoring, backups, etc.) to the composition

## Customization

You can customize the Composition to:
- Add default values for common parameters
- Include additional resources (monitoring, backups)
- Apply organization-specific policies
- Support multiple cloud providers (OpenStack, AWS, Azure)

## Verification

Check the created resources:

```bash
# Check the claim
kubectl get metakubecluster my-test-cluster

# Check the composite resource
kubectl get xmetakubecluster

# Check the actual MetaKube cluster
kubectl get clusters.metakube.metakube.syseleven.io

# Check all resources
kubectl get crossplane
```

## Troubleshooting

If the cluster doesn't appear:

```bash
# Check composite resource status
kubectl describe xmetakubecluster <name>

# Check events
kubectl get events --sort-by='.lastTimestamp'

# Check function logs
kubectl logs -n crossplane-system -l pkg.crossplane.io/function=function-patch-and-transform
```
