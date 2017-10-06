package cpw

// based on https://gist.github.com/rday/3504674
type InitFunction func() (interface{}, error)

type Cpw struct {
	conn chan interface{}
}

/**
 Call the init function size times. If the init function fails during any call, then
 the creation of the pool is considered a failure.
 We call the same function size times to make sure each connection shares the same
 state.
*/
func (p *Cpw) InitPool(size int, initfn InitFunction) error {
	// Create a buffered channel allowing size senders
	p.conn = make(chan interface{}, size)
	for x := 0; x < size; x++ {
		conn, err := initfn()
		if err != nil {
			return err
		}

		// If the init function succeeded, add the connection to the channel
		p.conn <- conn
	}
	return nil
}

func (p *Cpw) GetConnection() interface{} {
	return <-p.conn
}

func (p *Cpw) ReleaseConnection(conn interface{}) {
	p.conn <- conn
}

func NewCpw() *Cpw {
	return &Cpw{}
}
