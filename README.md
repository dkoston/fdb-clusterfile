# fdb-clusterfile

Builds a foundationdb fdb.cluster file from hostnames of cluster coordinators

## Purpose

FoundationDB doesn't allow hostnames in cluster files. In order to make a cluster,
one doesn't always know IPs ahead of time (i.e. in k8s). Use this tool to 
generate fdb.cluster on-the-fly with hostnames (k8s services) to easily create
a cluster when using automated tools like Helm or Terraform.

## Example

fdb-clusterfile --fdb_addr fdb:fdb@localhost:4500:tcp --cluster_file fdb.cluster

Output: `fdb:fdb@127.0.0.1:tcp`

File Location: `./fdb.cluster`

## Development

### Pre-requisites

- go 1.11.2+
- docker

### Releases

Releases can be built (via docker) with `./release.sh`