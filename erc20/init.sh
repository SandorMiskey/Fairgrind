#! /bin/bash

VENV="$PWD/venv"

# Create virtual python3 environment and activate it

python3 -m venv $VENV
source $VENV/bin/activate

# Install packages

$VENV/bin/pip3 install -U pip

# $VENV/bin/pip3 install urllib3
# $VENV/bin/pip3 install mysql.connector
# $VENV/bin/pip3 install mysql-connector-python --no-cache-dir
$VENV/bin/pip3 install pyyaml --no-cache-dir
$VENV/bin/pip3 install web3 --no-cache-dir
