from django.conf.urls import patterns, url

from irc import views

urlpatterns = patterns('',
    url(r'^$', views.index, name='index'),
    url(r'^(?P<category>[-A-Za-z0-9_]+)/$', views.index, name='type'),
    url(r'^(?P<category>[-A-Za-z0-9_]+)/(?P<page>\d+)/$', views.index, name='tpye_with_link')
)
