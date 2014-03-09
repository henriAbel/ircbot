from zope.interface import Interface, Attribute

class IMessage(Interface):

    def onMessage(caller, message):
    	pass
    def onLink(caller, message):
    	pass
