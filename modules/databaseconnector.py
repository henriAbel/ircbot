from configuration import Config
import MySQLdb

class DatabaseConnector():
	def __init__(self):
		self.config = Config("config")
		self.loop_counter = 0
		self.connect_database()

	def connect_database(self):
		self.connection = MySQLdb.connect(host=self.config.get_option("host"), 
        	user=self.config.get_option("username"),
        	passwd=self.config.get_option("password"),
        	db=self.config.get_option("database"))        
		self.cursor = self.connection.cursor()

	def make_query(self, query, param):
		if self.connection.open:
			self.loop_counter = 0
			self.cursor.execute(query, param)
			self.connection.commit()
			return self.cursor.lastrowid
		# If after 5 reconnect connection could not be restored, give up
		elif self.loop_counter < 5:
			self.connect_database();
			self.loop_counter += 1
			self.make_query(query, param)
		else:
			raise RuntimeError('Colud not connect to sql server')