from apscheduler.scheduler import Scheduler
from databaseconnector import DatabaseConnector
from logger import FileLogger
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
		for row in self.database.make_query("""SELECT l.id, m.id, l.content, l.last_checked, l.type FROM irc_message m INNER JOIN irc_link l ON l.message_id = m.id 
			WHERE julianday(Date('now')) - julianday(Date(l.last_checked)) > ? AND l.type != 'youtube';""", [DATABASE_DIFF], False).fetchall():
			
			if self.isBroken(row[2], (row[4] == "picture" or row[4] == "gif")):
				self.log.info("Deleting link: %s" % row[2])
				self.database.make_query("DELETE FROM irc_link WHERE id = ?", [row[0]], False)
				self.database.make_query("DELETE FROM irc_message WHERE id = ?", [row[1]], False)
			else:
				self.database.make_query("UPDATE irc_link SET last_checked = Date('now') WHERE id = ?", [row[0]], False)

		self.database.close()

	def isBroken(self, link, image = True):
		try:
			request = urllib2.Request(link)
			request.get_method = lambda : "HEAD"

			response = urllib2.urlopen(request)
			code = response.getcode()
			info = response.info()
			self.log.info("Link %s returned with code %s" % (link, code))
			if int(code) >= 400:
				return True
			elif image and info["Content-Type"].startswith("text/html"):
				return True
			
			return False
		except Exception, e:
			return True

timer = Timer()
sched = Scheduler()
sched.start()
sched.add_interval_job(timer.job, hours=CHECK_INTERVAL)