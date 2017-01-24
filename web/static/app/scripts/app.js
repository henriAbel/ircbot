'use strict';

var loginUrl = '/api/login';

angular
	.module('ircbotApp', [
		'ngAnimate',
		'ngResource',
		'ngRoute',
		'ngSanitize'
	])
	.config(['$routeProvider', '$httpProvider', '$locationProvider', function ($routeProvider, $httpProvider, $locationProvider) {
		$routeProvider
			.when('/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'image';}
				}
			})
			.when('/links/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'link';}
				}
			})
			.when('/youtube/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'youtube';}
				}
			})
			.when('/video/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'webm,gif';}
				}
			})
			.when('/stats/', {
				templateUrl: 'views/stats.html',
				controller: 'StatController'
			})
			.when('/login/', {
				templateUrl: 'login.html',
				controller: 'LoginController'
			})
			.otherwise({
				redirectTo: '/'
			});
		$httpProvider.interceptors.push('AuthInterceptor');
 		$locationProvider.hashPrefix('');
	}])
	.run(['$document', '$rootScope', 'ConfigProvider', function($document, $rootScope, ConfigProvider) {
		$document.bind('keyup', function(e) {
			$rootScope.$broadcast('keypress', e);
		});
		ConfigProvider.get().$promise.then(function(response) {
			$rootScope.config = response.data;
		}, function(err) {});
		
	}]);

// TODO Remove??
var formatUrl = function(url) {
	return url;
};

/*
 * Appends ?authorization={token} to url
 */
var addToken = function(url) {
	if (window.sessionStorage.token === undefined) return url;
	return url + '?authorization=' + window.sessionStorage.token;
}

String.prototype.capitalizeFirstLetter = function() {
    return this.charAt(0).toUpperCase() + this.slice(1);
}
