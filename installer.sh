#!/bin/bash
if [ "$(id -u)" != "0" ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

apt-get -y install python-pip python-flup python-mysqldb python-twisted python-imaging

sudo pip install Django==1.5.4
sudo pip install apscheduler
sudo pip install django-lockdown
