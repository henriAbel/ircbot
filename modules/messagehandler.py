import IMessage, plugins
from twisted.words.protocols import irc
from twisted.plugin import getPlugins
from logger import FileLogger
from dimensions import Dimensions
from irc_message import IrcMessage
import time, re, urlparse, os, uuid

class MessageHandler(irc.IRCClient):
    def connectionMade(self):
        irc.IRCClient.connectionMade(self)
        self.logger = FileLogger().log
        self.logger.info("Connected to %s" % self.factory.channel)

    def connectionLost(self, reason):
        irc.IRCClient.connectionLost(self, reason)
        self.logger.info("disconnected from %s" % self.factory.channel)

    def signedOn(self):
        self.channel = self.factory.channel.split(" ")[0]
        self.join(self.factory.channel)

    # Converts string to unicode(utf-8) and sends to channel
    def say_decoded(self, channel, message):
        unicoded = message.encode(encoding='utf8')
        self.say(channel, unicoded)

    def send_message(self, message):
        self.say_decoded(self.channel, message)

    def privmsg(self, user, channel, msg):
        m = IrcMessage()
        m.uuid = uuid.uuid1().hex
        m.msg = msg
        m.channel = channel
        user = user.split('!', 1)[0]
        m.user = user
        self.logger.info("<%s> %s" % (user, msg))
        for pl in getPlugins(IMessage.IMessage, plugins):
            pl.onMessage(self, m)

        urls = re.findall('http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+', msg)
        for url in urls:
            type = "link"
            path = urlparse.urlparse(url).path
            ext = os.path.splitext(path)[1]
            if ext != "" and any(ext in s for s in (".jpg", ".jpeg", ".png")):
                type = "picture"
            elif ext == ".gif":
                type = "gif"
            elif url.find("youtube.com/watch") != -1:
                type = "youtube"
                url_data = urlparse.urlparse(url)
                query = urlparse.parse_qs(url_data.query)
                url = query["v"][0]

            m.type = type
            m.url = url
            m.ext = ext

            for pl in getPlugins(IMessage.IMessage, plugins):
                pl.onLink(self, m)

    def irc_NICK(self, prefix, params):
        before = prefix.split('!')[0]
        now = params[0]
        self.logger.info("%s changed nick to %s" % (before, now))

    def lineReceived(self, line):
        irc.IRCClient.lineReceived(self, line)
