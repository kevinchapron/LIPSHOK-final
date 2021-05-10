<template>
    <div id="content" class="flex-container">
      <div class="receivers flex-container">
        <div class="receiver" v-for="receiver in receivers">
          <h2>{{ receiver.name }}</h2>
          <hr />
          <ul class="fa-ul">
            <li title="Addr"><span class="fa-li"><i class="fas fa-map-marker-alt"></i></span>{{ receiver.addr }}</li>
            <li title="Protocol"><span class="fa-li"><i class="fas fa-share-alt-square"></i></span>{{ receiver.protocol }}</li>
            <li title="Last Seen"><span class="fa-li"><i class="far fa-eye"></i></span>{{ formatTime(receiver.lastSeen) }}</li>
          </ul>
        </div>
      </div>
      <div class="vertical-line"></div>
      <div class="sensors">Sensors</div>
    </div>
</template>

<script>
	export default {
		name: "Content",
    data(){
		  return {
		    receivers: {}
      }
    },
    created(){
		  this.ws = new WebSocket("ws://127.0.0.1:5003/output");

      this.ws.onmessage = (event) => {
		    let data = JSON.parse(event.data);
		    let rawData = atob(data["data"]);
        this.receivers = JSON.parse(rawData)
      };

      this.ws.onopen = (event) => {
		    this.ws.send("status")
      };
    },
    methods:{
		  formatTime: function(time){
		    let t = time.split(" ")
        delete t[t.length-1]
        delete t[t.length-2]
        delete t[t.length-3]
        t[1] = t[1].split(".")
        delete t[1][1]
        t[1] = t[1].join(" ")
		    return t.join(" ")
      }
    }
	}
</script>

<style scoped lang="less">
  #content{
    min-height:calc(100% - 57px);


    & > *{
      flex:0;
      &.receivers{
        flex-basis:30%;
      }
      &.sensors{
        flex-basis:70%;
      }
      margin:20px;
      margin-right:0;

      &:last-child{
        margin-right:20px;
      }

      &.vertical-line{
        flex-basis:2px;
        background:rgba(0,0,0,0.5);
      }
    }

    .receivers{
      flex-wrap: nowrap;
      align-items:flex-start;

      .receiver{
        width:200px;
        background-color:#34495e;
        border:1px solid #95a5a6;
        margin:20px;
        padding:10px;
        h2{
          color:#ecf0f1;
          text-align:center;
          margin:7px 0;
        }
        hr {
          margin:7px 0;
          border:0;
          border-top:1px solid #bdc3c7;
        }

        ul{
          color:#ecf0f1;
          li{
            margin-bottom:10px;
            &:first-child{
              margin-top:20px;
            }
          }
        }
      }
    }
  }
</style>
