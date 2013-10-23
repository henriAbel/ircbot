from configuration import Config
import MySQLdb

class DatabaseConnector():
	def __init__(self):
		self.config = Config("config")
		self.loop_counter = 0
		self.connect_database()

	def connect_database(self):
		try:
			self.connection = MySQLdb.connect(host=self.config.get_option("host"), 
	        	user=self.config.get_option("username"),
	        	passwd=self.config.get_option("password"),
	        	db=self.config.get_option("database"))        
			self.cursor = self.connection.cursor()
		except MySQLdb.OperationalError, message:
			raise RuntimeError('Colud not connect to sql server')

	def make_query(self, query, param):
		try:
			self.cursor = self.connection.cursor()
			self.loop_counter = 0
			self.cursor.execute(query, param)
			self.connection.commit()
			self.cursor.close()
			return self.cursor.lastrowid
		except MySQLdb.OperationalError, message:
			# If after 5 reconnect connection could not be restored, give up
			if self.loop_counter < 5:
				self.connect_database();
				self.loop_counter += 1
				self.make_query(query, param)	