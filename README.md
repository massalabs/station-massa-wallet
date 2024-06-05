[![codecov](https://codecov.io/gh/massalabs/station-massa-wallet/branch/main/graph/badge.svg?token=RZ6AN1ISEA)](https://codecov.io/gh/massalabs/station-massa-wallet)

# Massa Wallet Plugin

This is the MassaStation plugin that implements the Massa wallet features.

## Developer guide

This section helps developer getting started.

If you want to contribute, please refer to our [CONTRIBUTING](CONTRIBUTING.md) guide.

### Install Task

Follow the installation instructions here:
[task-install](https://taskfile.dev/installation/)

On Windows, we recommend to run `go install github.com/go-task/task/v3/cmd/task@latest` and yo use task commands in a git bash terminal.

### Install dependencies

```shell
task install
```

### Build

Generate the projects: go-swagger, wails, web-frontend:

```shell
task generate
```

```shell
task build
```

### Test

```shell
task test
```

### Run

For development purpose, you can run the plugin in standalone mode: it will not try to register with MassaStation.

```shell
task run
```

All in one build & run:

```shell
task build-run
```

The `STANDALONE` environment variable is to run the plugin without MassaStation.

Now navigate into <http://localhost:8080>. Note that some features will not work if
[MassaStation-server](https://github.com/massalabs/station) is not running.

**Install manually the plugin for Massa Station:**

For development purpose, you can install the plugin manually:

```shell
task install-plugin
```

This will create MassaStation plugin directories and copy the binary file created in the previous step so that
MassaStation can detect the plugin and launch it.

### Postman collection

You can import the swagger file `api/walletApi-V0.yml` into Postman to test the API.
