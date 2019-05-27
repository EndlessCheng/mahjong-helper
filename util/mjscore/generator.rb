#! /usr/bin/ruby

# http://hp.vector.co.jp/authors/VA046927/mjscore/mjalgorism.html
# http://hp.vector.co.jp/authors/VA046927/mjscore/ptn.rb

# 順列生成
class Array
  def perms
    return [[]] if empty?
    uniq.inject([]) do |rs, h|
      tmp = self.dup
      tmp.delete_at(index(h))
      rs + tmp.perms.map {|t| [h] + t}
    end
  end

  def total
    self.inject(0) {|t, a| t += a}
  end
end

# a : [[1,1,1], [1,1,1], [1,1,1], [1,1,1], [2]]
def ptn(a)
  if a.size == 1 then
    return [a]
  end
  ret = Array.new
  # 重ならないパターン
  ret += a.perms
  # 重なるパターン
  h1 = Hash.new
  for i in 0..a.size - 1
    for j in i + 1..a.size - 1
      key = [a[i], 0, a[j]].to_s
      if !h1.key?(key) then
        h1.store(key, nil)
        h2 = Hash.new
        # a[i]とa[j]を範囲をずらしながら重ねる
        for k in 0..a[i].size + a[j].size
          t = [0] * a[j].size + a[i] + [0] * a[j].size
          for m in 0..a[j].size - 1
            t[k + m] += a[j][m]
          end
          # 余分な0を取り除く
          t.delete(0)
          # 4より大きい値がないかチェック
          next if t.any? {|v| v > 4}
          # 9より長くないかチェック
          next if t.size > 9
          # 重複チェック
          if !h2.key?(t.to_s) then
            h2.store(t.to_s, nil)
            # 残り
            t2 = a.dup
            t2.delete_at(i)
            t2.delete_at(j - 1)
            # 再帰呼び出し
            ret += ptn([t] + t2)
          end
        end
      end
    end
  end
  return ret
end

# キー値を計算
def calc_key(a)
  ret = 0
  len = -1
  for b in a
    for i in b
      len += 1
      case i
      when 2 then
        ret |= 0b11 << len
        len += 2
      when 3 then
        ret |= 0b1111 << len
        len += 4
      when 4 then
        ret |= 0b111111 << len
        len += 6
      end
    end
    ret |= 0b1 << len
    len += 1
  end
  return ret
end

# a : [[1,1,1], [1,1,1], [1,1,1], [1,1,1], [2]]
# ret
# 下位
#   3bit  0: 刻子の数(0～4)
#   3bit  3: 順子の数(0～4)
#   4bit  6: 頭の位置(1～13)
#   4bit 10: 面子の位置１(0～13)
#   4bit 14: 面子の位置２(0～13)
#   4bit 18: 面子の位置３(0～13)
#   4bit 22: 面子の位置４(0～13)
#   1bit 26: 七対子フラグ
#   1bit 27: 九蓮宝燈フラグ
#   1bit 28: 一気通貫フラグ
#   1bit 29: 二盃口フラグ
#   1bit 30: 一盃口フラグ
def find_hai_pos(a)
  ret_array = Array.new
  p_atama = 0
  for i in 0..a.size - 1
    for j in 0..a[i].size - 1
      # 頭を探す
      if a[i][j] >= 2 then
        # 刻子、順子の優先順位入れ替え
        for kotsu_shuntus in 0..1
          t = Marshal.load(Marshal.dump(a))
          t[i][j] -= 2

          p = 0
          p_kotsu = Array.new
          p_shuntsu = Array.new
          for k in 0..t.size - 1
            for m in 0..t[k].size - 1
              if kotsu_shuntus == 0 then
                # 刻子を先に取り出す
                # 刻子
                if t[k][m] >= 3 then
                  t[k][m] -= 3
                  p_kotsu.push(p)
                end
                # 順子
                while t[k].size - m >= 3 &&
                    t[k][m] >= 1 &&
                    t[k][m + 1] >= 1 &&
                    t[k][m + 2] >= 1 do
                  t[k][m] -= 1
                  t[k][m + 1] -= 1
                  t[k][m + 2] -= 1
                  p_shuntsu.push(p)
                end
              else
                # 順子を先に取り出す
                # 順子
                while t[k].size - m >= 3 &&
                    t[k][m] >= 1 &&
                    t[k][m + 1] >= 1 &&
                    t[k][m + 2] >= 1 do
                  t[k][m] -= 1
                  t[k][m + 1] -= 1
                  t[k][m + 2] -= 1
                  p_shuntsu.push(p)
                end
                # 刻子
                if t[k][m] >= 3 then
                  t[k][m] -= 3
                  p_kotsu.push(p)
                end
              end
              p += 1
            end
          end

          # 上がりの形か？
          if t.flatten.all? {|x| x == 0} then
            # 値を求める
            ret = p_kotsu.size + (p_shuntsu.size << 3) + (p_atama << 6)
            len = 10
            for x in p_kotsu
              ret |= x << len
              len += 4
            end
            for x in p_shuntsu
              ret |= x << len
              len += 4
            end
            if a.size == 1 then
              # 九蓮宝燈フラグ
              if a == [[4, 1, 1, 1, 1, 1, 1, 1, 3]] ||
                  a == [[3, 2, 1, 1, 1, 1, 1, 1, 3]] ||
                  a == [[3, 1, 2, 1, 1, 1, 1, 1, 3]] ||
                  a == [[3, 1, 1, 2, 1, 1, 1, 1, 3]] ||
                  a == [[3, 1, 1, 1, 2, 1, 1, 1, 3]] ||
                  a == [[3, 1, 1, 1, 1, 2, 1, 1, 3]] ||
                  a == [[3, 1, 1, 1, 1, 1, 2, 1, 3]] ||
                  a == [[3, 1, 1, 1, 1, 1, 1, 2, 3]] ||
                  a == [[3, 1, 1, 1, 1, 1, 1, 1, 4]] then
                ret |= 1 << 27
              end
            end
            # 一気通貫
            if a.size <= 3 && p_shuntsu.size >= 3 then
              p_ikki = 0
              for b in a
                if b.size == 9 then
                  b_ikki1 = false
                  b_ikki2 = false
                  b_ikki3 = false
                  for x_ikki in p_shuntsu
                    b_ikki1 |= (x_ikki == p_ikki)
                    b_ikki2 |= (x_ikki == p_ikki + 3)
                    b_ikki3 |= (x_ikki == p_ikki + 6)
                  end
                  if b_ikki1 && b_ikki2 && b_ikki3 then
                    ret |= 1 << 28
                  end
                end
                p_ikki += b.size
              end
            end
            # 二盃口
            if p_shuntsu.size == 4 &&
                p_shuntsu[0] == p_shuntsu[1] &&
                p_shuntsu[2] == p_shuntsu[3] then
              ret |= 1 << 29
            elsif p_shuntsu.size >= 2 && p_kotsu.size + p_shuntsu.size == 4 then
              # 一盃口
              if p_shuntsu.size - p_shuntsu.uniq.size >= 1 then
                ret |= 1 << 30
              end
            end
            ret_array.push(ret)
          end
        end
      end
      p_atama += 1
    end
  end
  if ret_array.size > 0 then
    ret_array.uniq!
    return ret_array
  end
  t = a.flatten
  # 七対子判定
  if t.total == 14 && t.all? {|x| x == 2} then
    return [1 << 26]
  end
end

chitoi = ptn([[2], [2], [2], [2], [2], [2], [2]])
chitoi.delete_if {|x|
  t = x.flatten
  t.any? {|y| y != 2}
}

(ptn([[1, 1, 1], [1, 1, 1], [1, 1, 1], [1, 1, 1], [2]]) +
    ptn([[1, 1, 1], [1, 1, 1], [1, 1, 1], [3], [2]]) +
    ptn([[1, 1, 1], [1, 1, 1], [3], [3], [2]]) +
    ptn([[1, 1, 1], [3], [3], [3], [2]]) +
    ptn([[3], [3], [3], [3], [2]]) +
    chitoi).uniq.each do |x|
  printf("%d", calc_key(x))
  find_hai_pos(x).map { |i| printf(" %d", i) }
  printf("\n")
end

(ptn([[1, 1, 1], [1, 1, 1], [1, 1, 1], [2]]) +
    ptn([[1, 1, 1], [1, 1, 1], [3], [2]]) +
    ptn([[1, 1, 1], [3], [3], [2]]) +
    ptn([[3], [3], [3], [2]])).uniq.each do |x|
  printf("%d", calc_key(x))
  find_hai_pos(x).map { |i| printf(" %d", i) }
  printf("\n")
end

(ptn([[1, 1, 1], [1, 1, 1], [2]]) +
    ptn([[1, 1, 1], [3], [2]]) +
    ptn([[3], [3], [2]])).uniq.each do |x|
  printf("%d", calc_key(x))
  find_hai_pos(x).map { |i| printf(" %d", i) }
  printf("\n")
end

(ptn([[1, 1, 1], [2]]) +
    ptn([[3], [2]])).uniq.each do |x|
  printf("%d", calc_key(x))
  find_hai_pos(x).map { |i| printf(" %d", i) }
  printf("\n")
end

(ptn([[2]])).uniq.each do |x|
  printf("%d", calc_key(x))
  find_hai_pos(x).map { |i| printf(" %d", i) }
  printf("\n")
end
