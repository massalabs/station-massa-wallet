[![codecov](https://codecov.io/gh/massalabs/thyra-plugin-wallet/branch/main/graph/badge.svg?token=RZ6AN1ISEA)](https://codecov.io/gh/massalabs/thyra-plugin-wallet)

# Massa Wallet Plugin

This is the Thyra plugin that implements the Massa wallet features.

## Developer guide

This section helps developer getting started.

If you want to contribute, please refer to our [CONTRIBUTING](CONTRIBUTING.md) guide.

### Install dependencies

```shell
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
```

#### On Linux

```shell
apt install libgl1-mesa-dev xorg-dev gcc-mingw-w64-x86-64 gcc-mingw-w64
```

### Build

```shell
./build.sh
```

This will create a binary file `thyra-plugin-wallet` in `build/wallet-plugin` folder.

### Run

For development purpose, you can run the plugin in standalone mode: it will not try to register against Thyra.

```shell
STANDALONE=1 go run cmd/massa-wallet/thyra-plugin-wallet.go
```

The `STANDALONE` environment variable is to run the plugin without Thyra.

Now navigate into <http://localhost:8080>. Note that some features will not work if
[MassaStation-server](https://github.com/massalabs/thyra) is not running.

**Install manually the plugin:**

For development purpose, you can install the plugin manually:

```shell
./manual-install.sh
```

This will create Thyra plugin directories and copy the binary file created in the previous step so that
Thyra can detect the plugin and launch it.

### Postman collection

You will find a postman collection in the `/api` directory.

Before testing this API, you must initialize the `baseURL` variable to <127.0.0.1:8080>.
