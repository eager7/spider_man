package middleware

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Entity interface {
	Id() uint32
}

type Pool interface {
	Take() (Entity, error)
	Return(entity Entity) error
	Total() uint32
	Used() uint32
}

type myPool struct {
	mutex sync.Mutex
	total uint32
	etype reflect.Type
	genEntity func()Entity
	container chan Entity
	idContainer map[uint32]bool
}

func NewPool(total uint32, entityType reflect.Type, genEntity func()Entity) (Pool, error) {
	if total == 0 {
		return nil, errors.New(fmt.Sprintf("The poll can't be initialized! (total=%d)\n", total))
	}
	size := int(total)
	container := make(chan Entity, size)
	idContainer := make(map[uint32]bool)
	for i := 0; i < size; i ++ {
		newEntity := genEntity()
		if reflect.TypeOf(newEntity) != entityType {
			return nil, errors.New(fmt.Sprintf("The type of result of function genEntity() is Not %s!\n", entityType))
		}
		container <- newEntity
		idContainer[newEntity.Id()] = true
	}
	pool := &myPool{sync.Mutex{},total, entityType, genEntity, container, idContainer}
	return pool, nil
}

func (p *myPool)Take() (Entity, error){
	entity, ok := <-p.container
	if !ok {
		return nil, errors.New("The inner container is invalid!")
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.idContainer[entity.Id()] = false
	return entity, nil
}

func (p *myPool)Return(entity Entity) error {
	if entity == nil {
		return errors.New("The returning entity is invalid!")
	}
	if p.etype != reflect.TypeOf(entity) {
		return errors.New(fmt.Sprintf("The type of returning entity is NOT %s!\n", p.etype))
	}
	entityId := entity.Id()
	casResult := p.compareAndSetForIdContainer(entityId, false, true)
	if casResult == 1 {
		p.container <- entity
		return nil
	} else if casResult == 0 {
		return errors.New(fmt.Sprintf("The entity (id=%d) is already in the pool!\n", entityId))
	} else {
		return errors.New(fmt.Sprintf("The entity (id=%d) is illegal!\n", entityId))
	}
}

func (p *myPool)Total()uint32 {
	return p.total
}

func (p *myPool)Used()uint32{
	return p.total - uint32(len(p.container))
}

func (p *myPool)compareAndSetForIdContainer (entityId uint32, oldValue bool, newValue bool) int8 {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	v, ok := p.idContainer[entityId]
	if !ok {
		return -1
	}
	if v != oldValue {
		return 0
	}
	p.idContainer[entityId] = newValue
	return 1
}