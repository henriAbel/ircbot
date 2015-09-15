'use strict';

 angular.module('ircbotApp').controller('ListController', ['$scope', 'LinkProvider', 'filter', '$compile', '$templateCache', function ($scope, LinkProvider, f, $compile, $templateCache) {
 	$scope.offset = 0;
 	$scope.itemsInPage = 30;
 	$scope.filter = f;
 	$scope.displayObject = undefined;
 	var tempObject;
 	
 	LinkProvider.getCount({filter: f}).$promise.then(function(o) {
 		$scope.listCount = Math.ceil(o.Count / $scope.itemsInPage)
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

 	$scope.showImage = function(e, o) {
 		if (tempObject !== undefined) {
 			if (o.model.Key === $scope.displayObject.model.Key) {
 				$scope.hideImage();
 				return;
 			}
 			$scope.hideImage();
 		}
 		var last = getLastElementInRow(e);
		tempObject = $compile($templateCache.get('displayBoxTemplate.html'))($scope);
		last.after(tempObject);
 		$scope.displayObject = o;
 	}

 	$scope.hideImage = function() {
 		$scope.displayObject = undefined;
 		tempObject.remove();
 		tempObject = undefined;
 	}

 	var getLastElementInRow = function(e) {
		var clickedParent = $(e.target).parents('.image-list-item');
		var allImages = clickedParent.nextAll(".image-list-item");
		for (var i = 0; i < allImages.length; i++) {
			if (clickedParent.position().top !== $(allImages[i]).position().top) {
				if (i == 0) {
					return clickedParent;
				}
				 return $(allImages[i-1]);
			}
			
		};
		return clickedParent;
	}

 	loadLinks();
 }]);
