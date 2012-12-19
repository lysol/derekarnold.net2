$(document).ready(function() {

		var ctx = document.getCSSCanvasContext('2d', 'butt', 800, 600);
		setInterval(function() {
				ctx.clearRect(0, 0, 800, 600);
				ctx.strokeStyle = 'black';
				ctx.lineWidth = '1';
				ctx.strokeRect(
						Math.random() * 800, 
						Math.random() * 600,
						Math.random() * 800, 
						Math.random() * 600
				);
				ctx.stroke();
		}, 1000);
});
		
