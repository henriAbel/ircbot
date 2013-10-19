import ConfigParser

class Config:
	def __init__(self, file):
		self.file = file
		self.config = ConfigParser.ConfigParser()
		self.config.read(self.file)
		self.mapper()

	def mapper(self):
		self.dictonary = {}
		for section in self.config.sections():
			for option in self.config.options(section):
				if len(option) > 0:
					self.dictonary[option] = self.config.get(section, option)

	def get_option(self, option):
		if option in self.dictonary:
			return self.dictonary[option]
