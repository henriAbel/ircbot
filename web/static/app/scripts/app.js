'use strict';

angular
	.module('ircbotApp', [
		'ngAnimate',
		'ngResource',
		'ngRoute',
		'ngSanitize'
	])
	.config(function ($routeProvider) {
		$routeProvider
			.when('/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'image,gif'}
				}
			})
			.when('/links/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'link'}
				}
			})
			.when('/youtube/', {
				templateUrl: 'views/list.html',
				controller: 'ListController',
				resolve: {
					filter: function(){return 'youtube'}
				}
			})
			.otherwise({
				redirectTo: '/'
			});
	});


var formatUrl = function(url) {
	return url;	
}
