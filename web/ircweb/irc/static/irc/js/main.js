var t;
var href;
var popup;
var initWaterfall = function() {
	$('.content').waterfall({
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
	$('.content').waterfall('pause', callback);	
}

var updateHref = function() {
	href = window.location.href;
	if (href.indexOf('/', href.length -1) !== -1) {
		href = href.substring(0, href.length -1);
	}
}

$(function() {
	updateHref();
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
			$('.content, #waterfall-message, #waterfall-loading').remove();
			$('.navbar').after('<div class="content"></div>')

			initWaterfall();	
		});
		e.preventDefault();
	});
});