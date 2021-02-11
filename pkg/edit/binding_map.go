package edit

import (
	"errors"
	"sort"

	"github.com/xiaq/persistent/hashmap"
	"src.elv.sh/pkg/eval"
	"src.elv.sh/pkg/eval/vals"
	"src.elv.sh/pkg/parse"
	"src.elv.sh/pkg/ui"
)

var errValueShouldBeFn = errors.New("value should be function")

// bindingMap is a special Map that converts its key to ui.Key and ensures
// that its values satisfy eval.CallableValue.
type bindingMap struct {
	hashmap.Map
}

var emptyBindingMap = bindingMap{vals.EmptyMap}

// Repr returns the representation of the binding table as if it were an
// ordinary map keyed by strings.
func (bt bindingMap) Repr(indent int) string {
	var keys ui.Keys
	for it := bt.Map.Iterator(); it.HasElem(); it.Next() {
		k, _ := it.Elem()
		keys = append(keys, k.(ui.Key))
	}
	sort.Sort(keys)

	builder := vals.NewMapReprBuilder(indent)

	for _, k := range keys {
		v, _ := bt.Map.Index(k)
		builder.WritePair(parse.Quote(k.String()), indent+2, vals.Repr(v, indent+2))
	}

	return builder.String()
}

// Index converts the index to ui.Key and uses the Index of the inner Map.
func (bt bindingMap) Index(index interface{}) (interface{}, error) {
	key, err := toKey(index)
	if err != nil {
		return nil, err
	}
	return vals.Index(bt.Map, key)
}

func (bt bindingMap) HasKey(k interface{}) bool {
	_, ok := bt.Map.Index(k)
	return ok
}

func (bt bindingMap) GetKey(k ui.Key) eval.Callable {
	v, ok := bt.Map.Index(k)
	if !ok {
		panic("get called when key not present")
	}
	return v.(eval.Callable)
}

// Assoc converts the index to ui.Key, ensures that the value is CallableValue,
// uses the Assoc of the inner Map and converts the result to a BindingTable.
func (bt bindingMap) Assoc(k, v interface{}) (interface{}, error) {
	key, err := toKey(k)
	if err != nil {
		return nil, err
	}
	f, ok := v.(eval.Callable)
	if !ok {
		return nil, errValueShouldBeFn
	}
	map2 := bt.Map.Assoc(key, f)
	return bindingMap{map2}, nil
}

// Dissoc converts the key to ui.Key and calls the Dissoc method of the inner
// map.
func (bt bindingMap) Dissoc(k interface{}) interface{} {
	key, err := toKey(k)
	if err != nil {
		// Key is invalid; dissoc is no-op.
		return bt
	}
	return bindingMap{bt.Map.Dissoc(key)}
}

func makeBindingMap(raw hashmap.Map) (bindingMap, error) {
	converted := vals.EmptyMap
	for it := raw.Iterator(); it.HasElem(); it.Next() {
		k, v := it.Elem()
		f, ok := v.(eval.Callable)
		if !ok {
			return emptyBindingMap, errValueShouldBeFn
		}
		key, err := toKey(k)
		if err != nil {
			return bindingMap{}, err
		}
		converted = converted.Assoc(key, f)
	}

	return bindingMap{converted}, nil
}
