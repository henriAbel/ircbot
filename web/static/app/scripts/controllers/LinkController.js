'use strict';

 angular.module('ircbotApp').controller('LinkController', ['$scope', 'LinkProvider', function ($scope, LinkProvider) {
 	$scope.links = LinkProvider.get({filter: 'link'});
 }]);
