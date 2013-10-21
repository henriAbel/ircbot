from twisted.words.protocols import irc
from logger import FileLogger, SqlLogger
from xml.dom import minidom
import time, re, urllib2, urlparse


class MessageHandler(irc.IRCClient):  
    def connectionMade(self):
        irc.IRCClient.connectionMade(self)
        self.logger = FileLogger(open(self.factory.filename, "a"))
        self.sqllogger = SqlLogger()

    def connectionLost(self, reason):
        irc.IRCClient.connectionLost(self, reason)
        self.logger.log_message("[disconnected at %s]" % 
                        time.asctime(time.localtime(time.time())))
        self.logger.close()


    def signedOn(self):
        self.join(self.factory.channel)

    # Converts string to unicode(utf-8) and sends to channel
    def say_decoded(self, channel, message):
        unicoded = message.encode(encoding='utf8')
        self.say(channel, unicoded)

    def privmsg(self, user, channel, msg):
        # Log every message once
        #if not self.sqllogger.message_exists(msg):
        user = user.split('!', 1)[0]
        self.logger.log_message("<%s> %s" % (user, msg)) 

        urls = re.findall('http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+', msg)
        for url in urls:
            id = self.sqllogger.log_message(msg, user)
            type = "normal"
            if any(url.endswith(s) for s in (".jpg", ".jpeg", ".gif", ".png")):
                type = "picture"
            elif url.find("youtube.com/watch") != -1:
                type = "youtube"

            self.sqllogger.log_url(url, id, type)

            if type == "youtube":
                url_data = urlparse.urlparse(url)
                query = urlparse.parse_qs(url_data.query)
                xml = urllib2.urlopen("http://gdata.youtube.com/feeds/api/videos/%s" % query["v"][0]) 
                xmldoc = minidom.parse(xml)
                self.say_decoded(channel, xmldoc.getElementsByTagName('title')[0].firstChild.nodeValue)

    def irc_NICK(self, prefix, params):
        self.sqllogger.log_nickchange(prefix.split('!')[0], params[0])

    def lineReceived(self, line):
        irc.IRCClient.lineReceived(self, line)
        #self.logger.log_message(line)
