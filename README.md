[![codecov](https://codecov.io/gh/massalabs/thyra-plugin-wallet/branch/main/graph/badge.svg?token=RZ6AN1ISEA)](https://codecov.io/gh/massalabs/thyra-plugin-wallet)

# Massa Wallet Plugin

This is the MassaStation plugin that implements the Massa wallet features.

## Developer guide

This section helps developer getting started.

If you want to contribute, please refer to our [CONTRIBUTING](CONTRIBUTING.md) guide.

### Install Task

Follow the installation instructions here:
[task-install](https://taskfile.dev/installation/)

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

This will create a binary file `wallet-plugin{.exe}` in `build/wallet-plugin` folder.

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
[MassaStation-server](https://github.com/massalabs/thyra) is not running.

**Install manually the plugin for Massa Station:**

For development purpose, you can install the plugin manually:

```shell
task install-plugin
```

This will create MassaStation plugin directories and copy the binary file created in the previous step so that
MassaStation can detect the plugin and launch it.

### Postman collection

You will find a postman collection in the `/api` directory.

Before testing this API, you must initialize the `baseURL` variable to <127.0.0.1:8080>.
