# Sensu VictorOps Handler

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Resource definition](#resource-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The Sensu VictorOps Handler is a [Sensu Handler][6] for sending events to the
[VictorOps][11] incident response platform.

## Files

N/A

## Usage examples

### Help

```The Sensu Go VictorOps handler for sending notifications

Usage:
  sensu-victorops-handler [flags]
  sensu-victorops-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -a, --api-url string      The URL for the VictorOps API (default "https://alert.victorops.com/integrations/generic/20131114/alert")
  -h, --help                help for sensu-victorops-handler
  -r, --routingkey string   The VictorOps Routing Key
```
### Environment Variables and Annotations

|Environment Variable|Setting|Annotation|
|--------------------|-------|----------|
|SENSU_VICTOROPS_ROUTINGKEY| same as -r / --routingkey|sensu.io/plugins/victorops/config/routingkey|
|SENSU_VICTOROPS_APIURL|same as -a / --api-url|sensu.io/plugins/victorops/config/api-url|

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not
using an asset, please consider doing so! If you're using sensuctl 5.13 with
Sensu Backend 5.13 or later, you can use the following command to add the asset:

```
sensuctl asset add nixwiz/sensu-victorops-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the
[Bonsai Asset Index][https://bonsai.sensu.io/assets/nixwiz/sensu-victorops-handler].

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
  - nixwiz/sensu-victorops-handler
  secrets:
  - name: SENSU_VICTOROPS_ROUTINGKEY
    secret: victorops-routingkey
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an
Asset. If you would like to compile and install the plugin from source or
contribute to it, download the latest version or create an executable script
from this source.

From the local path of the sensu-victorops-handler repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
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
