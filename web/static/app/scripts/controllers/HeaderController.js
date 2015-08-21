'use strict';

 angular.module('ircbotApp').controller('HeaderController', ['$scope', '$location', function ($scope, $location) {
 	$scope.isActive = function(linkLocation) {
 		return linkLocation == $location.path();
 	}
 }]);
