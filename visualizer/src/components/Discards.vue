<template>
  <div :style="`
    transform: rotate(${rotation}deg);
    height: 202px;
    width: 202px;
  `">
    <el-row
      v-for="(discards_in_row,index) in discards_aligned"
      :key="index"
    >
        <Tile
          v-for="(discard,index) in discards_in_row"
          :tile="discard.tile + (discard.isRedFive ? '-dora' : '')"
          :rotation="discard.isRiichi ? 90 : 0"
          :dim="discard.isTsumogiri"
          :key="index"
        />
    </el-row>
  </div>
</template>

<script>
import Tile from '../components/Tile.vue'

export default {
  name: 'Discards',
  components: {
    Tile,
  },
  props: {
    rotation: {
      type: Number,
      default: 0,
    },
    discards: {
      type: Array,
      default: () => {
        return []
      }
    }
  },
  computed: {
    discards_aligned: function() {
      var ret = []
      for (var i = 0; i < this.discards.length; i+=6) {
        ret.push(this.discards.slice(i, i+6))
      }
      return ret
    },
  }
}
</script>

<style scoped>
.el-card {
  width: 250px;
}
</style>
