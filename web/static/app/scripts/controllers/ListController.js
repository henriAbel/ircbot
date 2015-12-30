'use strict';

 angular.module('ircbotApp').controller('ListController', ['$scope', 'LinkProvider', 'filter', '$compile', '$templateCache', function ($scope, LinkProvider, f, $compile, $templateCache) {
	$scope.offset = 0;
	$scope.itemsInPage = 30;
	$scope.filter = f;
	$scope.displayObject = undefined;
    $scope.dragging = false;
	var tempObject;

	LinkProvider.getCount({filter: f}).$promise.then(function(o) {
		$scope.listCount = o.Count == 0 ? 1 : Math.ceil(o.Count / $scope.itemsInPage);
	});

	$scope.loadNext = function(callback) {
        if ($scope.links.length <  $scope.itemsInPage) return;
		$scope.offset += $scope.itemsInPage;
		loadLinks(callback);
	};

	$scope.loadPrev = function(callback) {
        if ($scope.offset == 0) return;
		$scope.offset -= $scope.itemsInPage;
		if ($scope.offset < 0) {
			$scope.offset = 0;
		}
		loadLinks(callback);
	};

	$scope.$on('keypress', function(event, e) {
        if (e.which == 39 || e.which == 37) {
            if ($scope.displayObject !== undefined) {
                var i = getElementPosition($scope.displayObject);
                var nextElement;
                if (e.which == 39) {
                    i++;
                    if (i >= $scope.links.length) {
                        $scope.loadNext(function() {
                            var nextModel = $scope.links[0];
                            $scope.$$postDigest(function() {
                                $scope.showImage($('.image-list-item:first'),
                                 {url: addToken("/api/raw/" + nextModel.Key + "/image"), model: nextModel})
                            });
                        });
                        $scope.hideImage();
                        return;
                    }
                    nextElement = $scope.displayObject.domElement.nextAll('.image-list-item:first');
                }
                else {
                    i--;
                    if (i < 0) {
                        $scope.loadPrev(function() {
                            var nextModel = $scope.links[$scope.links.length -1];
                            $scope.$$postDigest(function() {
                                $scope.showImage($('.image-list-item:last'),
                                 {url: addToken("/api/raw/" + nextModel.Key + "/image"), model: nextModel})
                            });
                        });
                        $scope.hideImage();
                        return;
                    }
                    nextElement = $scope.displayObject.domElement.prevAll('.image-list-item:first');
                }
                var nextModel = $scope.links[i];
                $scope.showImage(nextElement, {url: addToken("/api/raw/" + nextModel.Key + "/image"), model: nextModel})
            }
            else {
              if(e.which == 39) {
                $scope.loadNext();
              }
              else if (e.which == 37) {
                $scope.loadPrev()
              }
            }
        }
        else if (e.which == 27) { //ESC
            $scope.hideImage();
        }
	});

	var loadLinks = function(callback) {
		LinkProvider.get({filter: f, offset: $scope.offset, limit: $scope.itemsInPage}).$promise.then(function(e) {
            $scope.links = e;
            if (callback !== undefined) {
                callback.call(this);
            }
        });
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
        var clickedParent = e.target ? $(e.target).parents('.image-list-item') : e;
		var last = getLastElementInRow(clickedParent);
		tempObject = $compile($templateCache.get('displayBoxTemplate.html'))($scope);
		last.after(tempObject);
		$scope.displayObject = o;
        $scope.displayObject.domElement = clickedParent;
        wheelzoom(document.querySelectorAll('.display-box img'));
	};

	$scope.hideImage = function() {
        if ($scope.dragging) return;
        document.querySelector('.display-box img').dispatchEvent(new CustomEvent('wheelzoom.destroy'));
		$scope.displayObject = undefined;
		tempObject.remove();
		tempObject = undefined;
	};

	var getLastElementInRow = function(clickedParent) {
		var allImages = clickedParent.nextAll('.image-list-item');
		for (var i = 0; i < allImages.length; i++) {
			if (clickedParent.position().top !== $(allImages[i]).position().top) {
				if (i === 0) {
					return clickedParent;
				}
				return $(allImages[i-1]);
			}

		}

		if (allImages.length > 0) {
			return $(allImages[allImages.length-1]);
		}
		return clickedParent;
	};

    var getElementPosition = function(e) {
      for (var i = 0; i < $scope.links.length; i++) {
          var o = $scope.links[i];
          if (o.Key == e.model.Key) {
              return i;
          }
      }
    };

  loadLinks();
 }]);
