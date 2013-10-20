# Copyright (c) Twisted Matrix Laboratories.
# See LICENSE for details.



# twisted imports
from twisted.internet import reactor, protocol, ssl

# another imports
from configuration import Config
from message_handler import MessageHandler

class MessageHandlerFactory(protocol.ClientFactory):
    def __init__(self, channel, filename, channel_password = ""):
        self.channel = "%s %s" % (channel, channel_password)
        self.filename = filename

    def buildProtocol(self, addr):
        mhandler = MessageHandler()
        mhandler.nickname = "bot"
        mhandler.factory = self
        return mhandler

    def clientConnectionLost(self, connector, reason):
        connector.connect()

    def clientConnectionFailed(self, connector, reason):
        print reason
        reactor.stop()

if __name__ == '__main__':
    config = Config("config")
    factory = MessageHandlerFactory(config.get_option("channel"), config.get_option("log_file"), config.get_option("channel_password"))

    if config.get_option("ssl") == "true":
        reactor.connectSSL(config.get_option("server"), int(config.get_option("port")), factory, ssl.ClientContextFactory())
    else:
        reactor.connectTCP(config.get_option("server"), int(config.get_option("port")), factory)

    reactor.run()
