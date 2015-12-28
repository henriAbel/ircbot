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
			.when('/login/', {
				templateUrl: 'login.html',
				controller: 'LoginController'
			})
			.otherwise({
				redirectTo: '/'
			});
		$httpProvider.interceptors.push('AuthInterceptor');
	}])
	.run(['$document', '$rootScope', function($document, $rootScope) {
		$document.bind('keyup', function(e) {
			$rootScope.$broadcast('keypress', e);
		});
	}]);


var formatUrl = function(url) {
	return url;	
};
