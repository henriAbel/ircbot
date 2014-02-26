import os
from django.db import models
from datetime import datetime
from linkwrapper import Linkwrapper

""" Fixed width and height for thumbnail image. Only one value is used based on ratio """
THUMBWIDTH = 250
THUMBHEIGTH = 220

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
	width = models.IntegerField()
	height = models.IntegerField()
	hashLink = models.CharField(max_length=64)

	def youtubelink(self):
		return "http://www.youtube.com/embed/%s" % self.content

	def thumbLink(self):
		if self.hashLink is None:
			return self.content

		return "/static/irc/thumb/" + self.hashLink

	def fullLink(self):
		if self.hashLink is None:
			return self.content

		return "/static/irc/fullphoto/" + self.hashLink

	def thumbWidth(self):
		if not hasattr(self, "tWidth"):
			self.calculateThumbsize()

		return self.tWidth

	def thumbHeigth(self):
		if not hasattr(self, "tHeigth"):
			self.calculateThumbsize()

		return self.tHeigth


	def calculateThumbsize(self):
		ratio = self.width / float(self.height)
		#if self.width > self.height:
		self.tWidth = THUMBWIDTH
		self.tHeigth = self.tWidth / ratio
		#else:
		#	self.tHeigth = THUMBHEIGTH
		#	self.tWidth = self.tHeigth * ratio
