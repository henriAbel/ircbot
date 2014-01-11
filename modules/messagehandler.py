from twisted.words.protocols import irc
from logger import FileLogger, SqlLogger
from xml.dom import minidom
from gifextract import GifExtractor
from dimensions import Dimensions
import time, re, urllib2, urlparse, dimensions
from multiprocessing import Process


class MessageHandler(irc.IRCClient):  
    def connectionMade(self):
        irc.IRCClient.connectionMade(self)
        self.logger = FileLogger().log
        self.sqllogger = SqlLogger()
        self.logger.info("Connected to %s" % self.factory.channel)

    def connectionLost(self, reason):
        irc.IRCClient.connectionLost(self, reason)
        self.logger.info("disconnected from %s" % self.factory.channel)

    def signedOn(self):
        self.join(self.factory.channel)

    # Converts string to unicode(utf-8) and sends to channel
    def say_decoded(self, channel, message):
        unicoded = message.encode(encoding='utf8')
        self.say(channel, unicoded)

    def privmsg(self, user, channel, msg):
        user = user.split('!', 1)[0]
        self.logger.info("<%s> %s" % (user, msg)) 

        urls = re.findall('http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+', msg)
        for url in urls:
            type = "link"
            if any(url.endswith(s) for s in (".jpg", ".jpeg", ".png")):
                type = "picture"
            elif url.endswith(".gif"):
                type = "gif"
            elif url.find("youtube.com/watch") != -1:
                type = "youtube"
                url_data = urlparse.urlparse(url)
                query = urlparse.parse_qs(url_data.query)
                url = query["v"][0]
            # Log every message once
            if not self.sqllogger.message_exists(msg):
                id = self.sqllogger.log_message(msg, user)
                self.sqllogger.log_url(url, id, type)
                if type == "gif":
                    # Downloading and converting takes time, so start in new process
                    p = Process(target=GifExtractor, args=(url,))
                    p.start()
                elif type == "picture":
                    p = Process(target=Dimensions, args=(url,))
                    p.start()
            try:
                if type == "youtube":
                    xml = urllib2.urlopen("http://gdata.youtube.com/feeds/api/videos/%s" % url) 
                    xmldoc = minidom.parse(xml)
                    self.say_decoded(channel, xmldoc.getElementsByTagName('title')[0].firstChild.nodeValue)
            except urllib2.HTTPError, e:
                self.logger.info("Something went wrong while resolving youtube link %s", e)
            except Exception:
                import traceback
                self.logger.info("Fatal error! If you see this, please post following message to the github issuse resolver");
                self.logger.info(traceback.format_exc())

    def irc_NICK(self, prefix, params):
        before = prefix.split('!')[0]
        now = params[0]
        self.logger.info("%s changed nick to %s" % (before, now))
        #self.sqllogger.log_nickchange(befor, now)

    def lineReceived(self, line):
        irc.IRCClient.lineReceived(self, line)
        #self.logger.info(line)
