$(document).ready(function() {

	$("#img").click(function(event) {
		$(this).attr('src', '/loginPassCodeNew/' + Math.random());
		$("#code").val("").focus();
	});

})