from django.http import HttpResponse
from django.template import RequestContext, loader
from irc.models import Users, Link, Message, Linktype

def index(request):
    links = Link.objects.all()
    template = loader.get_template("irc/index.html")
    context = RequestContext(request, {
    	'links': links,
    	'linktype': Linktype
    	})
    return HttpResponse(template.render(context))
