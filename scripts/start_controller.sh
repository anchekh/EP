#!/bin/bash

cd $(dirname "$0")/..
go build -o controller ./controller
./controller