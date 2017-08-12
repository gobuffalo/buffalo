package pop

type beforeSaveable interface {
	BeforeSave(*Connection) error
}

func (m *Model) beforeSave(c *Connection) error {
	if x, ok := m.Value.(beforeSaveable); ok {
		return x.BeforeSave(c)
	}
	return nil
}

type beforeCreateable interface {
	BeforeCreate(*Connection) error
}

func (m *Model) beforeCreate(c *Connection) error {
	if x, ok := m.Value.(beforeCreateable); ok {
		return x.BeforeCreate(c)
	}
	return nil
}

type beforeUpdateable interface {
	BeforeUpdate(*Connection) error
}

func (m *Model) beforeUpdate(c *Connection) error {
	if x, ok := m.Value.(beforeUpdateable); ok {
		return x.BeforeUpdate(c)
	}
	return nil
}

type beforeDestroyable interface {
	BeforeDestroy(*Connection) error
}

func (m *Model) beforeDestroy(c *Connection) error {
	if x, ok := m.Value.(beforeDestroyable); ok {
		return x.BeforeDestroy(c)
	}
	return nil
}

type afterDestroyable interface {
	AfterDestroy(*Connection) error
}

func (m *Model) afterDestroy(c *Connection) error {
	if x, ok := m.Value.(afterDestroyable); ok {
		return x.AfterDestroy(c)
	}
	return nil
}

type afterUpdateable interface {
	AfterUpdate(*Connection) error
}

func (m *Model) afterUpdate(c *Connection) error {
	if x, ok := m.Value.(afterUpdateable); ok {
		return x.AfterUpdate(c)
	}
	return nil
}

type afterCreateable interface {
	AfterCreate(*Connection) error
}

func (m *Model) afterCreate(c *Connection) error {
	if x, ok := m.Value.(afterCreateable); ok {
		return x.AfterCreate(c)
	}
	return nil
}

type afterSaveable interface {
	AfterSave(*Connection) error
}

func (m *Model) afterSave(c *Connection) error {
	if x, ok := m.Value.(afterSaveable); ok {
		return x.AfterSave(c)
	}
	return nil
}
