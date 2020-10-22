# Vanity, vanity domain mapping for go-imports

This configures an Knative Service which will act as a router for `go` cli, this
means you don't have to reference all your `go` packages using the full
`github.com/{org}/{project}` path, you'll be able to use your own domain like
the Kubernetes (k8s.io, sigs.k8s.io), Etcd (etcd.io), etc.

## Why

There are a couple reasons why this might be useful for you.

1. If you plan on moving your project around but want to be independent of
   `github` domains and orgs.
2. Shorten the package imports for cleanliness.
3. You like customizing things :smile:

## Setup

The following assumes Kubernetes and Knative, but with a little work you can
modify the deployment to be anything you need.

### Prerequisites

- [`ko`](https://github.com/google/ko).
- If targeting Kubernetes, [Knative Serving](https://knative.dev/docs/serving)

### YAML config

A simple [example](./kodata/example.yaml) config is provided in the
[`kodata`](https://github.com/google/ko#including-static-assets) repo.

```yaml
host: example.com
paths:
  /foo:
    repo: https://github.com/example/foo

  /bar:
    repo: https://github.com/example/bar
```

Full options are:

```yaml
host: <vanity host>                                           # required
cache_max_age: <max cache age for http responses, in seconds> # optional, defaults to 24 hours
paths:                                                        # required
  <path relitive to vanity host>:                             # required
    repo: <repo location (without .git)>                      # required
    display: <go-import config for go-source metadata>        # optional
    vcs: <repo kind, one of: [git, bzr, git, hg, svn]>        # optional
```

## Running Locally

You can test out your config by running,

```shell script
go run .
```

And then you can poke around with `curl http://localhost:8080` or
`curl http://localhost:8080/foo`.

## Running on Knative Serving,

Edit the `service.yaml` file to suite your own needs and configuration. Then
publish it.

```shell script
ko apply -f service.yaml
```

Tie a [domain](https://knative.dev/docs/) to the service, and you are good to
go!

### Running on Cloud Run

`ko` lets you publish containers from go paths directly, so you can do something
like the following:

```shell script
ko publish github.com/n3wscott/vanity
```

Use the resulting image and set `VANITY_CONFIG` in the env vars to be your
config, set the vanity domain on the service, and you are good to go!

Or (experimentally) via the `gcloud` cli:

```shell script
gcloud beta run services replace service.yaml
```

## Thanks

This was based on code from
[christopherhein/go-path-router](https://github.com/christopherhein/go-path-router)
and
[GoogleCloudPlatform/govanityurls](https://github.com/GoogleCloudPlatform/govanityurls)
