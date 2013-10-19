from configuration import Config
import MySQLdb

if __name__ == '__main__':
	config = Config("config")
	connection = MySQLdb.connect(host=config.getOption("host"),
								 user=config.getOption("username"),
								 passwd=config.getOption("password"))
	cursor = connection.cursor()
	if config.getOption("database") is None:
		cursor.execute("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = %s", (config.getOption("username") + "_irc"))
		if cursor.rowcount == 0:
			print "not implemented"
			#TODO make database
	
