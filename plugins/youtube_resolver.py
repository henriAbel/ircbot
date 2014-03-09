from zope.interface import implements
from twisted.plugin import IPlugin
from modules.logger import FileLogger
from xml.dom import minidom
import IMessage, urllib2

class youtube(object):
	implements(IPlugin, IMessage.IMessage)

	def __init__(self):
		self.logger = FileLogger().log

	def onMessage(self, caller, message):
		pass

	def onLink(self, caller, message):
		try:
			if message.type == "youtube":
				xml = urllib2.urlopen("http://gdata.youtube.com/feeds/api/videos/%s" % message.url)
				xmldoc = minidom.parse(xml)
				caller.send_message(xmldoc.getElementsByTagName('title')[0].firstChild.nodeValue)
		except urllib2.HTTPError, e:
			self.logger.info("Something went wrong while resolving youtube link %s", e)
		except Exception:
			import traceback
			self.logger.info("Fatal error! If you see this, please post following message to the github issuse resolver");
			self.logger.info(traceback.format_exc())

m = youtube()
