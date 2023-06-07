#!/bin/bash -ex

PLUGIN=wallet-plugin

mkdir -p /usr/local/share/massastation/plugins/wallet-plugin
cp build/$PLUGIN/wallet-plugin /usr/local/share/massastation/plugins/wallet-plugin
cp manifest.json /usr/local/share/massastation/plugins/wallet-plugin
