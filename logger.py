from configuration import Config
import MySQLdb

class FileLogger:
    def __init__(self, file):
        self.file = file

    def log_message(self, message):
        timestamp = time.strftime("[%H:%M:%S]", time.localtime(time.time()))
        self.file.write('%s %s\n' % (timestamp, message))
        self.file.flush()
        self.file.close()

class SqlLogger:
    def __init__(self, host, username, password, database):
        self.config = Config("config")
        self.connection = MySQLdb.connect(host=config.getOption("host"),
                                 user=config.getOption("username"),
                                 passwd=config.getOption("password"))

