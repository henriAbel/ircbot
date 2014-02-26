from PIL import Image
from databaselayer import DatabaseLayer
from logger import FileLogger
import urllib2
import hashlib
import io
import os

THUMBDIR = os.path.join(os.getcwd(), "web/ircweb/irc/static/irc/thumb/")
IMAGEDIR = os.path.join(os.getcwd(), "web/ircweb/irc/static/irc/fullphoto/")

class Dimensions(object):
	"""Get photo/gif dimensions and save to database"""
	def __init__(self, path, ext):
		self.logger = FileLogger().log
		self.database = DatabaseLayer()

		# Get image
		self.img = self.getRemote(path)
		if self.img is None:
			return
		self.im = Image.open(self.img)
		size = self.im.size
		self.logger.info("{0} returned with size w:{1}px h:{2}px".format(path, size[0], size[1]))

		# Calculate hash from link
		nameHash = hashlib.md5()
		nameHash.update(path)
		self.dbfile = nameHash.hexdigest() + ext
		self.logger.debug("Link: {} MD5 filename is: {}".format(path, self.dbfile))

		# Save image width and height to database
		self.database.setDimensions(path, size[0], size[1])
		self.database.setHash(path, self.dbfile)

		# Save image to local drive
		self.saveFullphoto()
		self.saveThumbnail()

	def getRemote(self, path):
		try:
			file = urllib2.urlopen(path)
			image = io.BytesIO(file.read())
			return image
		except Exception, e:
			self.logger.info("Cannot open url {}".format(path))
			return None

	def saveThumbnail(self):
		# Make thumbnail image from original
		size = 210, 210
		self.im.thumbnail(size)
		self.save(THUMBDIR + self.dbfile)

	def saveFullphoto(self):
		# save byteIO to baseDir folder without modifications
		with open(IMAGEDIR + self.dbfile, "wb") as local:
			local.write(self.img.getvalue())

	def save(self, path):
		""" Saves current image in im object """
		try:
			self.im.save(path)
		except IOError:
			self.logger.info("Cannot save image: {}".format(self.dbfile))
