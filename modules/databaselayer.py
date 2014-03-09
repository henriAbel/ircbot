from databaseconnector import DatabaseConnector
from modules.logger import FileLogger

class DatabaseLayer():
	def __init__(self):
		self.database = DatabaseConnector()
		self.logger = FileLogger().log

	def add_user(self, username):
		return self.database.make_query("INSERT INTO irc_users (username) VALUES (?)", [username]).lastrowid

	def get_user(self, username):
		result = self.database.make_query("SELECT id FROM irc_users WHERE username = ?", [username], False).fetchone()
		self.database.close()
		if not result is None:
			return result[0]
		else:
			return None

	def add_message(self, message, userid, uuid):
		return self.database.make_query("INSERT INTO irc_message (content, user_id, uuid) VALUES (?, ?, ?)", [message, userid, uuid]).lastrowid

	def add_url(self, url, messageid, type):
		return self.database.make_query("INSERT INTO irc_link (content, message_id, type) VALUES (?, ?, ?)", [url, messageid, type]).lastrowid

	def change_nick(self, oldnick, newnick):
		self.database.make_query("UPDATE irc_users SET username = ? WHERE username = ? LIMIT 1", [oldnick, newnick])

	def link_exists(self, message, uuid = "unknown"):
		likeClause = '%' + message + '%'
		result = self.database.make_query("SELECT link.id FROM irc_link AS link JOIN irc_message AS message ON message.id = link.message_id WHERE link.content LIKE ? AND message.uuid != ?", [likeClause, uuid], False).fetchone()
		self.database.close()
		return not result is None

	def get_link(self, link):
		result = self.database.make_query("SELECT datetime(m.created, 'localtime'), u.username FROM irc_link as link JOIN irc_message as m on m.id = link.message_id JOIN irc_users as u ON u.id = m.user_id WHERE link.content = ?", [link], False).fetchone()
		self.database.close()
		return result

	def setDimensions(self, path, width, height):
		self.database.make_query("UPDATE irc_link SET width = ?, height = ? WHERE content = ?", [width, height, path])

	def setHash(self, path, hash):
		self.database.make_query("UPDATE irc_link SET hashLink = ? WHERE content = ?", [hash, path])

