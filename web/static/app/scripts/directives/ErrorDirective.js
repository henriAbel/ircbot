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
			})
			element.bind('load', function() {
				$interval.cancel(element.data('interval'));
				console.log('onLoad')
				element.next().remove()
				// Fixes 99% cases where Chrome isn't animating opacity 0 > 100
				$timeout(function() {
					element.removeClass('list-image-hidden')	
				}, 50);
				
			})
		}
	}
}]);