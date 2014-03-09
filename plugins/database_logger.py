from zope.interface import implements
from twisted.plugin import IPlugin
from modules.logger import FileLogger, SqlLogger
from modules.dimensions import Dimensions
from multiprocessing import Process
from xml.dom import minidom
import IMessage, urllib2

class databse_logger(object):
    implements(IPlugin, IMessage.IMessage)

    def __init__(self):
        self.logger = FileLogger().log
        self.sqllogger = SqlLogger()

    def onMessage(self, caller, message):
        pass

    def onLink(self, caller, message):
        # Log every message once
        if not self.sqllogger.database.link_exists(message.url):
            id = self.sqllogger.log_message(message.msg, message.user, message.uuid)
            self.sqllogger.log_url(message.url, id, message.type)
            if message.type == "gif" or message.type == "picture":
                # Downloading and converting takes time, so start in new process
                self.logger.debug("Starting subprocress {} {}".format(message.url, message.ext))
                p = Process(target=Dimensions, args=(message.url, message.ext))
                p.start()
        else:
            self.logger.info("Duplicate link: %s", message.url)

m = databse_logger()
