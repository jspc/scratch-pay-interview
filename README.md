# Cloud Architect Challenge

The master version of this application may be found at: http://34.67.159.16:3333/

This repo will build:

1. A new network in GCP
1. Create a firewall granting access to port 3333 (See below for ssh/ server access)
1. Create a VM running a simple echo-server, fronted by nginx
1. Create a runscope test to continually monitor the service for uptime

These services run with docker containers, which are built from this repo, and run on a [shielded](https://cloud.google.com/shielded-vm/) VM running Google's [Container Optimized OS](https://cloud.google.com/container-optimized-os/docs/) (which seems to be managed CoreOS anyway).

The containers are segregated into their own docker network, and log/ push metrics into StackDriver.

## Building

You will need the following things:

1. [`terraform`](https://www.terraform.io)
1. A Google Cloud Platform account, alongwith a [Service Account and set of creds](https://cloud.google.com/docs/authentication/getting-started)
1. A runscope account, and [API key](https://www.runscope.com/applications/create)

### Environment variables

```bash
$ export GOOGLE_CLOUD_KEYFILE_JSON=~/path/to/gcp/creds.json
$ export RUNSCOPE_ACCESS_TOKE=your-runscope-application-key
$ export TF_VAR_runscope_team_id=your-runscope-team-id
```

### Running `terraform`

First, initialise the project

`$ terraform init`

Then run plan to ensure that:

1. Your API access works as expected; and
1. What you expect to change, is being changed


`$ terraform plan`

Should the output be as expected (the command worked, and the resources we're about to create look like the resources your expect to be created)

`$ terraform apply`

On success, this will give you the address of your service.

## On SSH and Immutable Infrastructure

Instances created via this repo have no SSH access at all. Services are sandboxed in containers, though of course breaking out of containers is not impossible.

This is all by design; there is no good reason for infrasturcture in production systems to have ssh access in the modern world. In a world of custom VM images, log and metric shipping, and autoscaling, instances are little more than appliances.

While it may be useful to have access to instances when developing automations and installations, this should never be the norm, and there should be every precaution taken to ensure such exceptions never exceed development accounts. In fact: to be able to login to a running instance should be difficult.

Ideally I would have proven some kind of audit in this project, but it would be out of the scope of this service deployment.


## Services

Services are run on the VM via a `cloud-init` script, which may be found at [script/cloud-init.yml](script/cloud-init.yml). This config creates a series of `systemd` units to run docker operations, and removes SSH access (should an actor be able to breach the google firewalls).

### echo-service

This services lives in the [app/](app/) directory of this repo. It is a go app, with tests, which returns information about the request made.

Tests and build steps can be run via the included [Makefile](app/Makefile)

```bash
$ curl -XPOST -H 'Content-Type: application/json' -d '{"hello":"world!"}' http://34.67.159.16:3333/ | jq '.'
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   413  100   395  100    18    879     40 --:--:-- --:--:-- --:--:--   919
{
  "method": "POST",
  "uri": "http://34.67.159.16:3333/",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "Content-Length": [
      "18"
    ],
    "Content-Type": [
      "application/json"
    ],
    "Host": [
      "34.67.159.16:3333"
    ],
    "User-Agent": [
      "curl/7.66.0"
    ]
  },
  "connection": {
    "local": {
      "network": "tcp",
      "address": "172.18.0.3:8080"
    },
    "remote": {
      "network": "tcp",
      "address": "221.146.26.223:34548"
    }
  },
  "raw_payload": "{\"hello\":\"world!\"}",
  "payload": {
    "hello": "world!"
  }
}
```

### nginx

Configuration for this container image may be found in [nginx/](nginx/). We provide a simple server config, which creates a listener and a proxy, and bundle it into the latest _stable_ nginx tag.
