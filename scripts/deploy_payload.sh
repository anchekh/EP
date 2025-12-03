#!/bin/bash

curl -X POST http://127.0.0.1:8080/deploy -d '{"service":"payload","replicas":3}'