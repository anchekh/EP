#!/bin/bash

cd $(dirname "$0")/..
go build -o agent ./agent
./agent --controller http://127.0.0.1:8080