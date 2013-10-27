from django.http import HttpResponse
from django.template import RequestContext, loader
from irc.models import Users, Link, Message, Linktype
import json, math

ITEM_IN_PAGE = 15

def index(request, category = Linktype.NORMAL.db, page = 1):
	links = []
	youtube = []
	picture = []
	gif = []
	links_collection = getLinkCollection(category, page)
	total = Link.objects.filter(type=category).count()
	template = loader.get_template("irc/index.html")
	context = RequestContext(request, {
    	category: links_collection,
    	"total": math.ceil(total / float(ITEM_IN_PAGE)),
    	"current": int(page),
    	"category": category
    	})

	return HttpResponse(template.render(context))

def ajax(request, category, page = 1):
	if request.POST.get("ajax") is None:
		return index(request)

	if category == Linktype.GIF.db:
		template_file = "gif_panel"
	elif category == Linktype.PICTURE.db:
		template_file = "photo_panel"
	elif category == Linktype.YOUTUBE.db:
		template_file = "youtube_panel"
	else:
		template_file = "link_panel"

	page = int(page)
	if page < 1:
		page = 1

	total = Link.objects.filter(type=category).count()
	# Get html from template
	template = loader.get_template("irc/%s.html" % template_file)
	context = RequestContext(request, {
    	"links": getLinkCollection(category, page)
    	})
	# Put all togheter into json object
	resp = dict()
	# Check if there is next page
	if total <= limit:
		resp.update(more = False)
	else:
		resp.update(more = True)
	resp.update(total = total)
	resp.update(current = page)

	resp.update(data = template.render(context))
	return HttpResponse(json.dumps(resp), content_type="application/json")

def getLinkCollection(category, page):
	# Calcluate offset and limit
	page = int(page)
	offset = ITEM_IN_PAGE * (page -1)
	limit = ITEM_IN_PAGE * page
	return Link.objects.filter(type=category)[offset:limit]


