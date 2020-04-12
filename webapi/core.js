"use strict";

function tile_img_url(index) {
	function label(index) {
		let offset = [0, 9, 18, 27, 34];
		let l = ["m", "s", "p", "z"];

		for (let i = 0; i < l.length; i++) {
			if (index < offset[i+1]) {
				return (index - offset[i]) + l[i];
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

var data = {
    "counts": [
        0,
        0,
        0,
        0,
        1,
        0,
        1,
        1,
        1,
        0,
        0,
        0,
        1,
        1,
        1,
        1,
        1,
        1,
        0,
        0,
        0,
        0,
        0,
        0,
        1,
        0,
        1,
        0,
        2,
        0,
        0,
        0,
        0,
        0
    ],
    "risk": [
        6.8689162795508745,
        1.3130278422273782,
        14.364682863109048,
        26.086178013431933,
        0,
        3.996171693735499,
        26.609623463447978,
        10.24728615971008,
        4.707313155452437,
        1.62,
        12.041053092734451,
        4.955252900232019,
        4.019039743204979,
        31.297173259860788,
        22.897446084686777,
        13.089059187935034,
        21.93600580280426,
        20.047425288863106,
        2.1125103680057706,
        4.555635730858469,
        23.916679946635732,
        18.880975342227377,
        22.090670997703363,
        35.24169806753138,
        21.933339973317867,
        4.555635730858469,
        20.978551860920128,
        0,
        1.04143593387471,
        0,
        0.8999999999999999,
        2.519098836845261,
        0.14272041763341067,
        17.777878904753784
    ],
    "timestamp": 0
}

show_tiles(data)