from django.db import models
from datetime import datetime
from linkwrapper import Linkwrapper

class Linktype():
	NORMAL = Linkwrapper("link", "Links")
	PICTURE = Linkwrapper("picture", "Pictures")
	YOUTUBE = Linkwrapper("youtube", "Youtube")
	GIF = Linkwrapper("gif", "Gifs")

class Users(models.Model):
	username = models.CharField(max_length=100)
	last_seen = models.DateTimeField()
	created = models.DateTimeField(editable=False)

	def save(self):
	    if not self.id:
	        self.created = datetime.date.today()
	    super(Users, self).save()

class Message(models.Model):
	content = models.CharField(max_length=512)
	user = models.ForeignKey(Users)
	created = models.DateTimeField(editable=False)

	def save(self):
		if not self.id:
			self.created = datetime.date.today()
		super(Message, self).save()

	class Meta:
		get_latest_by = "created"

class Link(models.Model):
	content = models.CharField(max_length=512)
	message = models.ForeignKey(Message)
	type = models.CharField(max_length=50)
	
	def youtubelink(self):
		return "http://www.youtube.com/embed/%s" % self.content		