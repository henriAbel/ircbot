from django import template
from irc.models import Link, Linktype

register = template.Library()

@register.inclusion_tag('irc/youtube_panel.html', takes_context=True)
def youtube(context, links):
	return {"links": links}

@register.inclusion_tag('irc/gif_panel.html', takes_context=True)
def gif(context, links):
	return {"links": links}

@register.inclusion_tag('irc/photo_panel.html', takes_context=True)
def photo(context, links):
	return {"links": links}

@register.inclusion_tag('irc/link_panel.html', takes_context=True)
def url(context, links):
	return {"links": links}

@register.inclusion_tag('irc/page_panel.html', takes_context=True)
def page(context, total, current, category):
	print total
	return {
	"total": range(1,int(total) + 1), 
	"current": int(current),
	"category": category
	}

@register.inclusion_tag('irc/menu_panel.html', takes_context=True)
def menu(context):
	menus = [Linktype.NORMAL, Linktype.YOUTUBE, Linktype.GIF, Linktype.PICTURE]
	return {"menus": menus}