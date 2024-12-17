package pool

// PooledObject 需要放进池中的对象
type PooledObject interface {
	New()   //新建一个对象
	Reset() //放进池时，进行重置
}

// PooledObjectFactory 创建工厂
type PooledObjectFactory interface {
	Create() (PooledObject, error)
}

// Pool 实现接口
type Pool interface {
	Get() (PooledObject, error)
	Put(obj PooledObject) error
}
