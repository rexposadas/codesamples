// This essentially defers further processing asynchronously to Process()
func (o *SteamOrder) FireAtSteam() error {

	// initiate a non-blocking send
	select {
	case Cannons.Steam <- o:
		return nil // success
	default:
	}

	// return an error if our queue is full
	return fmt.Errorf("The Steam Cannon is full and has rejected %v", o)

}
