#!/bin/ash
apk --no-cache add npm go-bindata
cd chrono && ./prepare.sh && cd ..
