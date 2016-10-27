'use strict';

angular.module('ircbotApp').provider('ConfigProvider', function() {
	this.$get = ['$resource', function($resource) {
		var config = $resource(formatUrl('/api/config/'), {}, {
			get: {
				method: 'GET',
				isArray: false
			}
		});

		return config;
	}];
});
