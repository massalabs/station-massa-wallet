# Massa core plugin
This is a thyra plugin implementing Massa core features such as wallet.

## How to test ?

The first thing to do is to find the endpoint where the plugin listens.

When you start this plugin, which you can do using the command `go run cmd/massa-core/main.go` for example, you will get a message like the following in the terminal:
``shell
2022/11/21 22:11:43 Serving massa core at http://[::]:33049
warning: insecure HTTPS configuration.
	To fix this, use your own .crt and .key files using `--tls-certificate` and `--tls-key` flags
2022/11/21 22:11:43 Serving massa core at https://[::]:36869
```
You can find the port at the end of the first line. In this example, the port is 33049.
You can then access the service using the following URL: http://[::]:33049

### Signature

You can test the signature using the following curl command:
```shell
curl -X POST \
  'http://[::]:33049/rest/wallet/test/signOperation' \
  -H 'content-type: application/json' \
  -d '{"operation":"MjIzM3QyNHQ="}'
```
>**NOTE**: The port in the command must be modified according to your own instance of the plugin.
