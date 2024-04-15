#!/bin/sh

mysqldump \
    -h ${MYSQL_HOST} \
    -u${MYSQL_USER} \
    -p${MYSQL_PASSWORD} \
    ${MYSQL_DATABASE} > /mnt/${MYSQL_DATABASE}.sql

echo "/mnt"
