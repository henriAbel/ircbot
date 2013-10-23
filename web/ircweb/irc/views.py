from django.http import HttpResponse
from django.template import RequestContext, loader
from irc.models import Users, Link, Message, Linktype

def index(request):
	links = []
	youtube = []
	picture = []
	gif = []
	links_collection = Link.objects.all()

	for link in links_collection:
		if link.type == Linktype.NORMAL:
			links.insert(0, link)
		elif link.type == Linktype.PICTURE:
			picture.insert(0, link)
		elif link.type == Linktype.GIF:
			gif.insert(0, link)
		else:
			youtube.insert(0, link)

	template = loader.get_template("irc/index.html")
	context = RequestContext(request, {
    	'links': links,
    	'youtube': youtube,
    	'picture': picture,
    	'gif' : gif
    	})

	return HttpResponse(template.render(context))


