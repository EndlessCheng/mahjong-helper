// Copyrights C-EGG inc.
(function () {
    var u = function () {
        function b(a) {
            var b = a & 7, c = 0, d = 0;
            1 == b || 4 == b ? c = d = 1 : 2 == b && (c = d = 2);
            a >>= 3;
            b = (a & 7) - c;
            if (0 > b) return !1;
            c = d;
            d = 0;
            1 == b || 4 == b ? (c += 1, d += 1) : 2 == b && (c += 2, d += 2);
            a >>= 3;
            b = (a & 7) - c;
            if (0 > b) return !1;
            c = d;
            d = 0;
            1 == b || 4 == b ? (c += 1, d += 1) : 2 == b && (c += 2, d += 2);
            a >>= 3;
            b = (a & 7) - c;
            if (0 > b) return !1;
            c = d;
            d = 0;
            1 == b || 4 == b ? (c += 1, d += 1) : 2 == b && (c += 2, d += 2);
            a >>= 3;
            b = (a & 7) - c;
            if (0 > b) return !1;
            c = d;
            d = 0;
            1 == b || 4 == b ? (c += 1, d += 1) : 2 == b && (c += 2, d += 2);
            a >>= 3;
            b = (a & 7) - c;
            if (0 > b) return !1;
            c = d;
            d = 0;
            1 == b || 4 == b ? (c += 1, d += 1) : 2 == b && (c += 2, d += 2);
            a >>= 3;
            b = (a & 7) - c;
            if (0 > b) return !1;
            c = d;
            d = 0;
            1 == b || 4 == b ? (c += 1, d += 1) : 2 == b && (c += 2, d += 2);
            a >>= 3;
            b = (a & 7) - c;
            if (0 != b && 3 != b) return !1;
            b = (a >> 3 & 7) - d;
            return 0 == b || 3 == b
        }

        function a(a, d) {
            if (0 == a) {
                if (128 <= (d & 448) && b(d - 128) || 65536 <= (d & 229376) && b(d - 65536) || 33554432 <= (d & 117440512) && b(d - 33554432)) return !0
            } else if (1 == a) {
                if (16 <= (d & 56) && b(d - 16) || 8192 <= (d & 28672) && b(d - 8192) || 4194304 <= (d & 14680064) && b(d - 4194304)) return !0
            } else if (2 == a && (2 <= (d & 7) && b(d - 2) || 1024 <= (d & 3584) && b(d - 1024) || 524288 <= (d & 1835008) && b(d - 524288))) return !0;
            return !1
        }

        function g(a,
                   b) {
            return a[b + 0] << 0 | a[b + 1] << 3 | a[b + 2] << 6 | a[b + 3] << 9 | a[b + 4] << 12 | a[b + 5] << 15 | a[b + 6] << 18 | a[b + 7] << 21 | a[b + 8] << 24
        }

        function d(c) {
            var d = 1 << c[27] | 1 << c[28] | 1 << c[29] | 1 << c[30] | 1 << c[31] | 1 << c[32] | 1 << c[33];
            if (16 <= d) return !1;
            if (2 == (d & 3) && 2 == c[0] * c[8] * c[9] * c[17] * c[18] * c[26] * c[27] * c[28] * c[29] * c[30] * c[31] * c[32] * c[33] || !(d & 10) && 7 == (2 == c[0]) + (2 == c[1]) + (2 == c[2]) + (2 == c[3]) + (2 == c[4]) + (2 == c[5]) + (2 == c[6]) + (2 == c[7]) + (2 == c[8]) + (2 == c[9]) + (2 == c[10]) + (2 == c[11]) + (2 == c[12]) + (2 == c[13]) + (2 == c[14]) + (2 == c[15]) + (2 == c[16]) + (2 == c[17]) + (2 ==
                c[18]) + (2 == c[19]) + (2 == c[20]) + (2 == c[21]) + (2 == c[22]) + (2 == c[23]) + (2 == c[24]) + (2 == c[25]) + (2 == c[26]) + (2 == c[27]) + (2 == c[28]) + (2 == c[29]) + (2 == c[30]) + (2 == c[31]) + (2 == c[32]) + (2 == c[33])) return !0;
            if (d & 2) return !1;
            var r = c[0] + c[3] + c[6], m = c[1] + c[4] + c[7], n = c[9] + c[12] + c[15], e = c[10] + c[13] + c[16],
                q = c[18] + c[21] + c[24], k = c[19] + c[22] + c[25], p = (r + m + (c[2] + c[5] + c[8])) % 3;
            if (1 == p) return !1;
            var l = (n + e + (c[11] + c[14] + c[17])) % 3;
            if (1 == l) return !1;
            var t = (q + k + (c[20] + c[23] + c[26])) % 3;
            if (1 == t || 1 != (2 == p) + (2 == l) + (2 == t) + (2 == c[27]) + (2 == c[28]) + (2 ==
                c[29]) + (2 == c[30]) + (2 == c[31]) + (2 == c[32]) + (2 == c[33])) return !1;
            r = (1 * r + 2 * m) % 3;
            m = g(c, 0);
            n = (1 * n + 2 * e) % 3;
            e = g(c, 9);
            q = (1 * q + 2 * k) % 3;
            c = g(c, 18);
            return d & 4 ? !(p | r | l | n | t | q) && b(m) && b(e) && b(c) : 2 == p ? !(l | n | t | q) && b(e) && b(c) && a(r, m) : 2 == l ? !(t | q | p | r) && b(c) && b(m) && a(n, e) : 2 == t ? !(p | r | l | n) && b(m) && b(e) && a(q, c) : !1
        }

        return function (a, b) {
            if (34 == b) return d(a)
        }
    }();

    function w() {
        this.h = [-1, -1, -1, -1, -1, -1, -1];
        this.c = [{b: -1, a: 0}, {b: -1, a: 0}, {b: -1, a: 0}, {b: -1, a: 0}]
    }

    w.prototype = {};

    function x(b, a, g, d) {
        b = b.c;
        var c = b[0].a, f = [0, 0, 0], r = 7 << 24 - 3 * a, m = 2 << 24 - 3 * a, n = 0;
        (d & r) >= m && y(c, g, d - m, f) && (f[0] && (b[n].b = 9 * g + 8 - a, b[n].a = f[0], ++n), f[1] && (b[n].b = 9 * g + 8 - a, b[n].a = f[1], ++n), f[2] && (b[n].b = 9 * g + 8 - a, b[n].a = f[2], ++n));
        r >>= 9;
        m >>= 9;
        (d & r) >= m && y(c, g, d - m, f) && (f[0] && (b[n].b = 9 * g + 5 - a, b[n].a = f[0], ++n), f[1] && (b[n].b = 9 * g + 5 - a, b[n].a = f[1], ++n), f[2] && (b[n].b = 9 * g + 5 - a, b[n].a = f[2], ++n));
        m >>= 9;
        (d & r >> 9) >= m && y(c, g, d - m, f) && (f[0] && (b[n].b = 9 * g + 2 - a, b[n].a = f[0], ++n), f[1] && (b[n].b = 9 * g + 2 - a, b[n].a = f[1], ++n), f[2] && (b[n].b =
            9 * g + 2 - a, b[n].a = f[2], ++n));
        return 0 != n
    }

    function z(b, a, g) {
        b = b.c;
        var d = [0, 0, 0];
        if (!y(b[0].a, a, g, d)) return !1;
        a = 0;
        d[0] && (b[a].b = b[0].b, b[a].a = d[0], ++a);
        d[1] && (b[a].b = b[0].b, b[a].a = d[1], ++a);
        d[2] && (b[a].b = b[0].b, b[a].a = d[2], ++a);
        return 0 != a
    }

    function y(b, a, g, d) {
        var c = -1, f, r = g & 7, m = 0, n = 0;
        for (f = 0; 7 > f && 1755 != g; ++f) {
            switch (r) {
                case 4:
                    b <<= 8, b |= 7 * a + f + 1, m += 1, n += 1;
                case 3:
                    (g >> 3 & 7) >= 3 + m && (g >> 6 & 7) >= 3 + n ? (c = f, m += 3, n += 3) : (b <<= 8, b |= 21 + 9 * a + f + 1);
                    break;
                case 2:
                    b <<= 16;
                    b |= 257 * (7 * a + f + 1);
                    m += 2;
                    n += 2;
                    break;
                case 1:
                    b <<= 8;
                    b |= 7 * a + f + 1;
                    m += 1;
                    n += 1;
                    break;
                case 0:
                    break;
                default:
                    return 0
            }
            g >>= 3;
            r = (g & 7) - m;
            m = n;
            n = 0
        }
        if (7 > f) return d[0] = 16843009 * (21 + 9 * a + f + 1) + 66051, d[1] = 65793 * (7 * a + f + 1 + 1) | 21 + 9 * a + f + 0 + 1 << 24, d[2] = 65793 * (7 * a + f + 0 + 1) | 21 + 9 * a + f + 3 + 1 << 24, 3;
        if (3 == r) b = b << 8 | 9 * a + 29; else if (r) return 0;
        r = (g >> 3 & 7) - m;
        if (3 == r) b = b << 8 | 9 * a + 30; else if (r) return 0;
        if (-1 != c) return b <<= 24, d[0] = b | 65793 * (21 + 9 * a + c + 1) + 258, d[1] = b | 65793 * (7 * a + c + 1), d[2] = 0, 2;
        d[0] = b;
        d[1] = d[2] = 0;
        return 1
    }

    function A(b, a, g, d) {
        var c = 7 << 24 - 3 * a, f = 2 << 24 - 3 * a;
        if ((d & c) >= f && B(b, g, d - f)) return b.c[0].b = 9 * g + 8 - a, !0;
        c >>= 9;
        f >>= 9;
        if ((d & c) >= f && B(b, g, d - f)) return b.c[0].b = 9 * g + 5 - a, !0;
        f >>= 9;
        return (d & c >> 9) >= f && B(b, g, d - f) ? (b.c[0].b = 9 * g + 2 - a, !0) : !1
    }

    function B(b, a, g) {
        var d = b.c[0].a, c, f = g & 7, r = 0, m = 0;
        for (c = 0; 7 > c; ++c) {
            switch (f) {
                case 4:
                    d <<= 16;
                    d |= 21 + 9 * a + c + 1 << 8 | 7 * a + c + 1;
                    r += 1;
                    m += 1;
                    break;
                case 3:
                    d <<= 8;
                    d |= 21 + 9 * a + c + 1;
                    break;
                case 2:
                    d <<= 16;
                    d |= 257 * (7 * a + c + 1);
                    r += 2;
                    m += 2;
                    break;
                case 1:
                    d <<= 8;
                    d |= 7 * a + c + 1;
                    r += 1;
                    m += 1;
                    break;
                case 0:
                    break;
                default:
                    return !1
            }
            g >>= 3;
            f = (g & 7) - r;
            r = m;
            m = 0
        }
        if (3 == f) d = d << 8 | 9 * a + 29; else if (f) return !1;
        f = (g >> 3 & 7) - r;
        if (3 == f) d = d << 8 | 9 * a + 30; else if (f) return !1;
        b.c[0].a = d;
        return !0
    }

    function C(b, a) {
        var g, d = b.c, c = 1 << a[27] | 1 << a[28] | 1 << a[29] | 1 << a[30] | 1 << a[31] | 1 << a[32] | 1 << a[33];
        if (16 <= c) return !1;
        if (2 == (c & 3) && 2 == a[0] * a[8] * a[9] * a[17] * a[18] * a[26] * a[27] * a[28] * a[29] * a[30] * a[31] * a[32] * a[33]) {
            var f, c = [0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33];
            for (f = 0; 13 > f && 2 != a[c[f]]; ++f) ;
            d[0].b = c[f];
            d[0].a = 4294967295;
            return !0
        }
        if (c & 2) return !1;
        f = !1;
        if (!(c & 10) && 7 == (2 == a[0]) + (2 == a[1]) + (2 == a[2]) + (2 == a[3]) + (2 == a[4]) + (2 == a[5]) + (2 == a[6]) + (2 == a[7]) + (2 == a[8]) + (2 == a[9]) + (2 == a[10]) + (2 == a[11]) + (2 == a[12]) + (2 == a[13]) + (2 ==
            a[14]) + (2 == a[15]) + (2 == a[16]) + (2 == a[17]) + (2 == a[18]) + (2 == a[19]) + (2 == a[20]) + (2 == a[21]) + (2 == a[22]) + (2 == a[23]) + (2 == a[24]) + (2 == a[25]) + (2 == a[26]) + (2 == a[27]) + (2 == a[28]) + (2 == a[29]) + (2 == a[30]) + (2 == a[31]) + (2 == a[32]) + (2 == a[33])) {
            d[3].a = 4294967295;
            for (f = g = 0; 34 > f; ++f) 2 == a[f] && (b.h[g] = f, g += 1);
            f = !0
        }
        var r = a[0] + a[3] + a[6], m = a[1] + a[4] + a[7], n = a[2] + a[5] + a[8], e = a[9] + a[12] + a[15],
            q = a[10] + a[13] + a[16], k = a[11] + a[14] + a[17], p = a[18] + a[21] + a[24], l = a[19] + a[22] + a[25],
            t = a[20] + a[23] + a[26];
        g = (r + m + n) % 3;
        if (1 == g) return f;
        var v = (e + q + k) %
            3;
        if (1 == v) return f;
        var h = (p + l + t) % 3;
        if (1 == h || 1 != (2 == g) + (2 == v) + (2 == h) + (2 == a[27]) + (2 == a[28]) + (2 == a[29]) + (2 == a[30]) + (2 == a[31]) + (2 == a[32]) + (2 == a[33])) return f;
        c & 8 && (3 == a[27] && (d[0].a <<= 8, d[0].a |= 49), 3 == a[28] && (d[0].a <<= 8, d[0].a |= 50), 3 == a[29] && (d[0].a <<= 8, d[0].a |= 51), 3 == a[30] && (d[0].a <<= 8, d[0].a |= 52), 3 == a[31] && (d[0].a <<= 8, d[0].a |= 53), 3 == a[32] && (d[0].a <<= 8, d[0].a |= 54), 3 == a[33] && (d[0].a <<= 8, d[0].a |= 55));
        n = r + m + n;
        r = (1 * r + 2 * m) % 3;
        m = D(a, 0);
        k = e + q + k;
        e = (1 * e + 2 * q) % 3;
        q = D(a, 9);
        t = p + l + t;
        p = (1 * p + 2 * l) % 3;
        l = D(a, 18);
        if (c & 4) {
            if (g |
                r | v | e | h | p) return f;
            2 == a[27] ? d[0].b = 27 : 2 == a[28] ? d[0].b = 28 : 2 == a[29] ? d[0].b = 29 : 2 == a[30] ? d[0].b = 30 : 2 == a[31] ? d[0].b = 31 : 2 == a[32] ? d[0].b = 32 : 2 == a[33] && (d[0].b = 33);
            if (9 <= n) {
                if (B(b, 1, q) && B(b, 2, l) && z(b, 0, m)) return !0
            } else if (9 <= k) {
                if (B(b, 2, l) && B(b, 0, m) && z(b, 1, q)) return !0
            } else if (9 <= t) {
                if (B(b, 0, m) && B(b, 1, q) && z(b, 2, l)) return !0
            } else if (B(b, 0, m) && B(b, 1, q) && B(b, 2, l)) return !0
        } else if (2 == g) {
            if (v | e | h | p) return f;
            if (8 <= n) {
                if (B(b, 1, q) && B(b, 2, l) && x(b, r, 0, m)) return !0
            } else if (9 <= k) {
                if (B(b, 2, l) && A(b, r, 0, m) && z(b, 1, q)) return !0
            } else if (9 <=
                t) {
                if (A(b, r, 0, m) && B(b, 1, q) && z(b, 2, l)) return !0
            } else if (B(b, 1, q) && B(b, 2, l) && A(b, r, 0, m)) return !0
        } else if (2 == v) {
            if (h | p | g | r) return f;
            if (8 <= k) {
                if (B(b, 2, l) && B(b, 0, m) && x(b, e, 1, q)) return !0
            } else if (9 <= t) {
                if (B(b, 0, m) && A(b, e, 1, q) && z(b, 2, l)) return !0
            } else if (9 <= n) {
                if (A(b, e, 1, q) && B(b, 2, l) && z(b, 0, m)) return !0
            } else if (B(b, 2, l) && B(b, 0, m) && A(b, e, 1, q)) return !0
        } else if (2 == h) {
            if (g | r | v | e) return f;
            if (8 <= t) {
                if (B(b, 0, m) && B(b, 1, q) && x(b, p, 2, l)) return !0
            } else if (9 <= n) {
                if (B(b, 1, q) && A(b, p, 2, l) && z(b, 0, m)) return !0
            } else if (9 <= k) {
                if (A(b,
                    p, 2, l) && B(b, 0, m) && z(b, 1, q)) return !0
            } else if (B(b, 0, m) && B(b, 1, q) && A(b, p, 2, l)) return !0
        }
        d[0].a = 0;
        return f
    }

    function D(b, a) {
        return b[a + 0] << 0 | b[a + 1] << 3 | b[a + 2] << 6 | b[a + 3] << 9 | b[a + 4] << 12 | b[a + 5] << 15 | b[a + 6] << 18 | b[a + 7] << 21 | b[a + 8] << 24
    };var E = function () {
        function b(a) {
            e[a] -= 2;
            ++p
        }

        function a(a) {
            e[a] += 2;
            --p
        }

        function g(a) {
            --e[a];
            --e[a + 1];
            --e[a + 2];
            ++q
        }

        function d(a) {
            ++e[a];
            ++e[a + 1];
            ++e[a + 2];
            --q
        }

        function c(a) {
            --e[a];
            --e[a + 1];
            ++k
        }

        function f(a) {
            ++e[a];
            ++e[a + 1];
            --k
        }

        function r(a) {
            --e[a];
            --e[a + 2];
            ++k
        }

        function m(a) {
            ++e[a];
            ++e[a + 2];
            --k
        }

        var n = 0, e, q = 0, k = 0, p = 0, l = 0, t = 0, v = 0;
        return {
            g: 8, v: function () {
                var a = 8 - 2 * q - k - p, b = q + k;
                p ? b += p - 1 : t && v && (t | v) == t && ++a;
                4 < b && (a += b - 4);
                -1 != a && a < l && (a = l);
                a < this.g && (this.g = a)
            }, l: function (a, b) {
                e = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
                    0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
                v = t = l = p = k = q = 0;
                this.g = 8;
                if (136 == b) for (b = 0; 136 > b; ++b) a[b] && ++e[b >> 2]; else if (34 == b) for (b = 0; 34 > b; ++b) e[b] = a[b]; else for (--b; 0 <= b; --b) ++e[a[b] >> 2]
            }, j: function () {
                return e[0] + e[1] + e[2] + e[3] + e[4] + e[5] + e[6] + e[7] + e[8] + e[9] + e[10] + e[11] + e[12] + e[13] + e[14] + e[15] + e[16] + e[17] + e[18] + e[19] + e[20] + e[21] + e[22] + e[23] + e[24] + e[25] + e[26] + e[27] + e[28] + e[29] + e[30] + e[31] + e[32] + e[33]
            }, o: function () {
                var a = (2 <= e[0]) + (2 <= e[8]) + (2 <= e[9]) + (2 <= e[17]) + (2 <= e[18]) + (2 <= e[26]) + (2 <= e[27]) + (2 <= e[28]) +
                    (2 <= e[29]) + (2 <= e[30]) + (2 <= e[31]) + (2 <= e[32]) + (2 <= e[33]),
                    b = (0 != e[0]) + (0 != e[8]) + (0 != e[9]) + (0 != e[17]) + (0 != e[18]) + (0 != e[26]) + (0 != e[27]) + (0 != e[28]) + (0 != e[29]) + (0 != e[30]) + (0 != e[31]) + (0 != e[32]) + (0 != e[33]),
                    d = b + (0 != e[1]) + (0 != e[2]) + (0 != e[3]) + (0 != e[4]) + (0 != e[5]) + (0 != e[6]) + (0 != e[7]) + (0 != e[10]) + (0 != e[11]) + (0 != e[12]) + (0 != e[13]) + (0 != e[14]) + (0 != e[15]) + (0 != e[16]) + (0 != e[19]) + (0 != e[20]) + (0 != e[21]) + (0 != e[22]) + (0 != e[23]) + (0 != e[24]) + (0 != e[25]),
                    c = this.g,
                    d = 6 - (a + (2 <= e[1]) + (2 <= e[2]) + (2 <= e[3]) + (2 <= e[4]) + (2 <= e[5]) + (2 <= e[6]) +
                        (2 <= e[7]) + (2 <= e[10]) + (2 <= e[11]) + (2 <= e[12]) + (2 <= e[13]) + (2 <= e[14]) + (2 <= e[15]) + (2 <= e[16]) + (2 <= e[19]) + (2 <= e[20]) + (2 <= e[21]) + (2 <= e[22]) + (2 <= e[23]) + (2 <= e[24]) + (2 <= e[25])) + (7 > d ? 7 - d : 0);
                d < c && (c = d);
                d = 13 - b - (a ? 1 : 0);
                d < c && (c = d);
                return c
            }, m: function (a) {
                var b = 0, d = 0, c;
                for (c = 27; 34 > c; ++c) switch (e[c]) {
                    case 4:
                        ++q;
                        b |= 1 << c - 27;
                        d |= 1 << c - 27;
                        ++l;
                        break;
                    case 3:
                        ++q;
                        break;
                    case 2:
                        ++p;
                        break;
                    case 1:
                        d |= 1 << c - 27
                }
                l && 2 == a % 3 && --l;
                d && (v |= 134217728, (b | d) == b && (t |= 134217728))
            }, w: function (a) {
                var b = 0, d = 0, c;
                for (c = 27; 34 > c; ++c) switch (e[c]) {
                    case 4:
                        ++q;
                        b |= 1 << c - 18;
                        d |= 1 << c - 18;
                        ++l;
                        break;
                    case 3:
                        ++q;
                        break;
                    case 2:
                        ++p;
                        break;
                    case 1:
                        d |= 1 << c - 18
                }
                for (c = 0; 9 > c; c += 8) switch (e[c]) {
                    case 4:
                        ++q;
                        b |= 1 << c;
                        d |= 1 << c;
                        ++l;
                        break;
                    case 3:
                        ++q;
                        break;
                    case 2:
                        ++p;
                        break;
                    case 1:
                        d |= 1 << c
                }
                l && 2 == a % 3 && --l;
                d && (v |= 134217728, (b | d) == b && (t |= 134217728))
            }, s: function (a) {
                t |= (4 == e[0]) << 0 | (4 == e[1]) << 1 | (4 == e[2]) << 2 | (4 == e[3]) << 3 | (4 == e[4]) << 4 | (4 == e[5]) << 5 | (4 == e[6]) << 6 | (4 == e[7]) << 7 | (4 == e[8]) << 8 | (4 == e[9]) << 9 | (4 == e[10]) << 10 | (4 == e[11]) << 11 | (4 == e[12]) << 12 | (4 == e[13]) << 13 | (4 == e[14]) << 14 | (4 == e[15]) << 15 | (4 == e[16]) <<
                    16 | (4 == e[17]) << 17 | (4 == e[18]) << 18 | (4 == e[19]) << 19 | (4 == e[20]) << 20 | (4 == e[21]) << 21 | (4 == e[22]) << 22 | (4 == e[23]) << 23 | (4 == e[24]) << 24 | (4 == e[25]) << 25 | (4 == e[26]) << 26;
                q += a;
                this.u(0)
            }, u: function (h) {
                var k = arguments.callee;
                ++n;
                if (-1 != this.g) {
                    for (; 27 > h && !e[h]; ++h) ;
                    if (27 == h) return this.v();
                    var l = h;
                    8 < l && (l -= 9);
                    8 < l && (l -= 9);
                    switch (e[h]) {
                        case 4:
                            e[h] -= 3;
                            ++q;
                            7 > l && e[h + 2] && (e[h + 1] && (g(h), k.call(this, h + 1), d(h)), r(h), k.call(this, h + 1), m(h));
                            8 > l && e[h + 1] && (c(h), k.call(this, h + 1), f(h));
                            var p = h;
                            --e[p];
                            v |= 1 << p;
                            k.call(this, h + 1);
                            p = h;
                            ++e[p];
                            v &= ~(1 << p);
                            e[h] += 3;
                            --q;
                            b(h);
                            7 > l && e[h + 2] && (e[h + 1] && (g(h), k.call(this, h), d(h)), r(h), k.call(this, h + 1), m(h));
                            8 > l && e[h + 1] && (c(h), k.call(this, h + 1), f(h));
                            a(h);
                            break;
                        case 3:
                            e[h] -= 3;
                            ++q;
                            k.call(this, h + 1);
                            e[h] += 3;
                            --q;
                            b(h);
                            7 > l && e[h + 1] && e[h + 2] ? (g(h), k.call(this, h + 1), d(h)) : (7 > l && e[h + 2] && (r(h), k.call(this, h + 1), m(h)), 8 > l && e[h + 1] && (c(h), k.call(this, h + 1), f(h)));
                            a(h);
                            7 > l && 2 <= e[h + 2] && 2 <= e[h + 1] && (g(h), g(h), k.call(this, h), d(h), d(h));
                            break;
                        case 2:
                            b(h);
                            k.call(this, h + 1);
                            a(h);
                            7 > l && e[h + 2] && e[h + 1] && (g(h), k.call(this,
                                h), d(h));
                            break;
                        case 1:
                            6 > l && 1 == e[h + 1] && e[h + 2] && 4 != e[h + 3] ? (g(h), k.call(this, h + 2), d(h)) : (p = h, --e[p], v |= 1 << p, k.call(this, h + 1), p = h, ++e[p], v &= ~(1 << p), 7 > l && e[h + 2] && (e[h + 1] && (g(h), k.call(this, h + 1), d(h)), r(h), k.call(this, h + 1), m(h)), 8 > l && e[h + 1] && (c(h), k.call(this, h + 1), f(h)))
                    }
                }
            }
        }
    }();

    function F(b, a) {
        E.l(b, 34);
        var g = E.j();
        if (14 < g) return -2;
        !a && 13 <= g && (E.g = E.o(g));
        E.m(g);
        E.s(Math.floor((14 - g) / 3));
        return E.g
    }

    function G(b, a) {
        E.l(b, a);
        var g = E.j();
        if (!(14 < g)) {
            var d = [E.g, E.g];
            13 <= g && (d[0] = E.o(g));
            E.m(g);
            E.s(Math.floor((14 - g) / 3));
            d[1] = E.g;
            d[1] < d[0] && (d[0] = d[1]);
            return d
        }
    };

    function H(b) {
        var a = b >> 2;
        return (27 > a && 16 == b % 36 ? "0" : a % 9 + 1) + "mpsz".substr(a / 9, 1)
    }

    function J(b) {
        return b.replace(/\d(m|p|s|z)(\d\1)*/g, "$&:").replace(/(m|p|s|z)([^:])/g, "$2").replace(/:/g, "")
    }

    function aa(b) {
        b = b.replace(/(\d)m/g, "0$1").replace(/(\d)p/g, "1$1").replace(/(\d)s/g, "2$1").replace(/(\d)z/g, "3$1");
        var a, g = Array(136);
        for (a = 0; a < b.length; a += 2) {
            var d = b.substr(a, 2), c;
            d % 10 ? (c = 4 * (9 * Math.floor(d / 10) + (d % 10 - 1)), c = g[c + 3] ? g[c + 2] ? g[c + 1] ? c : c + 1 : c + 2 : c + 3) : c = 4 * (9 * d / 10 + 4) + 0;
            g[c] && document.write("err n=" + d + " k=" + c + "<br>");
            g[c] = 1
        }
        return g
    };

    function ba(b) {
        var a = parseInt(b.substr(0, 1));
        return (a ? a - 1 : 4) + 9 * "mpsz".indexOf(b.substr(1, 1))
    }

    function K(b) {
        var a, g = [];
        for (a = 0; 34 > a; ++a) 4 <= b[a] || (b[a]++, u(b, 34) && g.push(a), b[a]--);
        return g
    }

    function ca(b) {
        var a,
            g = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];
        for (a = 0; 136 > a; ++a) b[a] && ++g[a >> 2];
        return g
    }

    function da(b, a, g, d, c, f) {
        return '<a href="?' + g + "=" + d + '" class=D onmouseover="daFocus(this,' + c + "," + f + ');" onmouseout="daUnfocus();" >' + L(b, a) + "</a>"
    }

    function L(b, a) {
        return '<img src="/2/t/' + b + '.gif" border=0 ' + (a ? a : "") + " />"
    }

    function M(b) {
        return -1 == b ? "\u548c\u4e86" : 0 == b ? "\u8074\u724c" : b + "\u5411\u8074"
    }

    function N(b, a) {
        return a && b[0] != b[1] ? "\u6a19\u6e96\u5f62" + M(b[0]) + " / \u4e00\u822c\u5f62" + M(b[1]) : M(b[0])
    }

    function ea(b) {
        function a(a) {
            a &= 127;
            return 21 > a ? (a = 9 * Math.floor(a / 7) + a % 7, L(H(4 * a + 1)) + L(H(4 * a + 5)) + L(H(4 * a + 9))) : 55 > a ? (a = L(H(4 * (a - 21) + 1)), a + a + a) : 89 > a ? (a = L(H(4 * (a - 55) + 1)), a + a + a + a) : ""
        }

        function g(a) {
            a = L(H(4 * a + 1));
            return a + a
        }

        var d = new w;
        if (C(d, b)) {
            var c = "";
            for (b = 0; 4 > b; ++b) if (d.c[b].a) {
                if (0 == b && 4294967295 == d.c[0].a) {
                    var c = c + "\u56fd\u58eb\u5f62\u548c\u4e86 ", c = c + (g(d.c[b].b) + " "), f,
                        r = [0, 8, 9, 17, 18, 26, 27, 28, 29, 30, 31, 32, 33];
                    for (f = 0; 13 > f; ++f) d.c[b].b != r[f] && (c += L(H(4 * r[f] + 1)))
                } else 3 == b && 4294967295 == d.c[3].a ?
                    (c += "\u4e03\u5bfe\u5f62\u548c\u4e86 ", c += g(d.h[0]) + " " + g(d.h[1]) + " " + g(d.h[2]) + " " + g(d.h[3]) + " " + g(d.h[4]) + " " + g(d.h[5]) + " " + g(d.h[6])) : (f = [(d.c[b].a >> 0 & 255) - 1, (d.c[b].a >> 8 & 255) - 1, (d.c[b].a >> 16 & 255) - 1, (d.c[b].a >> 24 & 255) - 1], c += "\u4e00\u822c\u5f62\u548c\u4e86 ", c += g(d.c[b].b) + " " + a(f[3]) + " " + a(f[2]) + " " + a(f[1]) + " " + a(f[0]));
                c += "<br>"
            }
            return c
        }
    }

    function fa() {
        function b(a, b) {
            var c, d = 0;
            for (c = 0; c < a.length; ++c) d += 4 - b[a[c]];
            return d
        }

        var a = ga, g = O, d;
        d = "<hr size=1 color=#CCCCCC >";
        switch (a.substr(0, 1)) {
            case "q":
                d += '\u6a19\u6e96\u5f62(\u4e03\u5bfe\u56fd\u58eb\u3092\u542b\u3080)\u306e\u8a08\u7b97\u7d50\u679c / <a href="?p' + a.substr(1) + "=" + g + '">\u4e00\u822c\u5f62</a><br>';
                break;
            case "p":
                d += '\u4e00\u822c\u5f62(\u4e03\u5bfe\u56fd\u58eb\u3092\u542b\u307e\u306a\u3044)\u306e\u8a08\u7b97\u7d50\u679c / <a href="?q' + a.substr(1) + "=" + g + '">\u6a19\u6e96\u5f62</a><br>'
        }
        for (var c =
            "d" == a.substr(1, 1), a = a.substr(0, 1), g = g.replace(/(\d)(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{0,8})(\d{8})(m|p|s|z)/g, "$1$9$2$9$3$9$4$9$5$9$6$9$7$9$8$9").replace(/(\d?)(\d?)(\d?)(\d?)(\d?)(\d?)(\d)(\d)(m|p|s|z)/g, "$1$9$2$9$3$9$4$9$5$9$6$9$7$9$8$9").replace(/(m|p|s|z)(m|p|s|z)+/g, "$1").replace(/^[^\d]/, ""), g = g.substr(0, 28), f = aa(g), r = -1; r = Math.floor(136 * Math.random()), f[r];) ;
        var m = Math.floor(g.length / 2) % 3;
        2 == m || c || (f[r] = 1, g += H(r));
        var f = ca(f), n = "", e = G(f, 34), n = n + N(e, 28 == g.length), n = n + ("(" + Math.floor(g.length /
            2) + "\u679a)");
        -1 == e[0] && (n += ' / <a href="?" >\u65b0\u3057\u3044\u624b\u724c\u3092\u4f5c\u6210</a>');
        var n = n + "<br/>", q = "q" == a ? e[0] : e[1], k, p, l = Array(35);
        if (0 == q && 1 == m && c) k = 34, l[k] = K(f), l[k].length && (l[k] = {
            i: k,
            n: b(l[k], f),
            c: l[k]
        }); else if (0 >= q) for (k = 0; 34 > k; ++k) f[k] && (f[k]--, l[k] = K(f), f[k]++, l[k].length && (l[k] = {
            i: k,
            n: b(l[k], f),
            c: l[k]
        })); else if (2 == m || 1 == m && !c) for (k = 0; 34 > k; ++k) {
            if (f[k]) {
                f[k]--;
                l[k] = [];
                for (p = 0; 34 > p; ++p) k == p || 4 <= f[p] || (f[p]++, F(f, "p" == a) == q - 1 && l[k].push(p), f[p]--);
                f[k]++;
                l[k].length && (l[k] =
                    {i: k, n: b(l[k], f), c: l[k]})
            }
        } else {
            k = 34;
            l[k] = [];
            for (p = 0; 34 > p; ++p) 4 <= f[p] || (f[p]++, F(f, "p" == a) == q - 1 && l[k].push(p), f[p]--);
            l[k].length && (l[k] = {i: k, n: b(l[k], f), c: l[k]})
        }
        var t = [];
        for (k = 0; k < g.length; k += 2) {
            p = g.substr(k, 2);
            var v = ba(p),
                h = J(g.replace(p, "").replace(/(\d)(m|p|s|z)/g, "$2$1$1,").replace(/00/g, "50").split(",").sort().join("").replace(/(m|p|s|z)\d(\d)/g, "$2$1")),
                R = q + 1, I = l[v];
            I && I.n && (R = -1 == q ? 0 : q, void 0 == I.q && t.push(I), I.q = h);
            2 == m && (h += H(r));
            n += (2 == m || 2 != m && !c ? da : L)(p, 2 == k % 3 && k == g.length - 2 ? " hspace=3 " :
                "", a, h, v, R)
        }
        l[34] && l[34].n && (l[34].q = J(g), t.push(l[34]), n += '<br><br><a href="?' + a + "=" + l[34].q + '">\u6b21\u306e\u30c4\u30e2\u3092\u30e9\u30f3\u30c0\u30e0\u306b\u8ffd\u52a0</a>');
        t.sort(function (a, b) {
            return b.n - a.n
        });
        g = "" + (document.f.q.value + "\n");
        d += "<table cellpadding=2 cellspacing=0 >";
        q = 0 >= q ? "\u5f85\u3061" : "\u6478";
        for (k = 0; k < t.length; ++k) {
            v = t[k].i;
            d += "<tr id=mda" + v + " ><td>";
            34 > v && (d += "\u6253</td><td>" + ('<img src="/2/a/' + H(4 * v + 1) + '.gif" class=D />') + "</td><td>", g += "\u6253" + H(4 * v + 1) + " ");
            d += q + "[</td><td>";
            g += q + "[";
            l = t[k].c;
            c = t[k].q;
            for (p = 0; p < l.length; ++p) r = H(4 * l[p] + 1), d += '<a href="?' + a + "=" + (c + r) + '" class=D onmouseover="daFocus(this,' + v + ');" onmouseout="daUnfocus();"><img src="/2/a/' + r + '.gif" border=0 /></a>', g += H(4 * l[p] + 1);
            d += "</td><td>" + t[k].n + "\u679a</td><td>]</td></tr>";
            g += " " + t[k].n + "\u679a]\n"
        }
        d = d + "</table><br><hr><br>" + ('<textarea rows=10 style="width:100%;font-size:75%;">' + g + "</textarea>");
        -1 == e[0] && (d = d + "<hr size=1 color=#CCCCCC >" + ea(f));
        document.getElementById("tehai").innerHTML = n;
        document.getElementById("tips").innerHTML =
            "";
        document.getElementById("m2").innerHTML = d + "<br>"
    }

    document.write('<form name=f style="margin:0px;" ><a href="?" >\u65b0\u898f</a> | \u624b\u724c <input type=text name=q size=20 ><input type=submit value=" OK "><hr size=1 color=#CCCCCC ></form><div id=tehai align=center ></div><div id=tips align=center style="height:18px"></div><div id=m2 align=center ></div><hr size=1 color=#CCCCCC >- m=\u842c\u5b50, p=\u7b52\u5b50, s=\u7d22\u5b50, z=\u5b57\u724c, 0=\u8d64<br>- \u4e00\u822c\u5f62=\uff14\u9762\u5b50\uff11\u96c0\u982d / \u6a19\u6e96\u5f62=\u4e00\u822c\u5f62\uff0b\u4e03\u5bfe\u5f62\uff0b\u56fd\u58eb\u5f62<br>- \u30c4\u30e2\u306f\u305d\u306e\u6642\u70b9\u3067\u4f7f\u7528\u3057\u3066\u3044\u306a\u3044\u724c\u3092\u30e9\u30f3\u30c0\u30e0\u306b\u9078\u629e\u3057\u307e\u3059<br>- \u6709\u52b9\u724c\u3092\u30af\u30ea\u30c3\u30af\u3059\u308b\u3068\u6253\u724c\u5f8c\u306b\u305d\u306e\u724c\u3092\u30c4\u30e2\u724c\u3068\u3057\u3066\u8868\u793a\u3057\u307e\u3059<br>- (n*3+2)\u679a\u3067\u958b\u59cb\uff1a(n*3+2)\u679a\u76ee\u3092\u30c4\u30e2\u724c\u3068\u3057\u3066\u8868\u793a<br>- (n*3+1)\u679a\u3067\u958b\u59cb\uff1a\u30c4\u30e2\u306f\u30da\u30fc\u30b8\u306e\u30ed\u30fc\u30c9\u6642\u306b\u6bce\u56de\u5909\u5316<br>- \u548c\u4e86\u5f79\u306e\u5224\u5b9a\u306f\u3042\u308a\u307e\u305b\u3093<br>- \u6697\u69d3\u306f\u3067\u304d\u307e\u305b\u3093<br>');
    var O = window.location.search.substr(1), O = O.replace(/^([^=]+)=(.+)/, "$2"), ga = RegExp.$1;
    document.f.q.value = O;
    if (O.length) null == O.match(/^(\d+m|\d+p|\d+s|[1234567]+z)*$/) ? document.getElementById("tehai").innerHTML = "<font color=#FF0000 >INVALID QUERY</font>" : O.length && fa(); else {
        var P, Q = "", S = Math.floor(3 * Math.random()), Q = Q + "<table cellpadding=0 cellspacing=0 ><tr><td>";
        for (P = 9; 0 <= P; --P) {
            for (var T, ha = 3 < P ? 136 : 36, U = void 0, V = [], U = 0; U < ha; ++U) V.push(U);
            for (U = 0; U < V.length - 1; ++U) {
                var W = U + Math.floor(Math.random() * (V.length - U)), ia = V[U];
                V[U] = V[W];
                V[W] = ia
            }
            T = V;
            var X = T.splice(0, 9 == P ? 4 : 8 == P ? 7 : 7 == P ? 10 : 13).sort(function (b,
                                                                                           a) {
                return b - a
            }), Y = T.splice(0, 1)[0], Z, O = "";
            for (Z = 0; Z < X.length; ++Z) O += H(X[Z] + (3 < P ? 0 : 36 * S));
            -1 != Y && (O += H(Y + (3 < P ? 0 : 36 * S)), X.push(Y));
            var ja = G(X, X.length), Q = Q + ('\u25a0<a href="?q=' + J(O) + '" class=D >'),
                Q = Q + N(ja, 14 == X.length), Q = Q + "<br>";
            for (Z = 0; Z < O.length; Z += 2) Q += L(O.substr(Z, 2), 2 == Z % 6 && Z == O.length - 2 ? " hspace=3 " : "");
            Q += "</a>";
            Q += "<br><br>";
            4 == P && (document.f.q.value = J(O))
        }
        Q += "</td></tr></table>";
        document.getElementById("tehai").innerHTML = Q
    }
    ;
})();
//