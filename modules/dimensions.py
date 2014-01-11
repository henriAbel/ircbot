from PIL import Image
from databaselayer import DatabaseLayer
from logger import FileLogger
import urllib2
import io

class Dimensions(object):
	"""Get photo/gif dimensions and save to database"""
	def __init__(self, path, local = False, realpath = ""):
		self.logger = FileLogger().log
		self.database = DatabaseLayer()
		if local:
			im = Image.open(path)
		else:
			im = Image.open(self.getRemote(path))

		size = im.size
		self.logger.info("{0} returned with size w:{1}px h:{2}px".format(path, size[0], size[1]))
		dbpath = path if local == False else realpath
		self.database.setDimensions(dbpath, size[0], size[1])

	def getRemote(self, path):
		file = urllib2.urlopen(path)
		image = io.BytesIO(file.read())
		return image