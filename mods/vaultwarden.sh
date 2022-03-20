#!/bin/sh

# cold-store sqlite db
sqlite3 /mnt/data/db.sqlite3 ".backup '/mnt/bak.sqlite3'" 1>&2

# print output dir
echo /mnt
