#!/bin/bash

mkdir certs
rm certs/*

openssl req -new -nodes -x509 -out certs/server.crt -keyout certs/server.key -days 3650 -subj "/C=NO/ST=Some-State/L=City/O=ABC-Company/OU=IT/CN=www.abc.com/emailAddress=abc@abc.com"

