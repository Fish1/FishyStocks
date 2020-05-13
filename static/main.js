function updateStock(ticker) {
	$.ajax({
		url: "./stocks/" + ticker + ".json",
		type: "POST",
		success: function(data) {
			console.log(data)
			$("#aapl_t").html("Time: " + data[0]["t"]);
			$("#aapl_c").html("Price: " + data[0]["c"]);
			$("#aapl_v").html("Volume: " + data[0]["v"]);
		}
	});
}

setInterval(() => {
	updateStock("AAPL");
}, 1000);
