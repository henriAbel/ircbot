from django import template
from irc.models import Link

register = template.Library()

@register.inclusion_tag('irc/youtube_panel.html', takes_context=True)
def youtube(context, link):
	return {"link": link}

@register.inclusion_tag('irc/gif_panel.html', takes_context=True)
def gif(context, link):
	return {"link": link}

@register.inclusion_tag('irc/photo_panel.html', takes_context=True)
def photo(context, link):
	return {"link": link}

@register.inclusion_tag('irc/link_panel.html', takes_context=True)
def url(context, link):
	return {"link": link}