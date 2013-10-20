from configuration import Config
import MySQLdb, time

class FileLogger:
    def __init__(self, file):
        self.file = file

    def log_message(self, message):
        timestamp = time.strftime("[%H:%M:%S]", time.localtime(time.time()))
        self.file.write('%s %s\n' % (timestamp, message))
        self.file.flush()

    def close(self):
        self.file.close()

class SqlLogger:
    def __init__(self):
        self.users = {}
        self.config = Config("config")
        self.connection = MySQLdb.connect(host=self.config.get_option("host"),
                                 user=self.config.get_option("username"),
                                 passwd=self.config.get_option("password"),
                                 db=self.config.get_option("database"))
        self.cursor = self.connection.cursor()

    def log_url(self, url, message_id):
        type = "normal"
        if any(url.endswith(s) for s in (".jpg", ".jpeg", ".gif", ".png")):
            type = "picture"
        elif url.find("youtube.com/watch?v=") != -1:
            type = "youtube"
        self.make_insert("INSERT INTO irc_link (content, message_id, type) VALUES (%s, %s, %s)", (url, message_id, type))
        return self.cursor.lastrowid

    def log_message(self, message, user):
        id = self.check_user(user)
        self.make_insert("INSERT INTO irc_message (content, user_id) VALUES (%s, %s)", (message, id))
        return self.cursor.lastrowid

    # if new user, add it to the database
    def check_user(self, user):
        self.make_insert("SELECT id FROM irc_users WHERE username = %s", user)
        if self.cursor.rowcount > 0:
            if user not in self.users:
                self.users[user] = self.cursor.fetchone()[0]
        else:
            self.cursor.execute("INSERT INTO irc_users (username) VALUES (%s)", user)
            self.users[user] = self.cursor.lastrowid

        return self.users[user]

    def make_insert(self, query, param):
        self.cursor.execute(query, param)
        self.connection.commit()