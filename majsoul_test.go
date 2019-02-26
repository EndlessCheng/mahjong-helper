package main

import (
	"testing"
	"strings"
	"strconv"
)

func TestParse(t *testing.T) {
	raw := "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 136 1 8 1 18 14 65 99 116 105 111 110 78 101 119 82 111 117 110 100 26 116 8 0 16 0 24 0 34 2 50 109 34 2 50 109 34 2 53 109 34 2 49 112 34 2 49 112 34 2 49 115 34 2 52 115 34 2 55 115 34 2 55 115 34 2 57 115 34 2 50 122 34 2 54 122 34 2 55 122 42 2 52 112 50 12 168 195 1 168 195 1 168 195 1 168 195 1 64 0 88 0 98 32 55 68 51 52 65 67 55 69 49 56 57 57 70 56 49 70 52 67 51 69 69 49 53 53 65 53 57 69 49 49 57 65 104 69]"
	//raw := "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 136 1 8 0 18 14 65 99 116 105 111 110 78 101 119 82 111 117 110 100 26 116 8 0 16 0 24 1 34 2 52 109 34 2 54 109 34 2 55 109 34 2 49 112 34 2 51 112 34 2 52 112 34 2 53 112 34 2 56 112 34 2 49 115 34 2 51 115 34 2 51 115 34 2 52 115 34 2 53 115 42 2 50 109 50 12 132 207 1 204 183 1 156 199 1 204 183 1 64 1 88 0 98 32 54 49 57 70 56 55 65 68 53 52 51 68 51 70 52 67 55 53 67 57 68 48 65 52 67 56 51 56 49 55 55 66 104 69]"

	//raw =  "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 155 1 8 0 18 14 65 99 116 105 111 110 78 101 119 82 111 117 110 100 26 134 1 8 1 16 0 24 3 34 2 49 109 34 2 52 109 34 2 52 109 34 2 53 109 34 2 56 109 34 2 50 112 34 2 52 112 34 2 48 112 34 2 55 112 34 2 56 112 34 2 49 115 34 2 50 115 34 2 53 115 34 2 55 115 42 2 52 122 50 12 188 180 1 160 156 1 156 249 1 168 195 1 58 12 8 0 18 2 8 1 32 0 40 152 236 3 64 0 88 0 98 32 67 56 52 50 70 56 66 57 51 70 65 70 69 67 65 57 51 70 55 48 57 68 53 48 53 51 50 48 48 65 51 70 104 69]"
	raw = "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 154 1 8 0 18 14 65 99 116 105 111 110 78 101 119 82 111 117 110 100 26 133 1 8 0 16 1 24 4 34 2 49 109 34 2 52 112 34 2 56 112 34 2 57 112 34 2 50 115 34 2 50 115 34 2 53 115 34 2 57 115 34 2 49 122 34 2 49 122 34 2 49 122 34 2 51 122 34 2 53 122 34 2 54 122 42 2 55 112 50 11 188 255 1 168 217 2 216 4 228 175 1 58 12 8 1 18 2 8 1 32 0 40 152 236 3 64 0 88 0 98 32 51 49 57 48 48 49 54 65 57 68 48 52 56 51 67 70 50 54 68 55 54 50 51 55 48 55 53 65 55 57 51 51 104 69]"

	// 明杠，舍牌
	raw = "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 45 8 10 18 17 65 99 116 105 111 110 68 105 115 99 97 114 100 84 105 108 101 26 22 8 1 18 2 49 109 24 0 40 0 48 0 66 2 55 112 66 2 53 109 72 0]"

	raw = "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 72 8 7 18 17 65 99 116 105 111 110 68 105 115 99 97 114 100 84 105 108 101 26 49 8 0 18 2 49 122 24 0 34 33 8 1 18 9 8 3 18 5 49 122 124 49 122 18 12 8 5 18 8 49 122 124 49 122 124 49 122 32 0 40 224 212 3 40 1 48 0 72 0]"

	// ActionChiPengGang
	raw = "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 64 8 54 18 17 65 99 116 105 111 110 67 104 105 80 101 110 103 71 97 110 103 26 41 8 1 16 1 26 2 49 122 26 2 49 122 26 2 49 122 34 3 1 1 3 50 16 8 1 18 6 8 1 18 2 49 122 32 0 40 224 212 3 56 0]"
	raw = "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 51 8 8 18 17 65 99 116 105 111 110 67 104 105 80 101 110 103 71 97 110 103 26 28 8 1 16 2 26 2 49 122 26 2 49 122 26 2 49 122 26 2 49 122 34 4 1 1 1 0 56 0]"

	// 摸牌
	raw = "[1 10 19 46 108 113 46 65 99 116 105 111 110 80 114 111 116 111 116 121 112 101 18 44 8 63 18 14 65 99 116 105 111 110 68 101 97 108 84 105 108 101 26 24 8 3 18 2 51 122 24 38 34 12 8 3 18 2 8 1 32 0 40 224 212 3 56 0]"

	for i, rawByte := range strings.Split(raw[1:len(raw)-1], " ") {
		b, err := strconv.Atoi(rawByte)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(i, b, string([]byte{byte(b)}))
	}

	t.Log(string(raw))
}
