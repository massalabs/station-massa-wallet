# Massa wallet plugin
This is a thyra plugin implementing Massa wallet features.

## How to test ?

The first thing to do is to find the endpoint where the plugin listens.

When you start this plugin, which you can do using the command `go run cmd/massa-wallet/main.go` for example, you will get a message like the following in the terminal:
``shell
2022/11/21 22:11:43 Serving massa core at http://[::]:33049
```
You can find the port at the end of the first line. In this example, the port is 33049.
You can then access the service using the following URL: http://[::]:33049

### Postman collection

You will find a postman collection in the /api directory.

Only one variable is to set to test this API:
- baseURL: to be replaced with 127.0.0.1:<port>
