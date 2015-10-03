'use strict';

var loginUrl = '/api/login';

angular
	.module('ircbotApp', [
		'ngAnimate',
		'ngResource',
		'ngRoute',
		'ngSanitize'
	])
	.config(['$routeProvider', '$httpProvider', function ($routeProvider, $httpProvider) {
		$routeProvider
			.when('/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'image,gif';}
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
			.when('/login/', {
				templateUrl: 'login.html',
				controller: 'LoginController'
			})
			.otherwise({
				redirectTo: '/'
			});
		$httpProvider.interceptors.push('AuthInterceptor');
	}]);


var formatUrl = function(url) {
	return url;	
};
