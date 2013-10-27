$(function() {
	$('.youtubeLink').click(function() {
		$('.youtubeLink').each(function() {
			var e = $(this);
			if (e.has('iframe')) {
				e.find('iframe').remove();
				e.find('img, .youtubeIcon').show();
			}
		});

		var e = $(this);
		e.find('img, .youtubeIcon').hide();
		var url = e.attr('data-url');
		e.append('<iframe width="480" height="360" src="' + url + '?autoplay=1" frameborder="0" allowfullscreen></iframe>');
	});

	// deprecated
	/*$('.leftMenu a').click(function() {
		var i = $('.menuItems > a').index(this);
		i++
		$('.active').hide(400).removeClass('active');
		$('.content > div:nth-child(' + i + ')').slideDown(400).addClass('active');
  	});*/

  	$('.fancybox-media').fancybox({
		openEffect  : 'none',
		closeEffect : 'none',
		helpers : {
			media : {}
		}
	});
})