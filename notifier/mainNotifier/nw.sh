#!/bin/bash
ADDR="127.0.0.1:10101"
$*
curl -d '{"title":"Done!","message":"${*}"}' http://${ADDR}
