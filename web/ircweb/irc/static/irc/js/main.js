var t;
var href;
var popup;
var initWaterfall = function() {
	$('#waterfall').waterfall({
		itemCls: 'lazy',
		gutterWidth: 25,
		gutterHeight: 25,
		isFadeIn: true,
		checkImagesLoaded: true,
		path: function(page) {
			return href + '/' + page + '/ajax';
		},
		callbacks: {
			renderData: function (src, dataType) {
				if (src.count < 1) {
					stopWaterfall();
				}
				return src.data;
			},
			loadingFinished: function($loading, isBeyondMaxPage) {
				if ( !isBeyondMaxPage ) {
					$loading.fadeOut();
				} else {
					$loading.remove();
				}
				popup.reload();
			}
		}
	});
}

var stopWaterfall = function(callback) {
	$('#waterfall').waterfall('pause', callback);
}

var updateHref = function() {
	href = window.location.href;
	if (href.indexOf('/', href.length -1) !== -1) {
		href = href.substring(0, href.length -1);
	}
}

$(function() {
	updateHref();
	$('a[data-category="' + real_category + '"').parent().addClass('active');
	if (category.length > 1) {
		href += '/' + category;
	}
	popup = $('.lazy').popup();
	initWaterfall();

	$('#top-nav a').on('click', function(e) {
		history.pushState({}, '', this.href);
		updateHref();
		category = this.getAttribute('data-category');

		$('#top-nav li.active').removeClass('active');
		this.parentNode.className = 'active';
		stopWaterfall(function() {
			// Waterfall dosent provide any destroy/reload mehtods, so delete container object and reinit plugin
			$('#waterfall-message, #waterfall-loading').remove();
			var c = $('#waterfall');
			var parent = c.parent();
			c.remove();
			parent.append('<div id="waterfall"></div>')
			initWaterfall();
		});
		e.preventDefault();
	});
});