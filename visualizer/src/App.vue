<template>
  <div :style="`
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
  `">
    <el-card
      shadow="never"
      :style="`
        background-color: #25498D;
        width:645px;
        border: 0;
      `"
    >
      <el-row justify="center">
        <Discards :rotation="180" :discards="discards[2]"/>
      </el-row>
      <el-row justify="center">
        <Discards :rotation="90" :discards="discards[3]"/>
        <div style="width: 200px;"/>
        <Discards :rotation="270" :discards="discards[1]"/>
      </el-row>
      <el-row justify="center">
        <Discards :discards="discards[0]"/>
      </el-row>
    </el-card>
  </div>
</template>

<script>
import Discards from './components/Discards.vue'

export default {
  name: 'App',
  data() {
    return {
      //discards: [DISCARDS, DISCARDS, DISCARDS, DISCARDS]
      discards: [[], [], [], []]
    }
  },
  components: {
    Discards,
  },
  created: function() {
    this.ws = new WebSocket(`ws://${location.host}/ws`);

    this.ws.onmessage = (event) => {
      this.discards = JSON.parse(event.data).discards
    }

    // eslint-disable-next-line no-unused-vars
    this.ws.onopen = (event) => {
      console.log("Successfully connected")
    }
  }
}
// eslint-disable-next-line no-unused-vars
const DISCARDS = [
  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": false, "isRiichi": true},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
]
</script>
<style>
body {
  background-color: #A2A2A2;
}
</style>
