
# Contributing

We welcome and appreciate contributions from the community. Whether it's a bug report, feature request,
code contribution, or anything else, we're grateful for your effort. Here's a brief guide to get you started.

Any contribution you make will be reflected on [CONTRIBUTORS](CONTRIBUTORS.md).

To get an overview of the project, read the [README](README.md).

## Product Vision

This Plugin Wallet is MassaLabs's native wallet. Our goal is to provide maximum security to our end-users while still having a simple, enjoyable and user friendly experience. A wallet with no compromise.

Definition of a Massa account: it is composed of 2 distinct elements:

* A Massa account that is a smart contract
* Cryptographic pair of keys protecting this account

A wallet is the application that allows you to interact with your Massa account coupled with its key pair.

## What you can Contribute

We welcome contributions of all kinds! Here are some ways you can help:

* [Report bugs](#bug-reports) and [suggest features](#feature-requests)
* Improve the documentation
* [Write code, such as bug fixes and new features](#code-contributions)
* Submit design improvements

## Code of Conduct

Help us keep this project open and inclusive and keep our community approachable and respectable.

## Bug Reports

If you find any bugs, feel free to open an issue in the
[issue tracker](https://github.com/massalabs/station-massa-wallet/issues). Please include as much information as possible,
such as operating system, browser version, etc.

## Feature Requests

If you have an idea for a new feature, feel free to open an issue in the
[issue tracker](https://github.com/massalabs/station-massa-wallet/issues). Please include as much information as possible,
such as what the feature would do, why it would be useful, etc.

## Code Contributions

We are always looking for developers to help us improve the codebase. If you are interested in contributing code, please
take a look at the [open issues](https://github.com/massalabs/station-massa-wallet/issues) and follow the steps below:

1. [Fork](https://help.github.com/en/github/getting-started-with-github/fork-a-repo) the repository
2. Create a feature branch off of the `main` branch
3. Make your changes
4. Submit a pull request

When contributing:

* Ensure that all new features or changes are aligned with our product vision.
* Write clean, easy-to-understand code.
* Write tests for your code in order to ensure that it works as expected.
* Make sure all tests pass before submitting a pull request.
* Follow our code conventions and formatting standards.
* Document your code thoroughly.

## Questions

If you have any questions, feel free to open an issue in the
[issue tracker](https://github.com/massalabs/station-massa-wallet/issues).

*Thank you for taking the time to contribute to our project! We look forward to seeing your contributions!*

# Massa concepts

## Operation vs Transaction

In the Massa blockchain, an operation refers to an action that is included in blocks. Operations can take various forms, such as buying a roll, selling a roll, transferring funds, executing a smart contract, or deploying a smart contract. Within the wallet, we use the term **transaction** to describe these operations.

## Transaction vs Transfer

In the Massa blockchain, a transaction is a specific type of operation that enables the transfer of MAS coins from one account to another. Within the wallet, we specifically refer to this type of operation as a "transfer".

## Swagger file

The swagger file `api/walletApi-V0.yml` is the source of truth for the wallet API.
It is used to generate the API documentation and the API client.

Some endpoint responses are code 422 but the code itself never return such status code.
Go-swagger will internally generate 422 response when the inputs are invalid.

## Massa account file

This software implements this standard: <https://github.com/massalabs/massa-standards/blob/main/wallet/file-format.md>.

For account file extension, we recommend `.yaml`. However, `.yml` is also supported in the import feature.

Starting from version v0.2.1, we will implement retro-compatibility for all breaking changes in the account standard above.

## Frontend

### Web app

Go swagger URL prefix: `/web-app`

Folder: `web-frontend`.

This folder contains the source code of the web frontend app. It is a React app.

### Wails

Wails is a native GUI library for Go. It allows us to build a native GUI app with Go and React.

Folder: `wails-frontend`.

This folder contains the source code of the wails frontend app. It is a React app.

### Error handling

* Use sentinel error
* Functions don't return a middleware responder, but an error
* for http error response, use as much as possible `internal/handler/wallet/newErrorResponse`
* for wails error code: use `pkg/utils/WailsErrorCode`
