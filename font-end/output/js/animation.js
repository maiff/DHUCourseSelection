//tool.js
$.fn.isOnScreen = function(){
     
    var win = $(window);
     
    var viewport = {
        top : win.scrollTop(),
        left : win.scrollLeft()
    };
    viewport.right = viewport.left + win.width();
    viewport.bottom = viewport.top + win.height();
     
    var bounds = this.offset();
    bounds.right = bounds.left + this.outerWidth(true);
    bounds.bottom = bounds.top + this.outerHeight(true);
     
    return (!(viewport.right < bounds.left || viewport.left > bounds.right || viewport.bottom < bounds.top || viewport.top > bounds.bottom));
     
};

//end
$('#navmenu').on('mouseover','li',function(){
	$(this).find('em').removeClass('wi0');
});
$('#navmenu').on('mouseout','li',function(){
	$(this).find('em').addClass('wi0');
});

$(window).scroll(function(){
	if($(window).scrollTop()>$('.table-container').offset().top-55){
		$('#table-container').addClass('show');
	}
	else{
		$('#table-container').removeClass('show');
	}
});



if($('footer').isOnScreen()&&$(window).scrollTop()===0){
	$('footer').addClass('footer-fixed');
}