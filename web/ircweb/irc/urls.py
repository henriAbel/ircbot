from django.conf.urls import patterns, url

from irc import views

urlpatterns = patterns('',
    url(r'^$', views.index, name='index')
)
