(function($) {
	var IS_IOS = /iphone|ipad/i.test(navigator.userAgent);
	$.fn.nodoubletap = function() {
	if (IS_IOS)
		$(this).bind('touchstart', function preventDoubleTap(e) {
		var t2 = e.timeStamp
			, t1 = $(this).data('lastTouch') || t2
			, dt = t2 - t1
			, fingers = e.originalEvent.touches.length;
		$(this).data('lastTouch', t2);
		if (!dt || dt > 500 || fingers > 1) return; // not double-tap

		e.preventDefault();
		$(this).trigger("touchstart").trigger("touchend").trigger("click");

		});
	};
})(jQuery);

