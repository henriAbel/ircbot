import os, urllib, Image

GIFDIR = os.path.join(os.getcwd(), "web/ircweb/irc/static/irc/gif/")
THUMBDIR = os.path.join(os.getcwd(), "web/ircweb/irc/static/irc/thumb/")

class GifExtractor():
	def __init__(self, image):
		self.thumb = GIFDIR + os.path.basename(image)
		if not os.path.exists(self.thumb):
			urllib.urlretrieve(image, self.thumb)
			self.getFirstFrame()
			self.save()

	def getFirstFrame(self):
		self.im = Image.open(self.thumb)

	def save(self):
		try:
			# Save with .png extension
			self.im.save(THUMBDIR + os.path.basename(self.thumb).replace("gif", "png"), "PNG")
		except IOError:
			print "Cannot save gif: ", self.thumb
