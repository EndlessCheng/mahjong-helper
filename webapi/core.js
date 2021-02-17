"use strict";

let maj_risk = (function() {

// https://axonflux.com/handy-rgb-to-hsl-and-rgb-to-hsv-color-model-c

/**
 * Converts an HSL color value to RGB. Conversion formula
 * adapted from http://en.wikipedia.org/wiki/HSL_color_space.
 * Assumes h, s, and l are contained in the set [0, 1] and
 * returns r, g, and b in the set [0, 255].
 *
 * @param   {number}  h       The hue
 * @param   {number}  s       The saturation
 * @param   {number}  l       The lightness
 * @return  {Array}           The RGB representation
 */
function hslToRgb(h, s, l){
    var r, g, b;

    if(s == 0){
        r = g = b = l; // achromatic
    }else{
        var hue2rgb = function hue2rgb(p, q, t){
            if(t < 0) t += 1;
            if(t > 1) t -= 1;
            if(t < 1/6) return p + (q - p) * 6 * t;
            if(t < 1/2) return q;
            if(t < 2/3) return p + (q - p) * (2/3 - t) * 6;
            return p;
        }

        var q = l < 0.5 ? l * (1 + s) : l + s - l * s;
        var p = 2 * l - q;
        r = hue2rgb(p, q, h + 1/3);
        g = hue2rgb(p, q, h);
        b = hue2rgb(p, q, h - 1/3);
    }

    return [Math.round(r * 255), Math.round(g * 255), Math.round(b * 255)];
}

/**
 * Converts an RGB color value to HSL. Conversion formula
 * adapted from http://en.wikipedia.org/wiki/HSL_color_space.
 * Assumes r, g, and b are contained in the set [0, 255] and
 * returns h, s, and l in the set [0, 1].
 *
 * @param   {number}  r       The red color value
 * @param   {number}  g       The green color value
 * @param   {number}  b       The blue color value
 * @return  {Array}           The HSL representation
 */
function rgbToHsl(r, g, b){
    r /= 255, g /= 255, b /= 255;
    var max = Math.max(r, g, b), min = Math.min(r, g, b);
    var h, s, l = (max + min) / 2;

    if(max == min){
        h = s = 0; // achromatic
    }else{
        var d = max - min;
        s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
        switch(max){
            case r: h = (g - b) / d + (g < b ? 6 : 0); break;
            case g: h = (b - r) / d + 2; break;
            case b: h = (r - g) / d + 4; break;
        }
        h /= 6;
    }

    return [h, s, l];
}

let gradient = [
	{ risk: 0.0 , color: [255, 255,  255] }, // 绝安，白色
	{ risk: 1e-5 , color: [131, 179,  17] },  // 只要有一点风险，鸭绿色
	{ risk: 8 , color:   [250, 195,  0] },  // 还可以冲一冲， 橙黄色
	{ risk: 16 , color: [253, 99, 0] },  // 相当危险，粉色
	{ risk: 25.0, color: [253, 0,  135] }, // 绝对危险， 红色
];

// https://stackoverflow.com/questions/4856717/javascript-equivalent-of-pythons-zip-function
let zip = rows=>rows[0].map((_,c)=>rows.map(row=>row[c]));

return {
	get_rist_color(risk) {
		risk = Math.min(risk, gradient[gradient.length-1].risk);

		for (let i = 0; i + 1 < gradient.length; i++) {
			if (risk <= gradient[i + 1].risk) {
				// let left = rgbToHsl.apply(null, gradient[i].color);
				// let right = rgbToHsl.apply(null, gradient[i + 1].color);
				let left = gradient[i].color;
				let right = gradient[i+1].color;
				// interpolation in RGB color space
				// HSL 效果似乎并不如人意
				let ratio = (risk - gradient[i].risk) / (gradient[i + 1].risk - gradient[i].risk);
				let result = zip([left, right]).map(([x, y]) => x + (y - x) * ratio);

				// let rgba = hslToRgb.apply(null, (result)).join(',');
				let rgba = result;
				return `rgba(${rgba})`;
			}
		}

		throw "can not find a color for the risk " + risk;
	}
};

})(); // closure for maj_risk

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
	    
		let count = data["counts"][i];
		for (let c = 0; c < count; c++) {
			let img = document.createElement("img");
			img.src = tile_img_url(i);
			img.classList.add("tile");
			img.style.border = `solid 3px ${maj_risk.get_rist_color(risk)}`;
			container.appendChild(img);
		}
	}
}

function inline_img(line) {
	function img_tag(index) {
		return `<img src=${tile_img_url(index)} class="small-tile" />`;
	}

	let z_tiles = "东南西北白发中";
	
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

	let reg = new RegExp(`(\\d+)([${keys}])| ([${z_tiles}])`, 'g');

	return line.replace(reg, function(_, s, t, z) {
		if (z) {
			return img_tag(27 + z_tiles.indexOf(z));
		}

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
