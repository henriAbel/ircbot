'use strict';

angular.module('ircbotApp').provider('StatProvider', function() {
	this.$get = ['$resource', function($resource) {
		var stats = $resource(formatUrl('/api/stat/:param1/'), {param1: '@param1'}, {
			get: {
				method: 'GET',
				isArray: false,
        params: {
          param1: 'all'
        }
			}
		});

		return stats;
	}];
});
