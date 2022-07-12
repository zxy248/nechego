package app

import (
	"reflect"
	"testing"
)

func TestModsetAdd(t *testing.T) {
	ms := newModset()
	if ms.present(adminModifier) {
		t.Error("must not be present")
	}
	ms.add(adminModifier)
	if !ms.present(adminModifier) {
		t.Error("must be present")
	}
}

func TestModsetAddNoModifier(t *testing.T) {
	ms := newModset()
	ms.add(noModifier)
	if ms.present(noModifier) {
		t.Error("must not be present")
	}
}

func TestModsetModifiers(t *testing.T) {
	ms := newModset()
	ms.add(noModifier)

	wantMods := []*modifier{adminModifier, fisherModifier}
	for _, m := range wantMods {
		ms.add(m)
	}

	gotMods := ms.list()
	l := len(gotMods)
	if l != 2 {
		t.Error("got %v, want %v", l, 2)
	}
	if !reflect.DeepEqual(gotMods, wantMods) {
		t.Errorf("got %v, want %v", gotMods, wantMods)
	}
}
