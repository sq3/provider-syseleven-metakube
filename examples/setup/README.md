# Setup Guide

This guide explains how to configure the SysEleven MetaKube provider for Crossplane.

## Understanding Credentials

The provider uses a **single set of credentials** that contains:

### MetaKube Credentials
- **token**: Your OpenStack application credential secret (this is also used as the MetaKube API token)
- **project-id**: Your MetaKube project ID

These credentials are used for:
1. Authenticating with the MetaKube API (provider-level)
2. Provisioning infrastructure via OpenStack (cluster-level)

## Setup Steps

### Step 1: Create the MetaKube Credentials Secret

This secret contains both your token and project ID.

```bash
# Replace with your actual credentials
kubectl create secret generic metakube-creds \
  --from-literal=credentials='{"token":"your-openstack-app-credential-secret","project-id":"your-project-id"}' \
  -n crossplane-system
```

**Note**: The `credentials` key contains JSON with both token and project-id.

Or apply the example file after editing:
```bash
# Edit the file first to add your credentials
kubectl apply -f metakube-api-secret.yaml
```

### Step 2: Create the ProviderConfig

This tells the provider where to find the API credentials.

```bash
kubectl apply -f ../cluster/providerconfig/providerconfig.yaml
```

### Step 3: Create a Cluster

Now you can create a MetaKube cluster:

```bash
kubectl apply -f ../test-cluster.yaml
```

## Verification

Check that everything is working:

```bash
# Check provider is healthy
kubectl get providers

# Check ProviderConfig exists
kubectl get providerconfigs.metakube.syseleven.io

# Check cluster status
kubectl get clusters.metakube.metakube.syseleven.io

# View cluster details
kubectl describe cluster cluster-via-crossplane
```

## Troubleshooting

### Check provider logs
```bash
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-syseleven-metakube -f
```

### Common issues

1. **"no providerConfigRef provided"**: Make sure your Cluster resource has a `providerConfigRef` field
2. **"cannot get referenced ProviderConfig"**: The ProviderConfig doesn't exist or has the wrong name
3. **"cannot extract credentials"**: The secret referenced in ProviderConfig doesn't exist or has wrong format. Ensure it contains both `token` and `project-id` fields.
4. **Authentication errors**: Check that your OpenStack application credential secret (token) is valid
