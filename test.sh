#!/usr/bin/env bash

json="{\"JWT\": \"beta\", \"JWT_type\": \"google\"}"
curl --request POST -H "Content-type: text/json" --data "$json" http://localhost:5000/api/login

sleep 5

json="{\"JWT\": \"xyz\", \"JWT_type\": \"fackbook\"}"
curl --request POST -H "Content-type: text/json" --data "$json" http://localhost:5000/api/login
