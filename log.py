from twisted.internet import reactor, protocol, ssl
from modules.configuration import Config
from modules.messagehandler import MessageHandler
from modules.timer import Timer
import modules.databasemanager


class ClientFactory(protocol.ClientFactory):
    def __init__(self, channel, filename, channel_password = ""):
        self.channel = "%s %s" % (channel, channel_password)
        self.filename = "logs/%s" % filename

    def buildProtocol(self, addr):
        mhandler = MessageHandler()
        mhandler.nickname = self._nickname
        mhandler.password = self._password
        mhandler.factory = self
        return mhandler

    def clientConnectionLost(self, connector, reason):
        connector.connect()

    def clientConnectionFailed(self, connector, reason):
        print reason
        reactor.stop()

if __name__ == '__main__':
    modules.databasemanager.check_database()
    config = Config("config")
    factory = ClientFactory(config.get_option("channel"), config.get_option("log_file"), config.get_option("channel_password"))
    bot_name = config.get_option("bot_name")
    if bot_name is not None:
        factory._nickname = bot_name
    else:
        factory._nickname = "bot"

    factory._password = config.get_option("server_password")

    if config.get_option("ssl") == "true":
        reactor.connectSSL(config.get_option("server"), int(config.get_option("port")), factory, ssl.ClientContextFactory())
    else:
        reactor.connectTCP(config.get_option("server"), int(config.get_option("port")), factory)

    reactor.run()