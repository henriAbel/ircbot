$(function() {
	$('.content').on('click', '.youtubeLink', function() {
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

  	$('body').on('click', '.pageing > a, .menu_item', function(e) {
		e.preventDefault();
		var href = $(this).attr("href");
		if ($(e.target).hasClass('menu_item'))
			history.pushState({}, '', href.substring(0, href.length - 2));

		loadPage(href + "ajax");
		return false;
  	});

})

var loadPage = function(page) {
	try {
		var l = $('#loading');
		var c = $('.content');
		c.css('opacity', '0.2');
		l.addClass('loading');
		$.ajax({
		  type: 'GET',
		  url: page,
		  dataType: "json"
		}).done(function(msg) {
			$('.pageing').html(msg.page);
			$('.slides').html(msg.data);
			l.removeClass('loading');
			c.css('opacity', '1');
	 	}).fail(function(e) {
			console.log(e);
		}).complete(function(e, a) {
			console.log(a);
		});
	}
	catch (e) {
		console.log(e);
	}
}
