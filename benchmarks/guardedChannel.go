//START guardedChannel code

/*
  guide:
  1) replace TYPENAME with the name of the type eg. intarray
  2) replace TYPE with the actual type eg. []int
*/

type guarded<TYPENAME>Channel struct{
  dataChannel chan <TYPE>
  boolChannel chan bool
}

func newGuarded<TYPENAME>Channel(capacity int) guarded<TYPENAME>Channel {
  return guarded<TYPENAME>Channel{make(chan <TYPE>,capacity), make(chan bool,capacity)}
}

func (channel *guarded<TYPENAME>Channel) send(value <TYPE>) {
  channel.dataChannel <- value
  channel.boolChannel <- true
}

func (channel *guarded<TYPENAME>Channel) recv(value *<TYPE>) {
  *value = <- channel.dataChannel
  _ = <- channel.boolChannel
}

func (channel *guarded<TYPENAME>Channel) trysend(cond bool, value <TYPE>) {
  channel.dataChannel <- value
  channel.boolChannel <- cond
}

func (channel *guarded<TYPENAME>Channel) tryrecv(value *<TYPE>) bool {
  cond := <- channel.boolChannel
  if cond {
    *value = <- channel.dataChannel
  }
  return cond
}

//END guardedChannel code
