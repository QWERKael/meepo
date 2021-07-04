package taskchain

import "testing"

func TestNewTaskUnit(t *testing.T) {
	a := NewTaskUnit("a", "a", "a", "a")
	b := NewTaskUnit("b", "b", "b", "b")
	c := NewTaskUnit("c", "c", "c", "c")
	d := NewTaskUnit("d", "d", "d", "d")
	a.TaskUnitDisplay()
	b.TaskUnitDisplay()
	c.TaskUnitDisplay()
	d.TaskUnitDisplay()
}