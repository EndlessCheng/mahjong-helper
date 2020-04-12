"use strict";

function tile_img_url(index) {
	function label(index) {
		let offset = [0, 9, 18, 27, 34];
		let l = ["m", "p", "s", "z"];

		for (let i = 0; i < l.length; i++) {
			if (index < offset[i+1]) {
				return (index - offset[i] + 1) + l[i];
			}
		}
		throw new Error("invalid index");
	}

	let l = label(index);
	return `img/${l}.png`;
}

function show_tiles(data) {
	if (data["counts"] === null) {
		return;
	}

	let container = document.getElementById("my-tiles");
	container.innerHTML = '';

	for (let i = 0; i < 34; i++) {
		var risk = data["risk"] === null ? 0 : data["risk"][i];
	    risk = risk / 25;
		risk = Math.min(risk, 1.0);
		let count = data["counts"][i];
		for (let c = 0; c < count; c++) {
			let img = document.createElement("img");
			img.src = tile_img_url(i);
			img.classList.add("tile");
			img.style.border = `solid 3px rgba(255, 0, 0, ${risk})`;
			container.appendChild(img);
		}
	}
}

let interval = 1000;
var backoff = interval;

function update_tile() {
	$.getJSON("/api")
	.done(function(data) {
		show_tiles(data);
		backoff = interval;
		setTimeout(update_tile, interval);
	})
	.fail(function( jqxhr, textStatus, error ) {
		var err = textStatus + ", " + error;
		console.log("Request Failed: " + err);
		backoff = backoff * 2;
		setTimeout(update_tile, backoff);
	});
}

setTimeout(update_tile, interval);
