from django.http import HttpResponse
from django.template import RequestContext, loader
from irc.models import Users, Link, Message, Linktype

def index(request):
	# TODO, ask once from sql
    links = Link.objects.filter(type = Linktype.NORMAL)
    youtube = Link.objects.filter(type = Linktype.YOUTUBE)
    picture = Link.objects.filter(type = Linktype.PICTURE)
    template = loader.get_template("irc/index.html")
    context = RequestContext(request, {
    	'links': links,
    	'youtube': youtube,
    	'picture': picture
    	})
    return HttpResponse(template.render(context))
