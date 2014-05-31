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
					var v = passenger[th.val()]
					var o = $("li.p" + i).children();
					$(o[0]).val(v.passenger_name);
					$(o[2]).val(v.passenger_type).trigger("change");
					$(o[4]).val("3").trigger("change");
					$(o[6]).val(v.passenger_id_type_code).trigger("change");
					$(o[7]).val(v.passenger_id_no);
					break;
				}
			};
		}
	});

	$.post('/loadUser', function(data, textStatus, xhr) {
		if (data.r == true) {
			var p = data.o,
				str = "<option></option>";
			passenger = {}
			for (var i = 0; i < p.length; i++) {
				passenger[p[i].passenger_id_no] = p[i];
				str = str + '<option value="' + p[i].passenger_id_no + '">' + p[i].passenger_name + '</option>'
			};
			$('#passenger').html(str).trigger('change')
		}
	});

})