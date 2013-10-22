#!/bin/bash

if [[ $EUID -ne 0 ]]; then
   echo "This script must be run as root" 
   exit 1
fi

apt-get -y install python-pip python-flup python-mysqldb python-twisted

pip install Django==1.5.4