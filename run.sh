#!/bin/bash

cd src
echo "Running go fmt"
go fmt
 
cd ..


echo "Tidy up go mods"
go mod tidy

go run src/main.go
