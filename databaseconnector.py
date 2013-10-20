from configuration import Config
import MySQLdb

class Databaseconnector():
	def __init__(self):
		config = Config("config")
		self.connection = MySQLdb.connect(host=config.get_option("host"), 
        	user=config.get_option("username"),
        	passwd=config.get_option("password"),
        	db=config.get_option("database"))        
		self.cursor = self.connection.cursor()

	def add_user(self, username):
		return self.make_query("INSERT INTO irc_users (username) VALUES (%s)", username)

	def get_user(self, username):
		self.cursor.execute("SELECT id FROM irc_users WHERE username = %s", username)
		if self.cursor.rowcount > 0:
			return self.cursor.fetchone()[0]
		else:
			return None

	def add_message(self, message, userid):
		return self.make_query("INSERT INTO irc_message (content, user_id) VALUES (%s, %s)", (message, userid))

	def add_url(self, url, messageid, type):
		return self.make_query("INSERT INTO irc_link (content, message_id, type) VALUES (%s, %s, %s)", (url, messageid, type))

	def change_nick(self, oldnick, newnick):
		self.make_query("UPDATE irc_users SET username = %s WHERE username = %s LIMIT 1", oldnick, newnick)

	def message_exists(self, message):
		self.cursor.execute("SELECT id FROM irc_message WHERE content = %s", message)
		return self.cursor.rowcount > 0

	def make_query(self, query, param):
		self.cursor.execute(query, param)
		self.connection.commit()
		return self.cursor.lastrowid