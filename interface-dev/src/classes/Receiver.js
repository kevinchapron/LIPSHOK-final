export default class Receiver{

  constructor(fromWS) {
    for(let key in fromWS){
      this[key] = fromWS[key]
    }
  }
}
