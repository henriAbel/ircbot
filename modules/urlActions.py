import urllib2

def getContentType(url):
	try:
		request = urllib2.Request(url)
		request.get_method = lambda : "HEAD"

		response = urllib2.urlopen(request)
		c = response.getcode()
		info = response.info()
		return {"code": c, "type": info["Content-Type"]}
	except Exception, e:
		return False
