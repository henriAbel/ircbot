'use strict';

angular.module('ircbotApp').factory('AuthInterceptor', ['$rootScope', '$q', '$window', '$location',
		function($rootScope, $q, $window, $location) {
	return {
		request: function (config) {
			config.headers = config.headers || {};
			if ($window.sessionStorage.token) {
				config.headers.Authorization = 'Bearer ' + $window.sessionStorage.token;
			}
			return config;
		},
		response: function (response) {
			return response || $q.when(response);
		},
		responseError: function(response) {
			if (response.status === 401) {
				$location.path('/login/');
			}
			return $q.reject(response);
		}
	};
}]);