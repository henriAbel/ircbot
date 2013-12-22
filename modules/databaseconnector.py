from configuration import Config
import sqlite3
import time

class DatabaseConnector():
	#def __init__(self):
		#self.config = Config("config")
		#self.connect_database()

	def connect_database(self):
		self.connection = sqlite3.connect('./web/ircweb/ircweb/ircwebDatabase.db')
		self.connection.text_factory = str
		self.cursor = self.connection.cursor()

	def make_query(self, query, param):
		self.connect_database()
		self.cursor.execute(query, param)
		if not query.startswith("select"):
			self.connection.commit()
		self.connection.close()

		return self.cursor.lastrowid
