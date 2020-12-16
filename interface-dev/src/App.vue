<template>
  <div id="app">
    <TopHeader></TopHeader>
    
  </div>
</template>

<script>
  import TopHeader from './components/TopHeader'
export default {
  name: 'App',
  components:{TopHeader}
}
</script>

<style lang="less">
  @import (css) url('https://fonts.googleapis.com/css2?family=Roboto&display=swap');

* {
    margin:0;
    padding:0;
    font-family: 'Roboto',sans-serif;
  }
  html, body {
    width:100%;
    height:100%;
  }
  body {
    background: #ecf0f1;
  }

  .hidden{
    display:none !important;
  }
  .clearall{
    clear:both;
  }

  /* My Flex Lib */
  .flex-container {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    justify-content: flex-start;

    &.flex-wrap{ flex-wrap:wrap; }
    &.align-items-top{  align-items: flex-start;}
    &.align-items-stretch{align-items: stretch;}
    &.flex-justify-content-left{justify-content: flex-start;}
    &.flex-space-between{justify-content: space-between;}
    &.flex-space-around{justify-content: space-around;}
    &.flex-to-bottom{
      flex-direction: column;
    }

    & > * {
      flex-grow:1;
      flex-shrink:1;
    }
  }


</style>
