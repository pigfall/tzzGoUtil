#!/bin/bash
${*}
curl -d "{\"title\":\"Done!\",\"message\":\"${*}\"}" http://192.168.43.194:10101
