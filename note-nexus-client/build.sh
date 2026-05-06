#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "Downloading .NET install script..."
wget https://dot.net/v1/dotnet-install.sh -O dotnet-install.sh
chmod +x ./dotnet-install.sh

echo "Installing .NET 9.0 SDK..."
./dotnet-install.sh --channel 9.0

echo "Configuring PATH..."
export DOTNET_ROOT=$HOME/.dotnet
export PATH=$PATH:$DOTNET_ROOT:$DOTNET_ROOT/tools

echo "Publishing Blazor WebAssembly application..."
dotnet publish -c Release
