<template>
  <div id="content" class="flex-container">
    <div class="receivers">
      <h2>Connectors</h2>
      <div class="flex-container">
        <div class="receiver" v-for="receiver in receivers">
          <h2>{{ receiver.name }}</h2>
          <hr/>
          <ul class="fa-ul">
            <li title="Addr"><span class="fa-li"><i class="fas fa-map-marker-alt"></i></span>{{ receiver.addr }}</li>
            <li title="Protocol"><span class="fa-li"><i class="fas fa-share-alt-square"></i></span>{{ receiver.protocol }}
            </li>
            <li title="Last Seen"><span class="fa-li"><i class="far fa-eye"></i></span>{{ formatTime(receiver.lastSeen) }}
            </li>
          </ul>
        </div>
      </div>
    </div>
    <div class="vertical-line"></div>
    <div class="sensors">
      <h2>Sensors</h2>
      <div class="flex-container">
        <div class="sensor" v-for="sensor in sensors">
          <h2>{{ sensor.name }}</h2>
          <hr/>
          <ul class="fa-ul">
            <li title="Receiver"><span class="fa-li"><i class="fas fa-share-alt-square"></i></span><span v-if="sensor.receiverID !== -1">{{ receivers[sensor.receiverID].name }}</span><span v-else>Aucun</span></li>
            <li title="Value"><span class="fa-li"><i class="fas fa-database"></i></span>{{ sensor.value }}</li>
            <li title="Last Seen"><span class="fa-li"><i class="far fa-eye"></i></span>{{ formatTime(sensor.lastSeen) }}</li>
          </ul>
<!--          <ul class="fa-ul">-->
<!--            <li title="Addr"><span class="fa-li"><i class="fas fa-map-marker-alt"></i></span>{{ sensor.addr }}</li>-->
<!--            <li title="Protocol"><span class="fa-li"><i class="fas fa-share-alt-square"></i></span>{{ sensor.protocol }}-->
<!--            </li>-->
<!--            <li title="Last Seen"><span class="fa-li"><i class="far fa-eye"></i></span>{{ formatTime(sensor.lastSeen) }}-->
<!--            </li>-->
<!--          </ul>-->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Sensor from "../classes/Sensor";
import Receiver from "../classes/Receiver";

export default {
  name: "Content",
  data() {
    return {
      receivers: [],
      sensors: [],
    }
  },
  created() {
    this.ws = new WebSocket("ws://127.0.0.1:5003/output");

    this.ws.onmessage = (event) => {
      let data = JSON.parse(event.data);
      data["data"] = JSON.parse(atob(data["data"]))
      switch (data["type"]) {
        case 0:
          // normal data
          for(let i=0;i<this.sensors.length;i++){
            let sensor = this.sensors[i];
            if(sensor.name === data.from.Name){
              // traiter la data du sensor ici
              this.sensors[i].lastSeen = data.datetime;
              this.sensors[i].value = data.data.value;

              for(let j=0;j<this.receivers.length;j++){
                if(this.receivers[j].protocol === sensor.protocol){
                  this.receivers[j].lastSeen = data.datetime;
                  break
                }
              }

              console.log(this.sensors[i])
              break;
            }
          }
          break
        case 1:
          // auth
          this.buildReceivers(data["data"]["connectors"])
          this.buildSensors(data["data"]["sensors"])
      }

    };

    this.ws.onopen = (event) => {
      this.ws.send("status")
    };
  },
  methods: {
    formatTime: function (time) {
      if(time == null){
        return "No time";
      }
      return time.substr(0,10)+" "+time.substr(11,8)
    },
    buildSensors: function(arr){
      for(let i=0;i<arr.length;i++){
        this.sensors.push(new Sensor(arr[i], this.receivers))
      }
    },
    buildReceivers: function(arr){
      for(let i=0;i<arr.length;i++){
        this.receivers.push(new Receiver(arr[i]))
      }
    }
  },
}
</script>

<style scoped lang="less">
#content {
  min-height: calc(100% - 57px);


  & > * {
    flex: 0;

    margin: 20px;
    margin-right: 0;

    &:last-child {
      margin-right: 20px;
    }

    &.vertical-line {
      flex-basis:5px;
      background: rgba(0, 0, 0, 0.5);
    }
  }

  .receivers, .sensors {
    align-items: flex-start;
    flex:0;

    &.receivers{
      flex-basis:30vw;
    }
    &.sensors{
      flex-basis:69vw;
    }
    & > h2{
      margin-bottom:10px;
      padding:20px 0;
      border-bottom:2px solid gray;
      text-align:center;
    }
    .flex-container{
      flex-wrap: wrap;
      flex-basis:30vw;
    }
    .receiver, .sensor {
      min-width: 200px;
      &.receiver{ background-color: #34495e; }
      &.sensor{   background-color: fadeout(#34495e,20%); }
      border: 1px solid #95a5a6;
      margin: 20px;
      padding: 10px;
      flex-shrink:0;
      flex-grow:0.25;

      h2 {
        color: #ecf0f1;
        text-align: center;
        margin: 7px 0;
      }

      hr {
        margin: 7px 0;
        border: 0;
        border-top: 1px solid #bdc3c7;
      }

      ul {
        color: #ecf0f1;

        li {
          margin-bottom: 10px;

          &:first-child {
            margin-top: 20px;
          }
        }
      }
    }
  }
}
</style>
