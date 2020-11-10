# Chainsource Gateway
The gateway component for the Digital Bill of Materials

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [How to Use](#how-to-use)
  - [API](#api)
  - [Configuration](#configuration)
- [Helm Deployment](#helm-deployment)
- [Getting Help](#getting-help)
- [Getting Involved](#getting-involved)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## How to Use

### API

Latest OpenAPI Specifications and Postman Collection Files for this API is available on the [api-specs repository](https://github.com/DBOMproject/deployment/blob/master/api-specs/gateway)

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

Instructions for deploying the Chainsource Gateway using helm charts can be found [here](https://github.com/DBOMproject/deployment/blob/master/charts/chainsource-gateway)

## Getting Help

If you have any queries on insert-project-name, feel free to reach us on any of our [communication channels](https://github.com/DBOMproject/community/blob/master/COMMUNICATION.md) 

If you have questions, concerns, bug reports, etc, please file an issue in this repository's [issue tracker](https://github.com/DBOMproject/node-sdk/issues).

## Getting Involved

This section should detail why people should get involved and describe key areas you are
currently focusing on; e.g., trying to get feedback on features, fixing certain bugs, building
important pieces, etc.

General instructions on _how_ to contribute should be stated with a link to [CONTRIBUTING](CONTRIBUTING.md).