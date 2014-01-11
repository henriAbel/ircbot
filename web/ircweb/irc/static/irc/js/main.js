var initLazy = function() {
	$(".lazy img").lazyload({
	    effect : "fadeIn"
	});
}
$(function() {
	initLazy();
});