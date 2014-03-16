#!/bin/bash
if [[ $(id -u) -ne 0 ]] ; then echo "Please run as root" ; exit 1 ; fi

apt-get -y install python-pip python-flup python-twisted python-imaging

sudo pip install Django==1.5.4
sudo pip install apscheduler
sudo pip install django-lockdown
