$(function() {
	$('.youtubeLink').click(function() {
		$('.youtubeLink').each(function() {
			var e = $(this);
			if (e.has('iframe')) {
				e.find('iframe').remove();
				e.find('img').show();
			}
		});

		var e = $(this);
		e.find('img').hide();
		var url = e.attr('data-url');
		e.append('<iframe width="480" height="360" src="' + url + '?autoplay=1" frameborder="0" allowfullscreen></iframe>');
	})
})