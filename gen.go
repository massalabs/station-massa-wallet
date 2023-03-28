package main

// API swagger
//nolint:lll
//go:generate swagger generate server --target api/server --name MassaWallet --spec ./api/walletApi-V0.yml --principal interface{} --exclude-main
