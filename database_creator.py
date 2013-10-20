from configuration import Config
import MySQLdb

if __name__ == '__main__':
	config = Config("config")
	connection = MySQLdb.connect(host=config.get_option("host"),
								 user=config.get_option("username"),
								 passwd=config.get_option("password"))
	cursor = connection.cursor()
	database_name = ""
	if config.get_option("database") is None:
		database_name = config.get_option("username") + "_irc"
	else:
		database_name = config.get_option("database")
		
	#cursor.execute("CREATE DATABASE IF NOT EXISTS %s" % database_name);
	cursor.execute("use %s" % database_name) 
	# check if all tables exists
	cursor.execute("""CREATE TABLE IF NOT EXISTS `irc_users` (
				    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
				    `username` varchar(100) NOT NULL,
				    `last_seen` TIMESTAMP NULL,
				    `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
					)""")

	cursor.execute("""CREATE TABLE IF NOT EXISTS `irc_message` (
				    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
				    `content` varchar(512) NOT NULL,
				    `user_id` integer NOT NULL,
				    `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				    INDEX `irc_message_user` (`user_id`),
				    INDEX `irc_message_created` (`created`),
				    CONSTRAINT `user_id_refs_id_61cb1a8b` FOREIGN KEY (`user_id`) REFERENCES `irc_users` (`id`)
					)""")
	cursor.execute("""CREATE TABLE IF NOT EXISTS `irc_link` (
				    `id` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
				    `content` varchar(512) NOT NULL,
				    `message_id` integer NOT NULL,
				    `type` varchar(50) NOT NULL,
				    INDEX `irc_link_message` (`message_id`),
				    INDEX `irc_link_type` (`type`),
				    CONSTRAINT `message_id_refs_id_12db337f` FOREIGN KEY (`message_id`) REFERENCES `irc_message` (`id`)
					)""")

	print "database created"
