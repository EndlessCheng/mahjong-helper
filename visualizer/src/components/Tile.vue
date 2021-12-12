<template>
  <div 
    class="radius"
    :style="`
        border: 1px solid var(--el-border-color-base);
        border-radius: var(--el-border-radius-base);
        width: ${width};
        height: ${height};
        margin: ${margin};
        background-color: ${background};
      `">
    <img
      :src="img_src"
      :style="`
        width: ${width};
        height: ${height};
      `"
    >
  </div>
</template>

<script>
export default {
  name: 'Tile',
  props: {
    tile: {
      type: String,
      default: "1m",
    },
    rotation: {
      type: Number,
      default: 0,
    },
    dim: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    img_src: function() {
      if (this.rotation == 0) {
        return require(`../assets/${IMGFILE[this.tile]}.svg`)
      } else {
        return require(
          `../assets/${IMGFILE[this.tile]}-90.svg`)
      }
    },
    width: function() {
      return this.rotation % 180 == 0 ? "30px" : "40px"
    },
    height: function() {
      return this.rotation % 180 == 0 ? "40px" : "30px"
    },
    margin: function() {
      return this.rotation % 180 == 0 ? "0" : "5px 0"
    },
    background: function() {
      return this.dim ? "#909399" : ""
    }
  }
}
const IMGFILE = {
  "1m": "Man1",
  "2m": "Man2",
  "3m": "Man3",
  "4m": "Man4",
  "5m": "Man5",
  "5m-dora": "Man5-Dora",
  "6m": "Man6",
  "7m": "Man7",
  "8m": "Man8",
  "9m": "Man9",
  "1p": "Pin1",
  "2p": "Pin2",
  "3p": "Pin3",
  "4p": "Pin4",
  "5p": "Pin5",
  "5p-dora": "Pin5-Dora",
  "6p": "Pin6",
  "7p": "Pin7",
  "8p": "Pin8",
  "9p": "Pin9",
  "1s": "Sou1",
  "2s": "Sou2",
  "3s": "Sou3",
  "4s": "Sou4",
  "5s": "Sou5",
  "5s-dora": "Sou5-Dora",
  "6s": "Sou6",
  "7s": "Sou7",
  "8s": "Sou8",
  "9s": "Sou9",
  "1z": "Ton",
  "2z": "Nan",
  "3z": "Shaa",
  "4z": "Pei",
  "5z": "Haku",
  "6z": "Hatsu",
  "7z": "Chun",
}
</script>

<style scoped>
</style>
