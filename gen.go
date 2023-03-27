package main

// API swagger
//nolint:lll
//go:generate swagger generate server --target api/server --name MassaWallet --spec ./api/wallet_api-v0.yml --principal interface{} --exclude-main
