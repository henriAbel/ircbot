#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

apt-get -y install python-twisted
apt-get -y install python-mysqldb

python ./database_creator.py
