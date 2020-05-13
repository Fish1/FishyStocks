var ctx = $('#appl_chart');

var appl_chart = new Chart(ctx, {
	type: 'line',
	data: {
		labels: ['Label_A', 'Label_B', 'Label_C', 'Label_D'],
		datasets: [{
			label: '# of votes',
			data: [12, 19, 3, 5]
		}]
	},
	options: {
		responsive: false
	}
});
