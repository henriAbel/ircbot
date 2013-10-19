from configuration import Config
import MySQLdb

if __name__ == '__main__':
	config = Config("config")
	connection = MySQLdb.connect(host=config.get_option("host"),
								 user=config.get_option("username"),
								 passwd=config.get_option("password"))
	cursor = connection.cursor()
	if config.get_option("database") is None:
		# check if database exists
		database_name = config.get_option("username") + "_irc"
		cursor.execute("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = %s", database_name)
		if cursor.rowcount == 0:
			cursor.execute("CREATE DATABASE IF NOT EXISTS %s", database_name);
	else:
		database_name = config.get_option("database")

	print database_name
	cursor.execute("use %s" % database_name) 
	# check if all tables exists
	cursor.execute("""CREATE TABLE IF NOT EXISTS `users` (
  					  `id` int(15) NOT NULL AUTO_INCREMENT,
  					  `username` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  					  `last_seen` datetime NULL,
  					  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  					  PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1 """)
	cursor.execute("""CREATE TABLE IF NOT EXISTS `message` (
					  `id` int(25) NOT NULL AUTO_INCREMENT,
					  `content` varchar(512) COLLATE utf8_unicode_ci NOT NULL,
					  `user_id` int(15) NOT NULL,
					  PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1 ;""")
	cursor.execute("""CREATE TABLE IF NOT EXISTS `link` (
					  `id` int(25) NOT NULL AUTO_INCREMENT,
					  `content` varchar(512) COLLATE utf8_unicode_ci NOT NULL,
					  `message_id` int(25) NOT NULL,
					  `type` enum('picture','normal', 'youtube') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'normal',
					  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					  PRIMARY KEY (`id`)
					) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci AUTO_INCREMENT=1 ;""")
	print "database created"




	
