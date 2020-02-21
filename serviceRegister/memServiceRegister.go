package serviceRegister

import "vsFactoryAPIEngine/service"

//
// 将注册表放在内存中
//
type MemServiceRegister struct {
	// 存放UUID名称 -> 用户服务对象池
	s2u map[string]IUserServicePool
}

// 单例对象
var memRegister *MemServiceRegister

func init() {
	memRegister = &MemServiceRegister{s2u: make(map[string]IUserServicePool)}
}

func GetMemServiceRegister() *MemServiceRegister {
	return memRegister
}

func (m *MemServiceRegister) Register(uuid string, serviceName string, server service.Server) {
	if d, ok := m.s2u[uuid]; !ok {
		// 之前没有注册过, 新建对象池
		pool := NewMemUserServicePool()
		pool.RegisterServe(serviceName, server)
		m.s2u[uuid] = pool
	} else {
		// 注册服务
		d.RegisterServe(serviceName, server)
	}
}

func (m *MemServiceRegister) GetPool(uid string) (IUserServicePool, bool) {
	pool, ok := m.s2u[uid]
	return pool, ok
}

func (m *MemServiceRegister) Valid(uid string, servername string) bool {
	pool, ok := m.s2u[uid]
	if !ok {
		return false
	}
	_, ok2 := pool.GetServe(servername)
	return ok2
}

//
// 服务对象池
//
type MemUserServicePool struct {
	spools map[string]service.Server
}

func NewMemUserServicePool() *MemUserServicePool {
	return &MemUserServicePool{spools: make(map[string]service.Server)}
}

func (m *MemUserServicePool) RegisterServe(servername string, server service.Server) {
	m.spools[servername] = server
}

func (m *MemUserServicePool) GetServe(servername string) (service.Server, bool) {
	s, ok := m.spools[servername]
	return s, ok
}
