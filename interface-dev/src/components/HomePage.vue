<template>
    <div id="homepage">
      <TopHeader></TopHeader>
      <Content></Content>
    </div>
</template>

<script>
  import TopHeader from '@/components/TopHeader'
  import Content from '@/components/Content'
	export default {
		name: "HomePage",
    components: {
      TopHeader,
      Content
    }
	}
</script>

<style scoped>
  #homepage{
    width:100%;
    height:100%;
  }
</style>
