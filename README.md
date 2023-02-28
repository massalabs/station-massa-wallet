# Massa wallet plugin

This is the Thyra plugin that implements the Massa wallet features.

## How to test it?

The first thing to do is to find the endpoint that the plugin listens to.

When you start this plugin, which you can do using the command `go run cmd/massa-wallet/main.go` for example, you will get a message like the following in the terminal:

```shell
2022/11/21 22:11:43 Serving massa wallet at http://[::]:33049
```

The listening port is at the end of the first line. In this example, the port is `33049`.

You can then access the service using the following URL: <http://127.0.0.1:33049>

### Postman collection

You will find a postman collection in the `/api` directory.

Before testing this API, you must initialize the `baseURL` variable to 127.0.0.1:`port`, port being the listening port that the plugin listens to (`33049` in the previous example).

## Contribute

Generate:

```shell
go generate ./...
```

Run:

```shell
go run cmd/massa-wallet/thyra-plugin-wallet.go 1 --standalone
```

Build:

```shell
./build.sh
```

Install the plugin:

```shell
mkdir ~/.config/thyra/my_plugins/thyra-plugin-wallet
PLUGIN=wallet-plugin
mv build/$PLUGIN/thyra-plugin-wallet ~/.config/thyra/my_plugins/thyra-plugin-wallet
```
