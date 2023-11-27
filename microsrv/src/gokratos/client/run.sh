#/bin/bash

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

naming_dir=$SCRIPT_DIR/cache/naming
log_dir=$SCRIPT_DIR/log
rm -rf $naming_dir/* $log_dir/*
# find $naming_dir ! -name 'info.json' -type f -exec rm {} +
cd $SCRIPT_DIR
kratos run