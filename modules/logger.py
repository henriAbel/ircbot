from databaselayer import DatabaseLayer
from modules.configuration import Config
import time
import logging

class FileLogger:
    def __init__(self):
        config = Config("config")
        logging.basicConfig(filename="logs/" + config.get_option("log_file"), level=20, format='%(asctime)s %(levelname)-8s %(message)s')
        self.log = logging.getLogger("")

class SqlLogger:
    def __init__(self):
        self.users = {}
        self.database = DatabaseLayer()

    def log_url(self, url, messageid, type):
        return self.database.add_url(url, messageid, type)

    def log_message(self, message, user):
        id = self.check_user(user)
        return self.database.add_message(message, id)

    # if new user, add it to the database
    def check_user(self, user):
        id = self.database.get_user(user)
        if id is not None:
            if user not in self.users:
                self.users[user] = id
        else:
            self.users[user] = self.database.add_user(user)

        return self.users[user]

    def log_nickchange(self, oldnick, newnick):
        self.database.change_nick(oldnick, newnick)

    def message_exists(self, message):
        return self.database.message_exists(message)