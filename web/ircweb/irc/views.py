from django.http import HttpResponse
from django.template import RequestContext, loader
from irc.models import Users, Link, Message, Linktype
import json, math

ITEM_IN_PAGE = 5

def index(request, category = None):
	template = loader.get_template("irc/index.html")
	""" Set category only then no category is defined. Javascript concats category with request url"""
	if category is None:
		category = Linktype.NORMAL.db
	else:
		category = ""
	context = RequestContext(request, {
    	"category": category
   	})

	return HttpResponse(template.render(context))

def ajax(request, category, page = 1):
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

	limit = ITEM_IN_PAGE * page
	# Get html from template
	template = loader.get_template("irc/%s.html" % template_file)
	links = getLinkCollection(category, page)
	context = RequestContext(request, {
    	"links": links
    	})
	total = Link.objects.filter(type=category).count()
	total_pages = int(math.ceil(total / float(ITEM_IN_PAGE)))
	
	context2 = RequestContext(request, {
		"category": category,
		"next": getNextPage(total_pages, page),
		"prev": getPrevPage(total_pages, page)
		})
	# Put all togheter into json object
	resp = dict()
	resp.update(count = links.count())
	resp.update(data = template.render(context))
	return HttpResponse(json.dumps(resp), content_type="application/json")

def getLinkCollection(category, page):
	# Calcluate offset and limit
	page = int(page)
	offset = ITEM_IN_PAGE * (page -1)
	limit = ITEM_IN_PAGE * page
	return Link.objects.order_by("-id").filter(type=category)[offset:limit]

def getNextPage(total, current):
	try:
		return current + 1 if total - current > 0 else -1
	except TypeError, e:
		return -1

def getPrevPage(total, current):
	try:
		return current - 1 if current > 1 else -1
	except TypeError, e:
		return -1


