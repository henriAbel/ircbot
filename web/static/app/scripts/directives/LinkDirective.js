'use strict';

angular.module('ircbotApp').directive('ngLink', function($sce) {
	var getGifUrlFromModel = function(model, clicked) {
		if (clicked) {
			return "/api/raw/" + model.Key + "/gif";
		}
		return "/api/raw/" + model.Key + "/gif1";
	}

	var getValueFromQuery = function(query, val) {
		var vars = query.substring(query.indexOf('?') + 1 ).split('&');
		for (var i = 0; i < vars.length; i++) {
			var pair = vars[i].split('=');
			if (decodeURIComponent(pair[0]) == val) {
				return decodeURIComponent(pair[1]);
			}
		}
	}

	return {
		restrict: 'A',
		scope: {
			ngModel: '=',
		},
		link: function(scope, element, attrs) {
			scope.clicked = false;
			if (scope.ngModel.Link_type == "gif") {
				scope.contentUrl = formatUrl('/views/gifView.html')
				scope.gifUrl = getGifUrlFromModel(scope.ngModel, false);
				scope.gifClick = function() {
					scope.clicked = !scope.clicked;
					scope.gifUrl = getGifUrlFromModel(scope.ngModel, scope.clicked);
				}	
			}
			else if (scope.ngModel.Link_type == "image") {
				scope.contentUrl = formatUrl('/views/imageView.html')
			}
			else if (scope.ngModel.Link_type == "youtube") {
				scope.videoID = getValueFromQuery(scope.ngModel.Link, "v");
				scope.videoUrl = $sce.trustAsResourceUrl("http://www.youtube.com/embed/" + 
					 scope.videoID + "?autoplay=1");
				scope.contentUrl = formatUrl('/views/youtubeView.html')
				scope.youtubeClick = function() {
					scope.clicked = !scope.clicked;
				}
			}
			else if (scope.ngModel.Link_type == "link") {
				scope.contentUrl = formatUrl('/views/linkView.html')
			}
		},
		template: '<div class="content-wrap" ng-include="contentUrl"></div>',
		replace: true,
		controller: ['$scope', 'LinkProvider', function($scope, $sce) {

		}],
	}
});