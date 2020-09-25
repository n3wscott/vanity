# vanity Go Path Router

This configures an Knative Service which will act as a router for `go` cli, this means you don't have to reference all your `go` packages using the full `github.com/{org}/{project}` path, you'll be able to use your own domain like the Kubernetes (k8s.io, sigs.k8s.io), Etcd (etcd.io), etc.

## Why

There are a couple reasons why this might be useful for you. 

1. If you plan on moving your project around but want to be independent of `github` domains and orgs. 
2. Shorten the package imports for cleanliness. 
3. You like customizing things :smile:

## Setup

### Prerequisites
    
- [`ko`](https://github.com/google/ko).

### YAML config



## Install 

## Thanks

This was based on code from [christopherhein/go-path-router](https://github.com/christopherhein/go-path-router) and [GoogleCloudPlatform/govanityurls](https://github.com/GoogleCloudPlatform/govanityurls)