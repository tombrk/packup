#!/bin/sh

rclone sync --use-json-log -v data: /mnt > /dev/stderr

echo "/mnt"
