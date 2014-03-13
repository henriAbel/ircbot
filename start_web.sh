#!/bin/bash
python ./web/ircweb/manage.py syncdb --noinput
python ./web/ircweb/manage.py runfcgi host=127.0.0.1 port=8080
