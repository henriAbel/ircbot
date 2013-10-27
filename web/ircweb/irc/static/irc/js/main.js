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

  	// Is it necessary?
  	$('.fancybox-media').fancybox({
		openEffect  : 'none',
		closeEffect : 'none',
		helpers : {
			media : {}
		}
	});

	$('.menu_item').click(function(e) {
		e.preventDefault();
		var href = $(this).attr("href");
		history.pushState({}, '', href);

		loadPage(href + "ajax");
		return false;
	});

})

var loadPage = function(page) {
	try {
		$.ajax({
		  type: 'GET',
		  url: page,
		  dataType: "json"
		}).done(function(msg) {
			console.log(msg);
			$('.pageing').html(msg.page);
			$('.slides').html(msg.data);
	 	}).fail(function(e) {
			console.log(e);
	}
	catch (e) {
		console.log(e);
	}
}