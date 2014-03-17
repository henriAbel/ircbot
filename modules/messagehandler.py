import IMessage, plugins
from twisted.words.protocols import irc
from twisted.plugin import getPlugins
from logger import FileLogger
from dimensions import Dimensions
from irc_message import IrcMessage
from databaselayer import DatabaseLayer
import time, re, urlparse, os, uuid, urlActions

class MessageHandler(irc.IRCClient):
    def connectionMade(self):
        irc.IRCClient.connectionMade(self)
        self.logger = FileLogger().log
        self.logger.info("Connected to %s" % self.factory.channel)
        self.databaselayer = DatabaseLayer()

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
        if msg.startswith("!"):
            self.commandReceived(msg[1:])
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
            # Check if we are dealing with dropbox links
            if url.find("www.dropbox.com/s/") != -1:
                #https://www.dropbox.com/help/201/en
                url += "?dl=1"
            contentType = urlActions.getContentType(url)["type"]
            # Only end of content type is important
            contentType = contentType.replace("image/", "")
            if contentType != "" and any(contentType in s for s in ("jpg", "jpeg", "png")):
                type = "pictures"
            elif contentType == "gif":
                type = "gifs"
            elif url.find("youtube.com/watch") != -1:
                type = "youtube"
                url_data = urlparse.urlparse(url)
                query = urlparse.parse_qs(url_data.query)
                url = query["v"][0]
            else:
                type = "links"

            m.type = type
            m.url = url
            m.ext = contentType

            for pl in getPlugins(IMessage.IMessage, plugins):
                pl.onLink(self, m)

    def irc_NICK(self, prefix, params):
        before = prefix.split('!')[0]
        now = params[0]
        self.logger.info("%s changed nick to %s" % (before, now))

    def lineReceived(self, line):
        irc.IRCClient.lineReceived(self, line)

    def commandReceived(self, command):
        if command == "version":
            self.send_message(str(self.databaselayer.getVersion()))
