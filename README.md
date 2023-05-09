[![codecov](https://codecov.io/gh/massalabs/thyra-plugin-wallet/branch/main/graph/badge.svg?token=RZ6AN1ISEA)](https://codecov.io/gh/massalabs/thyra-plugin-wallet)

# Massa Wallet Plugin

This is the MassaStation plugin that implements the Massa wallet features.

## Developer guide

This section helps developer getting started.

If you want to contribute, please refer to our [CONTRIBUTING](CONTRIBUTING.md) guide.

### Install dependencies

```shell
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### On Linux

```shell
apt install libgl1-mesa-dev xorg-dev gcc-mingw-w64-x86-64 gcc-mingw-w64
```

### Build

Generate the projects: go-swagger, wails, web-frontend:

```shell
go generate ./...
```

On macos:

```shell
./build_darwin.sh
```

This will create a binary file `wallet-plugin` in `build/wallet-plugin` folder.

### Test

```shell
go test ./...
```

### Run

For development purpose, you can run the plugin in standalone mode: it will not try to register with MassaStation.

```shell
cd web-frontend && npm i && npm run build:standalone && cd ..
STANDALONE=1 wails dev
```

The `STANDALONE` environment variable is to run the plugin without MassaStation.

Now navigate into <http://localhost:8080>. Note that some features will not work if
[thyra-server](https://github.com/massalabs/thyra) is not running.

**Install manually the plugin:**

For development purpose, you can install the plugin manually:

```shell
./manual-install.sh
```

This will create MassaStation plugin directories and copy the binary file created in the previous step so that
MassaStation can detect the plugin and launch it.

### Postman collection

You will find a postman collection in the `/api` directory.

Before testing this API, you must initialize the `baseURL` variable to <127.0.0.1:8080>.
