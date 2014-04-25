from zope.interface import implements
from twisted.plugin import IPlugin
from modules.logger import FileLogger
from xml.dom import minidom
import IMessage, urllib2, datetime

class youtube(object):
	implements(IPlugin, IMessage.IMessage)

	def __init__(self):
		self.logger = FileLogger().log
		self.logger.info("youtube init")

	def onMessage(self, caller, message):
		pass

	def onLink(self, caller, message):
		self.logger.info("got youtube link")
		try:
			if message.type == "youtube":
				self.logger.info("opening youtube:")
				xml = urllib2.urlopen("http://gdata.youtube.com/feeds/api/videos/%s" % message.url)
				xmldoc = minidom.parse(xml)
				self.logger.info("got response, posting:")
				caller.send_message(xmldoc.getElementsByTagName('title')[0].firstChild.nodeValue)
				self.logger.info("post done:")
		except urllib2.HTTPError, e:
			self.logger.info("Something went wrong while resolving youtube link %s", e)
		except Exception:
			import traceback
			self.logger.info("Fatal error! If you see this, please post following message to the github issuse resolver");
			self.logger.info(traceback.format_exc())

m = youtube()
