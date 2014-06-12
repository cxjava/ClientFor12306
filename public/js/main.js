$(document).ready(function() {

	$("#img").click(function(event) {
		$(this).attr('src', '/submitPassCodeNew/' + Math.random());
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
				var v = passenger[th.val()]
				vals = $("li.p" + i).children("input").first().val();
				if (vals == v.passenger_name) {
					break;
				}
				if (vals == "") {
					var o = $("li.p" + i).children();
					console.log($(o[0]).val());

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
	var c = new WebSocket('ws://localhost:3000/sock');
	$("#code").keyup(function() {
		if ($(this).val().length == 4) {
			c.send("code#" + $(this).val());
		}
	});

	c.onopen = function() {
		c.onmessage = function(response) {
			console.log(response.data);
			if (response.data = "update") {
				$("#imageDiv").html('<img src="/submitPassCodeNew/' + Math.random() + '" id="img" title="单击刷新验证码">');
				$("#code").val("").focus();
			}
		};
		c.send("test");
	}

	$("#submit").click(function() {

		$.post('/query', $("form").serialize(), function(data, textStatus, xhr) {
			console.log(data);
		});
		return false;
	});

	$("#imageDiv").click(function(event) {
		$("#img").attr('src', '/submitPassCodeNew/' + Math.random());
		$("#code").val("").focus();
	});
});