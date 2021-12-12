<template>
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
</template>

<script>
import Discards from './components/Discards.vue'

export default {
  name: 'App',
  data() {
    return {
      discards: [[], [], [], []]
    }
  },
  components: {
    Discards,
  },
  created: function() {
    this.ws = new WebSocket(`ws://${location.host}/ws`);

    this.ws.onmessage = (event) => {
      console.log(event);
      console.log(event.data);
      //this.$set('discards', JSON.parse(event.data).discards)
      this.discards = JSON.parse(event.data).discards
      //vm.$set(this.discards, JSON.parse(event.data).discards)
    }

    this.ws.onopen = (event) => {
      console.log(event)
      console.log("Successfully connected")
    }
  }
}
//const DISCARDS = [
//  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": false, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": false, "isRiichi": true},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//  {"tile": "1m", "isTsumogiri": true, "isRiichi": false},
//]
</script>
