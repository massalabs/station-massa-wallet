#!/bin/bash -ex

PLUGIN=wallet-plugin

mkdir -p ~/.config/thyra/plugins/wallet-plugin
cp build/$PLUGIN/wallet-plugin ~/.config/thyra/plugins/wallet-plugin
cp manifest.json ~/.config/thyra/plugins/wallet-plugin
cp web/html/wallet.svg ~/.config/thyra/plugins/wallet-plugin
