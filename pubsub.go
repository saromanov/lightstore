package lightstore

import
(
	"time"
	"errors"
	"sync"
	"fmt"
)

type Pubsub struct {
	inner map[string]PubsubObject
	mutex *sync.RWMutex
	queue []PubsubObject
}

type SubscribeData struct{
	Title string
}

type PublishData struct{
	Title string
	Msg string
}

type PubsubObject struct {
	//Catch only once event
	once bool
	//msg
	msg string
}

func PubsubInit()*Pubsub{
	ps := new(Pubsub)
	ps.mutex = &sync.RWMutex{}
	ps.queue = []PubsubObject{}
	ps.inner = map[string]PubsubObject{}
	return ps
}

//Subscribe provides subscibtion to the specific key
//for example on the db title
func (ps *Pubsub) Subscribe(si* SubscribeData){
	ps.inner[si.Title] = PubsubObject{}
	var wg sync.WaitGroup
    wg.Add(1)
	go func(){
		for {
			item, err := ps.receive()
			if err == nil {
				fmt.Println(fmt.Sprintf("Receive from %s %s", si.Title, item.msg))
				if item.once {
					break
					wg.Done()
				}
				item.msg = ""
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	wg.Wait()
}

//Publish new message
func (ps *Pubsub) Publish(po *PublishData) {
	ps.publish(po.Title, po.Msg)
}

func (ps *Pubsub) receive()(PubsubObject, error){
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()
	defer ps.popElement()
	if len(ps.queue) == 0 {
		return PubsubObject{}, errors.New("Can't receive new object")
	}

	return ps.queue[0], nil
}


//If we get new message, construct new notification
func (ps *Pubsub) publish(title, msg string){
	item, ok := ps.inner[title]
	if !ok {
		//TODO set some error
	} else {
		ps.queue = append(ps.queue, PubsubObject{once:item.once, msg:msg})
	}
}

//Pop one element from the queue
func (ps *Pubsub) popElement(){
	if len(ps.queue) > 0 {
		ps.queue = append(ps.queue[:0], ps.queue[1:]...)
	}
}

//clearQueue from all elements
func (ps *Pubsub) clearQueue(){
	ps.mutex.RLock()
	ps.queue = ps.queue[:0]
	defer ps.mutex.RUnlock()
}
