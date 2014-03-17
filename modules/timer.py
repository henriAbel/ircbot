from apscheduler.scheduler import Scheduler
from databaseconnector import DatabaseConnector
from multiprocessing import Process
from modules.dimensions import Dimensions
from logger import FileLogger
import urlActions
import urllib2

# How often run timer job
CHECK_INTERVAL = 1
# Check links older then days
DATABASE_DIFF = 1

class Timer():
	def __init__(self):
		self.log = FileLogger().log
		self.database = DatabaseConnector()

	def job(self):
		self.log.info("Detecting broken links")
		# Get all links where last checked time is greater than 1 day
		for row in self.database.make_query("""SELECT l.id, m.id, l.content, l.last_checked, l.type, l.hashLink FROM irc_message m INNER JOIN irc_link l ON l.message_id = m.id
			WHERE julianday(Date('now')) - julianday(Date(l.last_checked)) > ? AND (l.type == 'links' OR (l.type == 'gifs' AND l.hashLink IS NULL)) ;""", [DATABASE_DIFF], False).fetchall():
			# Check if there are any gifs with empty hashLink
			if row[4] == "gifs":
				# If such links exists then try to populate hashLink or delete from database
				self.log.info("Found gif: {} with empty hahsLink".format(row[2]))
				content = urlActions.getContentType(row[2])
				if content is None or content["tpye"] != "image/gif":
					self.deleteRecord(row)
					self.log.info("Deleting gif: {} from database due to broken link".format(row[2]))
					continue
				p = Process(target=Dimensions, args=(row[2], content["type"].replace("image/", "")))
				p.start()
			elif self.isBroken(row[2]):
				self.deleteRecord(row)
			else:
				self.database.make_query("UPDATE irc_link SET last_checked = Date('now') WHERE id = ?", [row[0]], False)

		self.database.close()

	def isBroken(self, link):
		content = urlActions.getContentType(link)
		if not content:
			self.log.info("Something went wrong while checking link %s" % link)
			return False

		self.log.info("Link %s returned with code %s" % (link, content["code"]))
		if int(content["code"]) >= 400:
			return True

		#TODO Handle error codes < 400
		return False

	def deleteRecord(self, row):
		self.log.info("Deleting link: %s" % row[2])
		self.database.make_query("DELETE FROM irc_link WHERE id = ?", [row[0]], False)
		self.database.make_query("DELETE FROM irc_message WHERE id = ?", [row[1]], False)

timer = Timer()
sched = Scheduler()
sched.start()
sched.add_interval_job(timer.job, hours=CHECK_INTERVAL)
