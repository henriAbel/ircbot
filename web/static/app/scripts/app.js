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
				templateUrl: 'views/imageList.html',
				controller: 'ImageController'
			})
			.when('/links/', {
				templateUrl: 'views/listView.html',
				controller: 'LinkController'
			})     
			.otherwise({
				redirectTo: '/'
			});
	});


var formatUrl = function(url) {
	return url;	
}
