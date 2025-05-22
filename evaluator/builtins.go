package evaluator

import (
	"monkey/object"
)

var builtins = map[string]*object.Builtin{
	"len":   object.GetBuiltinByName("len"),
	"first": object.GetBuiltinByName("first"),
	"last":  object.GetBuiltinByName("last"),
	//bug todo: 对一个数组`[1, 2, 3, 4]`首先使用`push`操作`push(t, 5)`，再使用`rest`操作`rest(t)`，就会得到一个空数组`[2, 3, 4]`，而不是`[2, 3, 4, 5]`
	"rest": object.GetBuiltinByName("rest"),
	"push": object.GetBuiltinByName("push"),
	"puts": object.GetBuiltinByName("puts"),
}
