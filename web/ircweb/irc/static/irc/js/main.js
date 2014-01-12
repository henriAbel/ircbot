var t;
var initWaterfall = function() {
	$('.content').waterfall({
		itemCls: 'lazy',
		gutterWidth: 25,
		gutterHeight: 25,
		isFadeIn: true,
		checkImagesLoaded: false,
		debug: false,
		path: function(page) {
			return '/' + category + '/' + page + '/ajax';
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
				if (t) {clearTimeout(t)}
				t = setTimeout(function() {
					$('.lazy').popup();
					console.log('ok');
				}, 500)
			}
		}
	});
}

var stopWaterfall = function(callback) {
	$('.content').waterfall('pause', callback);	
}

$(function() {
	initWaterfall();

	$('#top-nav a').on('click', function(e) {
		e.preventDefault();
		var href = this.href;
		category = this.getAttribute('data-category');

		$('#top-nav li.active').removeClass('active');
		this.parentNode.className = 'active';
		stopWaterfall(function() {
			$('.content, #waterfall-message, #waterfall-loading').remove();
			$('.navbar').after('<div class="content"></div>')

			initWaterfall();	
		});
		
		history.pushState({}, '', href);
		
	});

	
});