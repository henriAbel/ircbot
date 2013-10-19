from twisted.words.protocols import irc
import time, re
from logger import FileLogger, SqlLogger

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

    def privmsg(self, user, channel, msg):
        user = user.split('!', 1)[0]
        self.logger.log_message("<%s> %s" % (user, msg)) 

        urls = re.findall('http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+', msg)
        for url in urls:
            id = self.sqllogger.log_message(msg, user)
            self.sqllogger.log_url(url, id)


    def irc_NICK(self, prefix, params):
        old_nick = prefix.split('!')[0]
        new_nick = params[0]
        self.logger.log_message("%s is now known as %s" % (old_nick, new_nick))

    def lineReceived(self, line):
        irc.IRCClient.lineReceived(self, line)
        #self.logger.log_message(line)