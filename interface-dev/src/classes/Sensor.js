export default class Sensor{

  constructor(fromWS, listReceivers) {
    this.name = fromWS.Name;
    this.protocol = fromWS.Protocol;
    this.receiverID = -1;
    this.lastSeen = null;
    this.value = null;

    for(let i=0;i<listReceivers.length;i++){
      if(listReceivers[i].protocol == this.protocol){
        this.receiverID = i
        break
      }
    }
  }
}
