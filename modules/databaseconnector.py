from configuration import Config
import sqlite3
import time

class DatabaseConnector():
	def __init__(self):
		#self.config = Config("config")
		self.connect_database()

	def connect_database(self):
		self.connection = sqlite3.connect('./web/ircweb/ircweb/ircwebDatabase.db')
		self.cursor = self.connection.cursor()
		'''
		try:
			self.connection = MySQLdb.connect(host=self.config.get_option("host"), 
	        	user=self.config.get_option("username"),
	        	passwd=self.config.get_option("password"),
	        	db=self.config.get_option("database"))        
			self.cursor = self.connection.cursor()
		except MySQLdb.OperationalError, message:
			raise RuntimeError('Colud not connect to sql server')
		'''

	def make_query(self, query, param):
		self.cursor.execute(query, param)
		if not query.startswith("select"):
			self.connection.commit()
		return self.cursor.lastrowid
