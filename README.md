![docker-ci](https://github.com/DBOMproject/chainsource-gateway/workflows/docker-ci/badge.svg)
![GitHub](https://img.shields.io/github/license/dbomproject/chainsource-gateway)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/dbomproject/chainsource-gateway)
![Docker Pulls](https://img.shields.io/docker/pulls/dbomproject/chainsource-gateway)
[![Coverage Status](https://coveralls.io/repos/github/DBOMproject/chainsource-gateway/badge.svg?branch=master)](https://coveralls.io/github/DBOMproject/chainsource-gateway?branch=master)

# Chainsource Gateway
The gateway component for the Digital Bill of Materials

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [How to Use](#how-to-use)
  - [API](#api)
  - [Configuration](#configuration)
- [Helm Deployment](#helm-deployment)
- [Platform Support](#platform-support)
- [Getting Help](#getting-help)
- [Getting Involved](#getting-involved)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## How to Use

### API

Latest OpenAPI Specifications and Postman Collection Files for this API is available on the [api-specs repository](https://github.com/DBOMproject/api-specs/tree/master/gateway)

### Configuration

| Environment Variable         | Default               | Description                                 |
|------------------------------|-----------------------|---------------------------------------------|
| LOG_LEVEL                    | `info`                | The verbosity of the logging                |
| PORT                         | `3000`                | Port on which the gateway listens           |
| JAEGER_ENABLED               | `false`               | Is jaeger tracing enabled                   |
| JAEGER_HOST                  | ``                    | The jaeger host to send traces to           |
| JAEGER_SAMPLER_PARAM         | `1`                   | The parameter to pass to the jaeger sampler |
| JAEGER_SAMPLER_TYPE          | `const`               | The jaeger sampler type to use              |
| JAEGER_SERVICE_NAME          | `Chainsource Gateway` | The name of the service passed to jaeger    |
| JAEGER_AGENT_SIDECAR_ENABLED | `false`               | Is jaeger agent sidecar injection enabled   |

Configure `agent-config.yaml` with the details of your agent(s)

## Helm Deployment

Instructions for deploying the Chainsource Gateway using helm charts can be found [here](https://github.com/DBOMproject/deployments/tree/master/charts/chainsource-gateway)

## Platform Support

Currently, we provide pre-built container images for linux amd64 and arm64 architectures via our Github Actions Pipeline. Find the images [here](https://hub.docker.com/r/dbomproject/chainsource-gateway)

## Getting Help

If you have any queries on chainsource-gateway, feel free to reach us on any of our [communication channels](https://github.com/DBOMproject/community/blob/master/COMMUNICATION.md) 

If you have questions, concerns, bug reports, etc, please file an issue in this repository's [issue tracker](https://github.com/DBOMproject/chainsource-gateway/issues).

## Getting Involved

Find the instructions on how you can contribute in [CONTRIBUTING](CONTRIBUTING.md).
