'use strict';

 angular.module('ircbotApp').controller('ListController', ['$scope', 'LinkProvider', 'filter', function ($scope, LinkProvider, f) {
 	$scope.offset = 0;
 	$scope.itemsInPage = 30;
 	$scope.filter = f;
 	
 	LinkProvider.getCount({filter: f}).$promise.then(function(o) {
 		$scope.listCount = Math.ceil(o.Count / 30)
 	});
 	
 	$scope.loadNext = function() {
 		$scope.offset += $scope.itemsInPage;
 		loadLinks();
 	};

 	$scope.loadPrev = function() {
		$scope.offset -= $scope.itemsInPage;
		if ($scope.offset < 0) {
			$scope.offset = 0;
		}
		loadLinks();
 	};

 	var loadLinks = function() {
 		$scope.links = LinkProvider.get({filter: f, offset: $scope.offset});
 		$scope.page = Math.ceil($scope.offset / $scope.itemsInPage) + 1;
 	};

 	loadLinks();
 }]);
