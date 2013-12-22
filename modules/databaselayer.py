from databaseconnector import DatabaseConnector

class DatabaseLayer():
	def __init__(self):
		self.database = DatabaseConnector()

	def add_user(self, username):
		return self.database.make_query("INSERT INTO irc_users (username) VALUES (?)", [username]).lastrowid

	def get_user(self, username):
		result = self.database.make_query("SELECT id FROM irc_users WHERE username = ?", [username], False).fetchone()
		self.database.close()
		if not result is None:
			return result[0]
		else:
			return None

	def add_message(self, message, userid):
		return self.database.make_query("INSERT INTO irc_message (content, user_id) VALUES (?, ?)", [message, userid]).lastrowid

	def add_url(self, url, messageid, type):
		return self.database.make_query("INSERT INTO irc_link (content, message_id, type) VALUES (?, ?, ?)", [url, messageid, type]).lastrowid

	def change_nick(self, oldnick, newnick):
		self.database.make_query("UPDATE irc_users SET username = ? WHERE username = ? LIMIT 1", [oldnick, newnick])

	def message_exists(self, message):
		result = self.database.make_query("SELECT id FROM irc_message WHERE content = ?", [message], False).fetchone()
		self.database.close()
		return not result is None

