[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/sensu/sensu-victorops-handler)
![Go Test](https://github.com/sensu/sensu-victorops-handler/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/sensu/sensu-victorops-handler/workflows/goreleaser/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/senssensuu/sensu-victorops-handler)](https://goreportcard.com/report/github.com/sensu/sensu-victorops-handler)

# Sensu VictorOps Handler

## Table of Contents
- [Overview](#overview)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Resource definition](#resource-definition)
- [Usage examples](#usage-examples)
  - [Help output](#help-output)
  - [Templates](#templates)
  - [Environment variables](#environment-variables)
  - [Argument annotations](#argument-annotations)
- [Installation from source](#installation-from-source)

## Overview

The Sensu VictorOps Handler is a [Sensu Go Handler][6] for sending events to the
[VictorOps][11] incident response platform.

As of the initial version (0.1.x), this is meant to work in the same fashion as
the prior Ruby based plugin, [sensu-plugins-victorops][12], with the following
changes:
- The environment variables for routing key and API URL are now
SENSU_VICTOROPS_ROUTINGKEY and SENSU_VICTOROPS_APIURL, respectively
- Since Sensu Go events do not have an action, the RECOVERY message_type is
based on event.check.status == 0

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not
using an asset, please consider doing so! If you're using sensuctl 5.13 with
Sensu Backend 5.13 or later, you can use the following command to add the asset:

```
sensuctl asset add sensu/sensu-victorops-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the
[Bonsai Asset Index][13].

### Resource definition

```yml
---
type: Handler
api_version: core/v2
metadata:
  name: sensu-victorops-handler
  namespace: default
spec:
  command: sensu-victorops-handler
  filters:
  - is_incident
  - not_silenced
  type: pipe
  runtime_assets:
  - sensu/sensu-victorops-handler
  secrets:
  - name: SENSU_VICTOROPS_ROUTINGKEY
    secret: victorops_routingkey
```

## Usage examples

### Help output

```The Sensu Go VictorOps handler for sending notifications

Usage:
  sensu-victorops-handler [flags]
  sensu-victorops-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -a, --api-url string              The URL for the VictorOps API (default "https://alert.victorops.com/integrations/generic/20131114/alert")
  -r, --routingkey string           The VictorOps Routing Key
  -e, --entity-id-template string   The template for the Entity ID sent to VictorOps (default "{{.Entity.Name}}/{{.Check.Name}}")
  -m, --message-template string     The template for the message sent to VictorOps (default "{{.Entity.Name}}:{{.Check.Name}}:{{.Check.Output}}")
  -h, --help                        help for sensu-victorops-handler
```

### Templates

This handler provides options for using templates to populate the values
provided by the event in the message sent via SNS. More information on
template syntax and format can be found in [the documentation][14].

### Environment variables

|Argument     |Environment Variable       |
|-------------|---------------------------|
|--routingkey |SENSU_VICTOROPS_ROUTINGKEY |
|--api-url    |SENSU_VICTOROPS_APIURL     |

**Security Note:** Care should be taken to not expose the routing key for this handler by specifying it
on the command line or by directly setting the environment variable in the handler definition.  It is
suggested to make use of [secrets management][17] to surface it as an environment variable.  The
handler definition above references it as a secret.  Below is an example secrets definition that make
use of the built-in [env secrets provider][18].

```yml
---
type: Secret
api_version: secrets/v1
metadata:
  name: victorops_routingkey
spec:
  provider: env
  id: SENSU_VICTOROPS_ROUTINGKEY
```

### Argument annotations

All arguments for this handler are tunable on a per entity or check basis based on annotations.  The
annotations keyspace for this handler is `sensu.io/plugins/victorops/config`.

**NOTE**: Due to [check token substituion][15], supplying a template value such
as for `message-template` as a check annotation requires that you place the
desired template as a [golang string literal][16] (enlcosed in backticks)
within another template definition.  This does not apply to entity annotations.

#### Examples

To change the message template for a particular check, for that check's metadata add the following:

```yml
type: CheckConfig
api_version: core/v2
metadata:
  annotations:
    sensu.io/plugins/victorops/config/message-template: "{{`{{.Entity.Name}}/{{.Check.Name}}: {{.Check.State}}, {{.Check.Occurrences}}`}}"
[...]
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an
Asset. If you would like to compile and install the plugin from source or
contribute to it, download the latest version or create an executable from
this source.

From the local path of the sensu-victorops-handler repository:

```
go build
```

[2]: https://github.com/sensu-community/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu-community/handler-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu-community/handler-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/handlers/
[7]: https://github.com/sensu-community/handler-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu-community/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[11]: https://victorops.com/
[12]: https://github.com/sensu-plugins/sensu-plugins-victorops
[13]: https://bonsai.sensu.io/assets/sensu/sensu-victorops-handler
[14]: https://docs.sensu.io/sensu-go/latest/observability-pipeline/observe-process/handler-templates/
[15]: https://docs.sensu.io/sensu-go/latest/observability-pipeline/observe-schedule/checks/#check-token-substitution
[16]: https://golang.org/ref/spec#String_literals
[17]: https://docs.sensu.io/sensu-go/latest/guides/secrets-management/
[18]: https://docs.sensu.io/sensu-go/latest/guides/secrets-management/#use-env-for-secrets-management
