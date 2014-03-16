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

	update = True
	while update:
		update = doDatabaseUpdate(cursor)
		connection.commit()

	connection.commit()
	connection.close()

def doDatabaseUpdate(cursor):
	cursor.execute("SELECT version FROM irc_versioning ORDER BY version DESC LIMIT 1")
	row = cursor.fetchone()
	if row is None:
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.1)")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN last_checked TIMESTAMP DEFAULT '2000-01-01T00:00:00.000'")
		return True
	if float(row[0]) == float(0.1):
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.2)")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN width integer NULL")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN height integer NULL")
		return True
	if float(row[0]) == float(0.2):
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.3)")
		cursor.execute("ALTER TABLE irc_link ADD COLUMN hashLink varchar(64) NULL")
		return True
	if float(row[0] == float(0.3)):
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.4)")
		# Delete duplicate links
		ids = cursor.execute("SELECT id, GROUP_CONCAT(id) FROM irc_link AS l WHERE (SELECT COUNT(*) FROM irc_link WHERE content = l.content) > 1 GROUP BY content;").fetchall()
		for id in ids:
			linkIds = id[1].split(",")
			first = True
			for delid in linkIds:
				if first:
					first = False
					continue
				cursor.execute("DELETE FROM irc_message WHERE id = (SELECT message_id FROM irc_link WHERE id = ?)", [delid])
				cursor.execute("DELETE FROM irc_link WHERE id = ?", [delid])

		cursor.execute("ALTER TABLE irc_message ADD COLUMN uuid integer NULL")
		return True
	if float(row[0] == float(0.4)):
		cursor.execute("INSERT INTO irc_versioning (version) VALUES (0.5)")
		cursor.execute("UPDATE irc_link SET type = 'gifs' WHERE type = 'gif'")
		cursor.execute("UPDATE irc_link SET type = 'pictures' WHERE type = 'picture'")
		cursor.execute("UPDATE irc_link SET type = 'links' WHERE type = 'link'")
		return True
	return False
