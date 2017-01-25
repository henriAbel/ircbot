'use strict';

var loginUrl = '/api/login';

angular
	.module('ircbotApp', [
		'ngAnimate',
		'ngResource',
		'ngRoute',
		'ngSanitize'
	])
	.config(['$routeProvider', '$httpProvider', '$locationProvider', function ($routeProvider, $httpProvider, $locationProvider) {
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
			.when('/stats/', {
				templateUrl: 'views/stats.html',
				controller: 'StatController'
			})
			.when('/login/', {
				templateUrl: 'login.html',
				controller: 'LoginController'
			})
			.otherwise({
				redirectTo: '/'
			});
		$httpProvider.interceptors.push('AuthInterceptor');
 		$locationProvider.hashPrefix('');
	}])
	.run(['$document', '$rootScope', 'ConfigProvider', function($document, $rootScope, ConfigProvider) {
		$document.bind('keyup', function(e) {
			$rootScope.$broadcast('keypress', e);
		});
		ConfigProvider.get().$promise.then(function(response) {
			$rootScope.config = response.data;
		}, function(err) {});
		
	}]);

// TODO Remove??
var formatUrl = function(url) {
	return url;
};

/*
 * Appends ?authorization={token} to url
 */
var addToken = function(url) {
	if (window.sessionStorage.token === undefined) return url;
	return url + '?authorization=' + window.sessionStorage.token;
}

String.prototype.capitalizeFirstLetter = function() {
    return this.charAt(0).toUpperCase() + this.slice(1);
}

/*!
	Wheelzoom 3.0.4
	license: MIT
	http://www.jacklmoore.com/wheelzoom
*/
window.wheelzoom = (function(){
	var defaults = {
		zoom: 0.10
	};

	var canvas = document.createElement('canvas');

	var main = function(img, options){
		if (!img || !img.nodeName || img.nodeName !== 'IMG') { return; }

		var settings = {};
		var width;
		var height;
		var bgWidth;
		var bgHeight;
		var bgPosX;
		var bgPosY;
		var previousEvent;
		var cachedDataUrl;
		var displayBox;
		var isDragged = false;

		function setSrcToBackground(img) {
			img.style.backgroundImage = 'url("'+img.src+'")';
			img.style.backgroundRepeat = 'no-repeat';
			canvas.width = img.naturalWidth;
			canvas.height = img.naturalHeight;
			cachedDataUrl = canvas.toDataURL();
			img.src = cachedDataUrl;
		}

		function updateBgStyle() {
			if (bgPosX > 0) {
				bgPosX = 0;
			} else if (bgPosX < width - bgWidth) {
				bgPosX = width - bgWidth;
			}

			if (bgPosY > 0) {
				bgPosY = 0;
			} else if (bgPosY < height - bgHeight) {
				bgPosY = height - bgHeight;
			}

			img.style.backgroundSize = bgWidth+'px '+bgHeight+'px';
			img.style.backgroundPosition = bgPosX+'px '+bgPosY+'px';
		}

		function reset() {
			bgWidth = width;
			bgHeight = height;
			bgPosX = bgPosY = 0;
			updateBgStyle();
		}

		function onwheel(e) {
			var deltaY = 0;

			e.preventDefault();

			if (e.deltaY) { // FireFox 17+ (IE9+, Chrome 31+?)
				deltaY = e.deltaY;
			} else if (e.wheelDelta) {
				deltaY = -e.wheelDelta;
			}

			// As far as I know, there is no good cross-browser way to get the cursor position relative to the event target.
			// We have to calculate the target element's position relative to the document, and subtrack that from the
			// cursor's position relative to the document.
			var rect = img.getBoundingClientRect();
			var offsetX = e.pageX - rect.left - window.pageXOffset;
			var offsetY = e.pageY - rect.top - window.pageYOffset;

			// Record the offset between the bg edge and cursor:
			var bgCursorX = offsetX - bgPosX;
			var bgCursorY = offsetY - bgPosY;

			// Use the previous offset to get the percent offset between the bg edge and cursor:
			var bgRatioX = bgCursorX/bgWidth;
			var bgRatioY = bgCursorY/bgHeight;

			// Update the bg size:
			if (deltaY < 0) {
				bgWidth += bgWidth*settings.zoom;
				bgHeight += bgHeight*settings.zoom;
			} else {
				bgWidth -= bgWidth*settings.zoom;
				bgHeight -= bgHeight*settings.zoom;
			}

			// Take the percent offset and apply it to the new size:
			bgPosX = offsetX - (bgWidth * bgRatioX);
			bgPosY = offsetY - (bgHeight * bgRatioY);

			// Prevent zooming out beyond the starting size
			if (bgWidth <= width || bgHeight <= height) {
				reset();
			} else {
				updateBgStyle();
			}
		}

		function drag(e) {
			e.preventDefault();
			isDragged = true;
			bgPosX += (e.pageX - previousEvent.pageX);
			bgPosY += (e.pageY - previousEvent.pageY);
			previousEvent = e;
			updateBgStyle();
		}

		function removeDrag(e) {
			document.removeEventListener('mouseup', removeDrag);
			document.removeEventListener('mousemove', drag);
			var $scope = angular.element(displayBox).scope();
			if (!isDragged) {
				$scope.dragging = false;
			}
			else {
				$scope.$$postDigest(function(){
					$scope.dragging = false;
				});
			}
		}

		// Make the background draggable
		function draggable(e) {
			e.preventDefault();
			previousEvent = e;
			document.addEventListener('mousemove', drag);
			document.addEventListener('mouseup', removeDrag);
			angular.element(displayBox).scope().dragging = true;
			isDragged = false;
		}

		function load() {
			if (img.src === cachedDataUrl) return;

			var computedStyle = window.getComputedStyle(img, null);

			width = parseInt(computedStyle.width, 10);
			height = parseInt(computedStyle.height, 10);
			bgWidth = width;
			bgHeight = height;
			bgPosX = 0;
			bgPosY = 0;

			setSrcToBackground(img);

			img.style.backgroundSize =  width+'px '+height+'px';
			img.style.backgroundPosition = '0 0';
			img.addEventListener('wheelzoom.reset', reset);

			img.addEventListener('wheel', onwheel);
			img.addEventListener('mousedown', draggable);
			displayBox = $(img).parents('.display-box')[0]
		}

		var destroy = function (originalProperties) {
			img.removeEventListener('wheelzoom.destroy', destroy);
			img.removeEventListener('wheelzoom.reset', reset);
			img.removeEventListener('load', load);
			img.removeEventListener('mouseup', removeDrag);
			img.removeEventListener('mousemove', drag);
			img.removeEventListener('mousedown', draggable);
			img.removeEventListener('wheel', onwheel);

			img.style.backgroundImage = originalProperties.backgroundImage;
			img.style.backgroundRepeat = originalProperties.backgroundRepeat;
			img.src = originalProperties.src;
		}.bind(null, {
			backgroundImage: img.style.backgroundImage,
			backgroundRepeat: img.style.backgroundRepeat,
			src: img.src
		});

		img.addEventListener('wheelzoom.destroy', destroy);

		options = options || {};

		Object.keys(defaults).forEach(function(key){
			settings[key] = options[key] !== undefined ? options[key] : defaults[key];
		});

		if (img.complete) {
			load();
		}

		img.addEventListener('load', load);
	};

	// Do nothing in IE8
	if (typeof window.getComputedStyle !== 'function') {
		return function(elements) {
			return elements;
		};
	} else {
		return function(elements, options) {
			if (elements && elements.length) {
				Array.prototype.forEach.call(elements, main, options);
			} else if (elements && elements.nodeName) {
				main(elements, options);
			}
			return elements;
		};
	}
}());

'use strict';

 angular.module('ircbotApp').controller('ListController', ['$scope', 'LinkProvider', 'filter', '$compile', '$templateCache', function ($scope, LinkProvider, f, $compile, $templateCache) {
	$scope.offset = 0;
  /*
    If value is < 0, then itemsInPage function will recalculate value
    Window resize event will change value to -1
  */
  var itemsInPageLast = -1;
	$scope.itemsInPage = function() {
    if (itemsInPageLast > 0) return itemsInPageLast;
    var w = $('.content').width();
    var h = $(window).height();
    var columns = Math.floor(w / 125);
    var rows = Math.floor(h / 125);
    itemsInPageLast = columns * rows;
    return itemsInPageLast;
  };
	$scope.filter = f;
	$scope.displayObject = undefined;
    $scope.dragging = false;
	var tempObject;

	LinkProvider.getCount({filter: f}).$promise.then(function(o) {
		$scope.listCount = o.Count == 0 ? 1 : Math.ceil(o.Count / $scope.itemsInPage());
	});

	$scope.loadNext = function(callback) {
        if ($scope.links.length <  $scope.itemsInPage()) return;
		$scope.offset += $scope.itemsInPage();
		loadLinks(callback);
	};

	$scope.loadPrev = function(callback) {
        if ($scope.offset == 0) return;
		$scope.offset -= $scope.itemsInPage();
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
		LinkProvider.get({filter: f, offset: $scope.offset, limit: $scope.itemsInPage()}).$promise.then(function(e) {
            $scope.links = e;
            if (callback !== undefined) {
                callback.call(this);
            }
        });
		$scope.page = Math.ceil($scope.offset / $scope.itemsInPage()) + 1;
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
  var t = undefined;
  $(window).resize(function() {
    $scope.$apply(function(){
      if (t !== undefined) {
          clearTimeout(t);
      }
      t = setTimeout(function() {
        itemsInPageLast = -1;
        loadLinks();
      }, 250)
    });
  });

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

'use strict';

 angular.module('ircbotApp').controller('HeaderController', ['$scope', '$location', function ($scope, $location) {
 	$scope.isActive = function(linkLocation) {
 		return linkLocation === $location.path();
 	};
 }]);

'use strict';

 angular.module('ircbotApp').controller('LoginController', ['$scope', '$location', '$http', '$window', function ($scope, $location, $http, $window) {
	$scope.login = function(password) {
		$scope.loginError = "";
		$scope.loading = true;
		$http({
            method: 'POST',
            url: loginUrl,
            data: {
                password: password
            }
        }).then(function(response) {
        	$window.sessionStorage.token = response.data.token;
        	$location.path('/');
        }, function(error) {
        	$scope.loginError = error.Error;
			$scope.loading = false;
        });
	};
 }]);

'use strict';

angular.module('ircbotApp').controller('StatController', ['$scope', 'StatProvider', function ($scope, StatProvider) {
    var colors = ['#1BE7FF', '#6EEB83', '#E4FF1A', '#E8AA14', '#FF5714', '#50514F', '#F25F5C', '#247BA0', '#70C1B3'];
    var shuffle = function(o) {
        for(var j, x, i = o.length; i; j = Math.floor(Math.random() * i), x = o[--i], o[i] = o[j], o[j] = x);
        return o;
    }
    shuffle(colors);
    StatProvider.get().$promise.then(function(data) {
        var linkData = [];
        var userData = [];
        var duplicateData = [];
        for (var i = 0; i < data.GroupLink.length; i++) {
            var e = data.GroupLink[i]
            linkData.push({value: e.Count, label: e.Type.capitalizeFirstLetter(),
                color: getColor(i)})
        }
        for (var i = 0; i < data.GroupUser.length; i++) {
            var e = data.GroupUser[i];
            userData.push({value: e.Count, label: e.User_name.capitalizeFirstLetter(),
                color: getColor(i + data.GroupLink.length)})
        }
        for (var i = 0; i < data.Duplicates.length; i++) {
            var e = data.Duplicates[i];
            duplicateData.push({value: e.Count, label: e.User_name.capitalizeFirstLetter(),
                color: getColor(i + data.GroupLink.length + data.GroupUser.length)})
        }
        var ctx = document.getElementById("linkCanvas").getContext("2d");
        var ctx2 = document.getElementById("userCanvas").getContext("2d");
        var ctx3 = document.getElementById("duplicateCanvas").getContext("2d");
        var linkChart = new Chart(ctx).Doughnut(linkData);
        var userChart = new Chart(ctx2).Doughnut(userData);
        var duplicateChart = new Chart(ctx3).Doughnut(duplicateData);
    });

    var getColor = function(i) {
        if (colors.length <= i) {
            var start = i - colors.length * Math.floor(i / colors.length);
            return colors[start];
        }
        return colors[i];
    }
}]);

'use strict';

angular.module('ircbotApp').provider('LinkProvider', function() {
	this.$get = ['$resource', function($resource) {
		var links = $resource(formatUrl('/api/links/:param1/:param2'), {param1: '@param1', param2: '@param2'}, {
			get: {
				method: 'GET',
				isArray: true
			},
			getRaw: {
				method: 'GET',
			},
			getCount: {
				method: 'GET',
				params: {
					param1: 'count'
				}
			}
		});

		return links;
	}];
});
'use strict';

angular.module('ircbotApp').provider('StatProvider', function() {
	this.$get = ['$resource', function($resource) {
		var stats = $resource(formatUrl('/api/stat/:param1/'), {param1: '@param1'}, {
			get: {
				method: 'GET',
				isArray: false,
        params: {
          param1: 'all'
        }
			}
		});

		return stats;
	}];
});

'use strict';

angular.module('ircbotApp').provider('ConfigProvider', function() {
	this.$get = ['$resource', function($resource) {
		var config = $resource(formatUrl('/api/config/'), {}, {
			get: {
				method: 'GET',
				isArray: false
			}
		});

		return config;
	}];
});

'use strict';

angular.module('ircbotApp').directive('ngLink', ['$sce', '$rootScope', function($sce, $rootScope) {
	var getValueFromQuery = function(query, val) {
		var vars = query.substring(query.indexOf('?') + 1 ).split('&');
		for (var i = 0; i < vars.length; i++) {
			var pair = vars[i].split('=');
			if (decodeURIComponent(pair[0]) == val) {
				return decodeURIComponent(pair[1]);
			}
		}
	};

	return {
		restrict: 'A',
		scope: {
			ngModel: '=',
		},
		link: function(scope, element, attrs) {
			var showImage = function(e, url) {
				scope.$parent.showImage(e, {url: url, model: scope.ngModel});
			};
			if (scope.ngModel.Link_type == "image") {
				scope.contentUrl = formatUrl('/views/imageView.html')
				if ($rootScope.config.ExternalLibraries) {
					scope.imageUrl = addToken("/api/raw/" + scope.ngModel.Key + "/thumb")
					scope.imageClick = function(e) {
						showImage(e, addToken("/api/raw/" + scope.ngModel.Key + "/image"));
					}
				}
				else {
					scope.imageUrl = addToken("/api/raw/" + scope.ngModel.Key + "/image")
				}
			}
			else if (scope.ngModel.Link_type == "youtube") {
				scope.clicked = false;
				scope.videoID = getValueFromQuery(scope.ngModel.Link, "v");
				scope.videoUrl = $sce.trustAsResourceUrl("http://www.youtube.com/embed/" +
					 scope.videoID + "?autoplay=1");
				scope.contentUrl = formatUrl('/views/youtubeView.html')
				scope.youtubeClick = function() {
					scope.clicked = !scope.clicked;
					/*
						Something goes wrong when using iframe inside ng-if. Iframes
						are created for a short time and page loading time goes up. No time
						to search issue, so quick fix is to create iframes dynamically only when
						needed.
					*/
					if (scope.clicked) {
						var element = document.getElementById("video-" + scope.videoID);
						var youtubeFrame = document.createElement("iframe");
						youtubeFrame.type = "text/html";
						youtubeFrame.width = 640;
						youtubeFrame.height = 390;
						youtubeFrame.src = scope.videoUrl;
						youtubeFrame.frameborder = 0;
						youtubeFrame.id = "ytplayer-" + scope.videoID;

						element.parentNode.replaceChild(youtubeFrame, element);
					}
					else {
						var element = document.getElementById("ytplayer-" + scope.videoID);
						var placeholder = document.createElement("div");
						placeholder.id = "video-" + scope.videoID;
						element.parentNode.replaceChild(placeholder, element);
					}
				}
			}
			else if (scope.ngModel.Link_type == "link") {
				scope.contentUrl = formatUrl('/views/linkView.html');
			}
			else if (scope.ngModel.Link_type == "webm" || scope.ngModel.Link_type == "gif") {
				if ($rootScope.config.ExternalLibraries) {
						scope.contentUrl = formatUrl('/views/webmView.html');
						scope.thumbUrl = addToken("/api/raw/" + scope.ngModel.Key + "/webm1");
						scope.videoUrl = addToken("/api/raw/" + scope.ngModel.Key + "/webm");
				}
				else {
					scope.contentUrl = formatUrl('/views/gifView.html');
					scope.videoUrl = addToken("/api/raw/" + scope.ngModel.Key + "/gif");
				}
			}
		},
		template: '<div class="content-wrap" ng-include="contentUrl"></div>',
		replace: true,
	}
}]);

'use strict';

angular.module('ircbotApp').directive('ngErr', ['$interval', '$timeout', function($interval, $timeout) {
	return {
		link: function(scope, element, attrs) {
			element.bind('error', function() {
				if (element.data('interval') === undefined) {
					element.addClass('list-image-hidden')
					element.next().remove()
					element.after('<div class="sk-folding-cube"><div class="sk-cube1 sk-cube"></div><div class="sk-cube2 sk-cube"></div><div class="sk-cube4 sk-cube"></div><div class="sk-cube3 sk-cube"></div></div>');
					var interval = $interval(function() {
						attrs.$set('src', attrs.src);
					}, 5000);
					element.data('interval', interval);	
				}
			});
			element.bind('load', function() {
				$interval.cancel(element.data('interval'));
				element.next().remove()
				// Fixes 99% cases where Chrome isn't animating opacity 0 > 100
				$timeout(function() {
					element.removeClass('list-image-hidden')	
				}, 50);
				
			});
			element.on('$destroy', function() {
				if (element.data('interval') !== undefined) {
					$interval.cancel(element.data('interval'));
				}
			});
		}
	}
}]);
'use strict';

angular.module('ircbotApp').factory('AuthInterceptor', ['$rootScope', '$q', '$window', '$location',
		function($rootScope, $q, $window, $location) {
	return {
		request: function (config) {
			config.headers = config.headers || {};
			if ($window.sessionStorage.token) {
				config.headers.Authorization = 'Bearer ' + $window.sessionStorage.token;
			}
			return config;
		},
		response: function (response) {
			return response || $q.when(response);
		},
		responseError: function(response) {
			if (response.status === 401) {
				$location.path('/login/');
			}
			return $q.reject(response);
		}
	};
}]);