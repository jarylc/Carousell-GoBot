#!/bin/ash

set -e
[[ "$DEBUG" == "true" ]] && set -x

cd /app

getent group cgb >/dev/null || addgroup -g ${GID} cgb
getent passwd cgb >/dev/null || adduser -h /data -s /bin/sh -G cgb -D -u ${UID} cgb
chmod o+rx /data
chown -R cgb:cgb /data

exec su-exec cgb:cgb $@
