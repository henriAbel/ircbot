'use strict';

 angular.module('ircbotApp').controller('ImageController', ['$scope', 'LinkProvider', function ($scope, LinkProvider) {
 	$scope.links = LinkProvider.get({filter: 'image,gif,youtube'});
 }]);
