/////////////////////////////////////////////////////////////////////////
// The following are code samples integration with Steam using GoLang.
//
// The code samples demostrates the following:
// 1. How to create high throughput API services.
// 2. Benefits of a highly concurrent language like Golang.
//
// Performance:
// 600+/sec throughput on a linux VM with 512MB of memorry running on
// a Windows 7 machine.
//

// The focal point of this ample is this channel.
// 1. We load the cannon with and a steam order.
// 2. Steam order have states associated with them.
// 3. Depending on their state we fire the cannon at Steam
func (o *SteamOrder) LoadCannon() error {

	// initiate a non-blocking send
	select {
	case Cannons.Steam <- o:
		return nil // success
	default:
	}

	// return an error if our queue is full
	return fmt.Errorf("The Steam Cannon is full and has rejected %v", o)
}

///////////////////////////////////////////////////////////////////////////////
// Grab orders from the channel. Depending on the status of the order,
// we either initialize or finalize and order in Steam
//
// The number of workers is tied to the size of the queue,
// which will block once the limit is hit.
func steamWorker(id int) {
	fmt.Printf("STEAM:WORKER(%d) START\n", id)

	for {
		order, ok := <-Cannons.Steam
		if !ok {
			return
		}

		order.Process()
	}

	fmt.Printf("STEAM:WORKER(%d) STOP\n", id)
}

// Initializes the channel and then spawns & invokes our worker pool
func SpawnWorkers() {
	Cannons.Steam = make(chan *SteamOrder, STEAM_QUEUE_SIZE)

	for i := 0; i < STEAM_WORKERS; i++ {
		go steamWorker(i)
	}
}

// This function is ran by the workers.
// Depending on the state of each order an action is taken.
// This is how we use Golang's builtin queue system. This is also
// a way to simulate an assembly line for your processes.
func (o *SteamOrder) Process() {

	switch o.State {
	case STEAM_ORDER_STATE_NEW:
		o.initialize()
		return
	case STEAM_ORDER_STATE_INITIALIZED:
		o.finalize()
		return
	}

	// we have nothing to return and nothing to do here
	fmt.Printf("Attempt to process an order in an invalid state: %v", o)
}

