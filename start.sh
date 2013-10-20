#!/bin/bash


python log.py &
python ./web/ircweb/manage.py runserver 0.0.0.0:8000 &
