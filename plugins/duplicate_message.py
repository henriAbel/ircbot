from zope.interface import implements
from twisted.plugin import IPlugin
from modules.logger import SqlLogger
import IMessage

class message_plugin(object):
	implements(IPlugin, IMessage.IMessage)

	def __init__(self):
		self.sqllogger = SqlLogger()

	def onMessage(self, caller, message):
		pass

	def onLink(self, caller, message):
		if self.sqllogger.database.link_exists(message.url, message.uuid):
			link = self.sqllogger.database.get_link(message.url)
			caller.send_message("Old link {} {}".format(link[0], link[1]))

m = message_plugin()
