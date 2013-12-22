from databaseconnector import DatabaseConnector

class DatabaseLayer():
	def __init__(self):
		self.database = DatabaseConnector()

	def add_user(self, username):
		return self.database.make_query("INSERT INTO irc_users (username) VALUES (?)", [username])

	def get_user(self, username):
		self.database.make_query("SELECT id FROM irc_users WHERE username = ?", [username])
		if self.database.cursor.rowcount > 0:
			return self.database.cursor.fetchone()[0]
		else:
			return None

	def add_message(self, message, userid):
		return self.database.make_query("INSERT INTO irc_message (content, user_id) VALUES (?, ?)", [message, userid])

	def add_url(self, url, messageid, type):
		return self.database.make_query("INSERT INTO irc_link (content, message_id, type) VALUES (?, ?, ?)", [url, messageid, type])

	def change_nick(self, oldnick, newnick):
		self.database.make_query("UPDATE irc_users SET username = ? WHERE username = ? LIMIT 1", [oldnick, newnick])

	def message_exists(self, message):
		self.database.make_query("SELECT id FROM irc_message WHERE content = ?", [message])
		return self.database.cursor.rowcount > 0

