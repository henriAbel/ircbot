from databaseconnector import DatabaseConnector
from configuration import Config

def check_database():
	config = Config("config")
	database = DatabaseConnector()
	database.connect_database()
	cursor = database.cursor
	connection = database.connection

	# check if all tables exists
	cursor.execute("""CREATE TABLE IF NOT EXISTS `irc_users` (
				    `id` integer PRIMARY KEY AUTOINCREMENT,
				    `username` varchar(100) NOT NULL,
				    `last_seen` TIMESTAMP NULL,
				    `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
					)""")

	cursor.execute("""CREATE TABLE IF NOT EXISTS `irc_message` (
				    `id` integer PRIMARY KEY AUTOINCREMENT,
				    `content` varchar(512) NOT NULL,
				    `user_id` integer NOT NULL,
				    `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				    CONSTRAINT `user_id_refs_id_61cb1a8b` FOREIGN KEY (`user_id`) REFERENCES `irc_users` (`id`)
					)""")
	cursor.execute("""CREATE TABLE IF NOT EXISTS `irc_link` (
				    `id` integer PRIMARY KEY AUTOINCREMENT,
				    `content` varchar(512) NOT NULL,
				    `message_id` integer NOT NULL,
				    `type` varchar(50) NOT NULL,
				    CONSTRAINT `message_id_refs_id_12db337f` FOREIGN KEY (`message_id`) REFERENCES `irc_message` (`id`)
					)""")

	cursor.execute("""CREATE TABLE IF NOT EXISTS `irc_versioning` (
				    `id` integer PRIMARY KEY AUTOINCREMENT,
				    `version` decimal(2, 1) NOT NULL,
				    `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
					)""")

	cursor.execute("CREATE INDEX IF NOT EXISTS `irc_message_user_index` ON `irc_message` (`user_id`);")
	cursor.execute("CREATE INDEX IF NOT EXISTS `irc_message_created_index` ON `irc_message` (`created`);")
	cursor.execute("CREATE INDEX IF NOT EXISTS `irc_link_message_index` ON `irc_link` (`message_id`);")
	cursor.execute("CREATE INDEX IF NOT EXISTS `irc_link_type_index` ON `irc_link` (`type`);")
	connection.commit()

	# Simple versioning

	cursor.execute("SELECT version FROM irc_versioning ORDER BY created DESC LIMIT 1")
	row = cursor.fetchone()
	if row is None:
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.1)")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN last_checked TIMESTAMP DEFAULT '2000-01-01T00:00:00.000'")
	if float(row[0]) == float(0.1):
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.2)")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN width integer NULL")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN height integer NULL")
	if float(row[0]) == float(0.2):
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.3)")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN hashLink varchar(64) NULL")
	if float(row[0] == float(0.3)):
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.4)")
		# Delete duplicate links
		ids = cursor.execute("SELECT id, message_id FROM irc_link AS l WHERE (SELECT COUNT(*) FROM irc_link WHERE content = l.content) > 1 LIMIT -1 OFFSET 1;").fetchall()
		for id in ids:
			cursor.execute("DELETE FROM irc_link WHERE id = ?", [id[0]])
			cursor.execute("DELETE FROM irc_message WHERE id = ?", [id[1]])
		cursor.execute("ALTER TABLE irc_message ADD COLUMN uuid integer NULL")

	connection.commit()
	connection.close()
