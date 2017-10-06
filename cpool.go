package cpool

import "github.com/mkrufky/cpool/cpw"
import "fmt"

type Cpool struct {
	cpw *cpw.Cpw
	idp *cpw.Cpw
	m map[int] interface{}
}

func (p *Cpool) InitPool(size int, initfn cpw.InitFunction) error {
	ret := p.cpw.InitPool(size, initfn)
	if ret == nil {
		sz := 0
		p.idp.InitPool(size, func() (interface{}, error) {
			id := sz
			sz++
			return id, nil
		})
	}
	return ret
}

func (p *Cpool) Connect() (int, interface{}) {
	conn := p.cpw.GetConnection()
	id := p.idp.GetConnection().(int)

	p.m[id] = conn
	return id, conn
}

func (p *Cpool) GetConnection(id int) (interface{}, error) {
	c, ok := p.m[id]
	if !ok {
		panic(fmt.Errorf("GetConnection: invalid map key: %d", id))
	}
	return c, nil
}

func (p *Cpool) ReleaseConnection(id int) error {
	c, ok := p.m[id]
	if !ok {
		panic(fmt.Errorf("ReleaseConnection: invalid map key: %d", id))
	}
	delete(p.m, id)
	p.idp.ReleaseConnection(id)
	p.cpw.ReleaseConnection(c)
	return nil
}

func NewCpool() *Cpool {
	return &Cpool{cpw: cpw.NewCpw(), idp: cpw.NewCpw(), m: make(map[int]interface{})}
}
