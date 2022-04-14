#!/bin/sh

pg_dump $PGDATABASE > /mnt/db.sql

echo "/mnt"
