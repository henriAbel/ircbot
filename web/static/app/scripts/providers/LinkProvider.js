'use strict';

angular.module('ircbotApp').provider('LinkProvider', function() {
	this.$get = ['$resource', function($resource) {
		var links = $resource(formatUrl('/api/links/:param1/:param2'), {param1: '@param1', param2: '@param2'}, {
			get: {
				method: 'GET',
				isArray: true
			},
			getRaw: {
				method: 'GET',
			},
			getCount: {
				method: 'GET',
				params: {
					param1: 'count'
				}
			}
		});

		return links;
	}];
});