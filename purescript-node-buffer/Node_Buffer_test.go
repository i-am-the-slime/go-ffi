package purescript_node_buffer

import (
	"testing"

	. "github.com/purescript-native/go-runtime"
)

func TestCreate(t *testing.T) {
	exports := Foreign("Node.Buffer")
	create := exports["create"].(func(Any) Any)
	size := exports["size"].(func(Any) Any)
	
	effect := create(10).(func() Any)
	buf := effect()
	
	sizeEffect := size(buf).(func() Any)
	actualSize := sizeEffect().(int)
	
	if actualSize != 10 {
		t.Errorf("Expected size 10, got %d", actualSize)
	}
}

func TestFromString(t *testing.T) {
	exports := Foreign("Node.Buffer")
	fromString := exports["fromString"].(func(Any, Any) Any)
	toString := exports["toString"].(func(Any, Any) Any)
	
	effect := fromString("Hello", "UTF8").(func() Any)
	buf := effect()
	
	toStrEffect := toString("UTF8", buf).(func() Any)
	result := toStrEffect().(string)
	
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

func TestFromStringBase64(t *testing.T) {
	exports := Foreign("Node.Buffer")
	fromString := exports["fromString"].(func(Any, Any) Any)
	toString := exports["toString"].(func(Any, Any) Any)
	
	// "Hello" in base64 is "SGVsbG8="
	effect := fromString("SGVsbG8=", "Base64").(func() Any)
	buf := effect()
	
	toStrEffect := toString("UTF8", buf).(func() Any)
	result := toStrEffect().(string)
	
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

func TestFromArray(t *testing.T) {
	exports := Foreign("Node.Buffer")
	fromArray := exports["fromArray"].(func(Any) Any)
	toArray := exports["toArray"].(func(Any) Any)
	
	arr := []Any{72, 101, 108, 108, 111} // "Hello" in ASCII
	
	effect := fromArray(arr).(func() Any)
	buf := effect()
	
	toArrEffect := toArray(buf).(func() Any)
	result := toArrEffect().([]Any)
	
	if len(result) != len(arr) {
		t.Fatalf("Expected length %d, got %d", len(arr), len(result))
	}
	
	for i, v := range arr {
		if result[i] != v {
			t.Errorf("At index %d: expected %v, got %v", i, v, result[i])
		}
	}
}

func TestConcat(t *testing.T) {
	exports := Foreign("Node.Buffer")
	fromString := exports["fromString"].(func(Any, Any) Any)
	concat := exports["concat"].(func(Any) Any)
	toString := exports["toString"].(func(Any, Any) Any)
	
	buf1 := fromString("Hello", "UTF8").(func() Any)()
	buf2 := fromString(" ", "UTF8").(func() Any)()
	buf3 := fromString("World", "UTF8").(func() Any)()
	
	concatEffect := concat([]Any{buf1, buf2, buf3}).(func() Any)
	result := concatEffect()
	
	toStrEffect := toString("UTF8", result).(func() Any)
	str := toStrEffect().(string)
	
	if str != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", str)
	}
}

func TestSlice(t *testing.T) {
	exports := Foreign("Node.Buffer")
	fromString := exports["fromString"].(func(Any, Any) Any)
	slice := exports["slice"].(func(Any, Any, Any) Any)
	toString := exports["toString"].(func(Any, Any) Any)
	
	buf := fromString("Hello World", "UTF8").(func() Any)()
	
	sliceEffect := slice(0, 5, buf).(func() Any)
	sliced := sliceEffect()
	
	toStrEffect := toString("UTF8", sliced).(func() Any)
	result := toStrEffect().(string)
	
	if result != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result)
	}
}

func TestFill(t *testing.T) {
	exports := Foreign("Node.Buffer")
	create := exports["create"].(func(Any) Any)
	fill := exports["fill"].(func(Any, Any, Any, Any) Any)
	toArray := exports["toArray"].(func(Any) Any)
	
	buf := create(5).(func() Any)()
	
	fillEffect := fill(65, 0, 5, buf).(func() Any) // 'A' = 65
	fillEffect()
	
	arrEffect := toArray(buf).(func() Any)
	arr := arrEffect().([]Any)
	
	for i, v := range arr {
		if v != 65 {
			t.Errorf("At index %d: expected 65, got %v", i, v)
		}
	}
}

func TestCopy(t *testing.T) {
	exports := Foreign("Node.Buffer")
	fromString := exports["fromString"].(func(Any, Any) Any)
	create := exports["create"].(func(Any) Any)
	copyFn := exports["copy"].(func(Any, Any, Any, Any, Any) Any)
	toString := exports["toString"].(func(Any, Any) Any)
	
	src := fromString("Hello", "UTF8").(func() Any)()
	dst := create(10).(func() Any)()
	
	copyEffect := copyFn(0, 5, src, 0, dst).(func() Any)
	copied := copyEffect().(int)
	
	if copied != 5 {
		t.Errorf("Expected 5 bytes copied, got %d", copied)
	}
	
	toStrEffect := toString("UTF8", dst).(func() Any)
	result := toStrEffect().(string)
	
	if result[:5] != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result[:5])
	}
}

func TestWriteString(t *testing.T) {
	exports := Foreign("Node.Buffer")
	create := exports["create"].(func(Any) Any)
	writeString := exports["writeString"].(func(Any, Any, Any, Any, Any) Any)
	toString := exports["toString"].(func(Any, Any) Any)
	
	buf := create(10).(func() Any)()
	
	writeEffect := writeString("UTF8", 0, 5, "Hello", buf).(func() Any)
	written := writeEffect().(int)
	
	if written != 5 {
		t.Errorf("Expected 5 bytes written, got %d", written)
	}
	
	toStrEffect := toString("UTF8", buf).(func() Any)
	result := toStrEffect().(string)
	
	if result[:5] != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", result[:5])
	}
}

