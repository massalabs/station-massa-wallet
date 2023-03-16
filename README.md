# Massa wallet plugin

This is the Thyra plugin that implements the Massa wallet features.

## Developer guide

This section helps developer getting started.

If you want to contribute, please refer to our [CONTRIBUTING](CONTRIBUTING.md) guide.

### Generate

```shell
go generate ./...
```

This will generate go swagger file and a javascript file with constants variable for the frontend.

### Run

For development purpose, you can run the plugin in standalone mode: it will not try to register against Thyra.

```shell
STANDALONE=1 go run cmd/massa-wallet/thyra-plugin-wallet.go
```

The `STANDALONE` environment variable is to run the plugin without Thyra.

Now navigate into <http://localhost:8080>. Note that some feature will not work if
[thyra-server](https://github.com/massalabs/thyra) is not running.

### Build

```shell
./build.sh
```

This will create a binary file `thyra-plugin-wallet` in `build/wallet-plugin` folder.

**Install manually the plugin:**

For development purpose, you can install the plugin manually:

```shell
./manual-install.sh
```

This will create Thyra plugin directories and move the binary file created in the previous step so that
Thyra can detect the plugin and launch it.

### Postman collection

You will find a postman collection in the `/api` directory.

Before testing this API, you must initialize the `baseURL` variable to <127.0.0.1:8080>.
