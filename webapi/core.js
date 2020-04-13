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

function inline_img(line) {
	function img_tag(index) {
		return `<img src=${tile_img_url(index)} class="small-tile" />`;
	}

	let pattern = [/ ()东/, / ()南/, / ()西/, / ()北/,
		/ ()白/, / ()发/, / ()中/];
	let offset_p = [27, 28, 29, 30, 31, 32, 33];

	for (let i = 0; i < pattern.length; i++) {
		line = line.replace(pattern[i], function(_, j) {
			return img_tag(offset_p[i] + (parseInt(j) || 0))
		});
	}
	
	let offset = {
		'm': 0,
		'p': 9,
		's': 18,
		'z': 27,
		'万': 0,
		'饼': 9,
		'索': 18,
	};

	var keys = "";
	for (const [key, _] of Object.entries(offset)) {
		keys = keys + key;
	}

	let reg = new RegExp(`(\\d+)([${keys}])`, 'g');

	return line.replace(reg, function(_, s, t) {
		var result = "";
		for (var i = 0; i < s.length; i++) {
			result += img_tag(parseInt(s[i]) - 1 + offset[t]);
		}
		return result;
	});
}

function show_outputs(outputs) {
	let lines = outputs.replace(/\u001b\[\d*m/g, '').split('\n')
	let result = lines.map(function(line) {
		return `<p>${ inline_img(line) }</p>`;
	}).join('');
	document.getElementById("outputs").innerHTML = result;
}

// show_outputs("\u001b[96m39\u001b[0m\u001b[97m[47.04]\u001b[0m 切 南 => \u001b[96m22.95\u001b[0m两向听 [12345789m 7p 8s 37z]\n\u001b[96m39\u001b[0m\u001b[97m[47.04]\u001b[0m 切 西 => \u001b[96m22.95\u001b[0m两向听 [12345789m 7p 8s 27z]\n\u001b[96m31\u001b[0m[40.21] 切9万 => \u001b[96m20.10\u001b[0m两向听 [12345m 7p 8s 237z]\n\u001b[96m39\u001b[0m\u001b[97m[47.04]\u001b[0m 切 中 => \u001b[96m22.95\u001b[0m两向听 [12345789m 7p 8s 23z]\n\u001b[96m33\u001b[0m[39.99] 切3万 => \u001b[96m18.33\u001b[0m两向听 [14789m 7p 8s 237z]\n\u001b[96m28\u001b[0m[37.66] 切2万 => \u001b[96m14.00\u001b[0m两向听 [3789m 7p 8s 237z]\n\u001b[96m65\u001b[0m\u001b[97m[68.44]\u001b[0m 切9索 => \u001b[96m45.00\u001b[0m三向听 [12345789m 2347p 56789s 237z]\n\u001b[96m58\u001b[0m[63.13] 切7饼 => \u001b[96m44.83\u001b[0m三向听 [12345789m 56789p 8s 237z]\n\u001b[96m57\u001b[0m[61.50] 切7索 => \u001b[96m39.95\u001b[0m三向听 [12345789m 2347p 789s 237z]\n\u001b[96m61\u001b[0m[64.97] 切4饼 => \u001b[96m34.84\u001b[0m三向听 [12345789m 12347p 789s 237z]\n\u001b[96m61\u001b[0m[64.44] 切2饼 => \u001b[96m34.84\u001b[0m三向听 [12345789m 23457p 789s 237z]\n\u001b[96m57\u001b[0m[61.50] 切3饼 => \u001b[96m31.60\u001b[0m三向听 [12345789m 2347p 789s 237z]\n");

let interval = 1000;
var backoff = interval;

function update_tile() {
	$.getJSON("/api")
	.done(function(data) {
		show_tiles(data);
		show_outputs(data["outputs"]);
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
