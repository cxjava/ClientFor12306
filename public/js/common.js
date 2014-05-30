$(document).ready(function() {

	$("#img").click(function(event) {
		$(this).attr('src', '/loginPassCodeNew/' + Math.random());
	});

	$('.form_datetime').datetimepicker({
		endDate: "+19d",
		startDate: "+0d",
		initialDate: "+19d",
		format: 'yyyy-mm-dd',
		todayBtn: 1,
		autoclose: 1,
		todayHighlight: 1,
		minView: 2,
		startView: 2,
		autoclose: true,
		language: 'zh-CN'
	});

	$("select").select2({
		language: "zh_CN",
		allowClear: true
	});
})