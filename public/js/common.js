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


	$("button.btn-warning").click(function() {
		var o = $(this).prevAll();
		$(o[0]).val("");
		$(o[1]).val("1").trigger("change");
		$(o[3]).val("3").trigger("change");
		$(o[5]).val("1").trigger("change");
		$(o[7]).val("");
	});

	$("#passenger").change(function() {
		th = $(this);
		if (th.val() !== "") {
			for (var i = 1; i <= 5; i++) {

				if ($("li.p" + i).children("input").first().val() == "") {
					var o = $("li.p" + i).children();
					$(o[0]).val(th.val());
					$(o[7]).val(th.val() + th.val());
					break;
				}
			};
		}
	});
})