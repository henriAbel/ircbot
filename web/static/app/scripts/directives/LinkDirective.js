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
			var showImage = function(e, url) {
				scope.$parent.showImage(e, {url: url, model: scope.ngModel});
			}
			if (scope.ngModel.Link_type == "gif") {
				scope.contentUrl = formatUrl('/views/gifView.html')
				scope.gifUrl = getGifUrlFromModel(scope.ngModel, false);
				scope.gifClick = function(e) {
					showImage(e, getGifUrlFromModel(scope.ngModel, true));
				}	
			}
			else if (scope.ngModel.Link_type == "image") {
				scope.contentUrl = formatUrl('/views/imageView.html')
				scope.imageUrl = "/api/raw/" + scope.ngModel.Key + "/thumb"
				scope.imageClick = function(e) {
					showImage(e, "/api/raw/" + scope.ngModel.Key + "/full");
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
				scope.contentUrl = formatUrl('/views/linkView.html')
			}
		},
		template: '<div class="content-wrap" ng-include="contentUrl"></div>',
		replace: true,
		controller: ['$scope', 'LinkProvider', function($scope, $sce) {

		}],
	}
});