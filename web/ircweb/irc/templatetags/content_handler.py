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
def page(context, next, prev, category):
	return {
        "next": next, 
        "prev": prev,
        "category": category
    }

@register.inclusion_tag('irc/menu_panel.html', takes_context=True)
def menu(context, active):
	menus = [Linktype.NORMAL, Linktype.YOUTUBE, Linktype.GIF, Linktype.PICTURE]
	return {
		"menus": menus,
		"active": active
	}