# Ketch Provider for Terraform

The Ketch provider for Terraform builds an application-level interface that enables developers to deploy and manage jobs and applications on Kubernetes using Terraform

# Building the provider

1. Build the Terraform provider

```
make install
```

2. Run the examples

```
cd example
terraform init && terraform apply --auto-approve   
```

For the example above to work, make sure you:
- Have Ketch installed in your cluster
- Have local kubectl access to the cluster 
- Updated the ingress controller configuration in the framework definition

# Detailed documentation

Detailed documentation on how to use this provider can be found here: https://learn.theketch.io/docs/terraform
