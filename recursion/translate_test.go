package recursion

//import "fmt"
import "testing"

func TestPrefix2Postfix(t *testing.T) {
	testPrefix2Infix(t, Prefix2InfixRecursive, "prefix to infix recursive")
	testPrefix2Postfix(t, Prefix2PostfixRecursive, "prefix to postfix recursive")
	testPrefix2Infix(t, Prefix2Infix, "prefix to infix with stack")
	testPrefix2Postfix(t, Prefix2Postfix, "prefix to postfix with stack")
	testInfix2Prefix(t, Infix2PrefixRecursive, "infix to prefix recursive")
	testInfix2Postfix(t, Infix2PostfixRecursive, "infix to postfix recursive")
	testInfix2Prefix(t, Infix2Prefix, "infix to prefix with a stack")
	testInfix2Postfix(t, Infix2Postfix, "infix to postfix with a stack")
	testPostfix2Prefix(t, Postfix2PrefixRecursive, "postfix to prefix recursive")
	testPostfix2Infix(t, Postfix2InfixRecursive, "postfix to infix recursive")
	testPostfix2Prefix(t, Postfix2Prefix, "postfix to prefix with a stack")
	testPostfix2Infix(t, Postfix2Infix, "postfix to infix with a stack")
}

func testPrefix2Postfix(t *testing.T, translate func(string) (string, error), name string) {
	if result, err := translate(""); err == nil {
		t.Errorf("%v fails on empty string with result %v", name, result)
	}
	if result, err := translate("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "5" {
		t.Errorf("%v fails on argument 5 with result %v", name, result)
	}
	if result, err := translate("56"); err == nil {
		t.Errorf("%v fails on 56 with result %v", name, result)
	}
	if result, err := translate("+56"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "56+" {
		t.Errorf("%v fails on +56 with result %v", name, result)
	}
	if result, err := translate("5+6"); err == nil {
		t.Errorf("%v fails on 5+6 with result %v", name, result)
	}
	if result, err := translate("*+56-72"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "56+72-*" {
		t.Errorf("%v fails on *+56-72 with result %v", name, result)
	}
	if result, err := translate("-*+7229"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "72+2*9-" {
		t.Errorf("%v fails on -*+7229 with result %v", name, result)
	}
	if result, err := translate("+-8*72%-643"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "872*-64-3%+" {
		t.Errorf("%v fails on +-8*72%-643 with result %v", name, result)
	}
}

func testPrefix2Infix(t *testing.T, translate func(string) (string, error), name string) {
	if result, err := translate(""); err == nil {
		t.Errorf("%v fails on empty string with result %v", name, result)
	}
	if result, err := translate("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "5" {
		t.Errorf("%v fails on argument 5 with result %v", name, result)
	}
	if result, err := translate("56"); err == nil {
		t.Errorf("%v fails on 56 with result %v", name, result)
	}
	if result, err := translate("+56"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "(5+6)" {
		t.Errorf("%v fails on +56 with result %v", name, result)
	}
	if result, err := translate("5+6"); err == nil {
		t.Errorf("%v fails on 5+6 with result %v", name, result)
	}
	if result, err := translate("*+56-72"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "((5+6)*(7-2))" {
		t.Errorf("%v fails on *+56-72 with result %v", name, result)
	}
	if result, err := translate("-*+7229"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "(((7+2)*2)-9)" {
		t.Errorf("%v fails on -*+7229 with result %v", name, result)
	}
	if result, err := translate("+-8*72%-643"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "((8-(7*2))+((6-4)%3))" {
		t.Errorf("%v fails on +-8*72%-643 with result %v", name, result)
	}
}

func testInfix2Prefix(t *testing.T, translate func(string) (string, error), name string) {
	if result, err := translate(""); err == nil {
		t.Errorf("%v fails on empty string with result %v", name, result)
	}
	if result, err := translate("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "5" {
		t.Errorf("%v fails on argument 5 with result %v", name, result)
	}
	if result, err := translate("56"); err == nil {
		t.Errorf("%v fails on 56 with result %v", name, result)
	}
	if result, err := translate("5+6"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "+56" {
		t.Errorf("%v fails on 5+6 with result %v", name, result)
	}
	if result, err := translate("+56"); err == nil {
		t.Errorf("%v fails on +56 with result %v", name, result)
	}
	if result, err := translate("(56)"); err == nil {
		t.Errorf("%v fails on (56) with result %v", name, result)
	}
	if result, err := translate("(5+6"); err == nil {
		t.Errorf("%v fails on (5+6 with result %v", name, result)
	}
	if result, err := translate("(5+6)*(7-2)"); err != nil {
		t.Errorf("%v fails on (5+6)*(7-2): %v", name, err)
	} else if result != "*+56-72" {
		t.Errorf("%v fails on (5+6)*(7-2) with result %v", name, result)
	}
	if result, err := translate("5+6*7-2"); err != nil {
		t.Errorf("%v fails on 5+6*7-2: %v", name, err)
	} else if result != "-*+5672" {
		t.Errorf("%v fails on 5+6*7-2 with result %v", name, result)
	}
}

func testInfix2Postfix(t *testing.T, translate func(string) (string, error), name string) {
	if result, err := translate(""); err == nil {
		t.Errorf("%v fails on empty string with result %v", name, result)
	}
	if result, err := translate("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "5" {
		t.Errorf("%v fails on argument 5 with result %v", name, result)
	}
	if result, err := translate("56"); err == nil {
		t.Errorf("%v fails on 56 with result %v", name, result)
	}
	if result, err := translate("5+6"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "56+" {
		t.Errorf("%v fails on 5+6 with result %v", name, result)
	}
	if result, err := translate("+56"); err == nil {
		t.Errorf("%v fails on +56 with result %v", name, result)
	}
	if result, err := translate("(56)"); err == nil {
		t.Errorf("%v fails on (56) with result %v", name, result)
	}
	if result, err := translate("(5+6"); err == nil {
		t.Errorf("%v fails on (5+6 with result %v", name, result)
	}
	if result, err := translate("(5+6)*(7-2)"); err != nil {
		t.Errorf("%v fails on (5+6)*(7-2): %v", name, err)
	} else if result != "56+72-*" {
		t.Errorf("%v fails on (5+6)*(7-2) with result %v", name, result)
	}
	if result, err := translate("5+6*7-2"); err != nil {
		t.Errorf("%v fails on 5+6*7-2: %v", name, err)
	} else if result != "56+7*2-" {
		t.Errorf("%v fails on 5+6*7-2 with result %v", name, result)
	}
}

func testPostfix2Prefix(t *testing.T, translate func(string) (string, error), name string) {
	if result, err := translate(""); err == nil {
		t.Errorf("%v fails on empty string with result %v", name, result)
	}
	if result, err := translate("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "5" {
		t.Errorf("%v fails on argument 5 with result %v", name, result)
	}
	if result, err := translate("56"); err == nil {
		t.Errorf("%v fails on 56 with result %v", name, result)
	}
	if result, err := translate("56+"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "+56" {
		t.Errorf("%v fails on 56+ with result %v", name, result)
	}
	if result, err := translate("5+6"); err == nil {
		t.Errorf("%v fails on 5+6 with result %v", name, result)
	}
	if result, err := translate("567*"); err == nil {
		t.Errorf("%v fails on 567* with result %v", name, result)
	}
	if result, err := translate("+56"); err == nil {
		t.Errorf("%v fails on +56 with result %v", name, result)
	}
	if result, err := translate("567*+"); err != nil {
		t.Errorf("%v fails on 567*+: %v", name, err)
	} else if result != "+5*67" {
		t.Errorf("%v fails on 567*+ with result %v", name, result)
	}
	if result, err := translate("56*71-3+*"); err != nil {
		t.Errorf("%v fails on 56*71-3+*: %v", name, err)
	} else if result != "**56+-713" {
		t.Errorf("%v fails on 56*71-3+* with result %v", name, result)
	}
	if result, err := translate("12+31/43%+42**+"); err != nil {
		t.Errorf("%v fails on 12+31/43%+42**+: %v", name, err)
	} else if result != "++12*+/31%43*42" {
		t.Errorf("%v fails on 12+31/43%+42**+ with result %v", name, result)
	}
}

func testPostfix2Infix(t *testing.T, translate func(string) (string, error), name string) {
	if result, err := translate(""); err == nil {
		t.Errorf("%v fails on empty string with result %v", name, result)
	}
	if result, err := translate("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "5" {
		t.Errorf("%v fails on argument 5 with result %v", name, result)
	}
	if result, err := translate("56"); err == nil {
		t.Errorf("%v fails on 56 with result %v", name, result)
	}
	if result, err := translate("56+"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if result != "(5+6)" {
		t.Errorf("%v fails on 56+ with result %v", name, result)
	}
	if result, err := translate("5+6"); err == nil {
		t.Errorf("%v fails on 5+6 with result %v", name, result)
	}
	if result, err := translate("567*"); err == nil {
		t.Errorf("%v fails on 567* with result %v", name, result)
	}
	if result, err := translate("+56"); err == nil {
		t.Errorf("%v fails on +56 with result %v", name, result)
	}
	if result, err := translate("567*+"); err != nil {
		t.Errorf("%v fails on 567*+: %v", name, err)
	} else if result != "(5+(6*7))" {
		t.Errorf("%v fails on 567*+ with result %v", name, result)
	}
	if result, err := translate("56*71-3+*"); err != nil {
		t.Errorf("%v fails on 56*71-3+*: %v", name, err)
	} else if result != "((5*6)*((7-1)+3))" {
		t.Errorf("%v fails on 56*71-3+* with result %v", name, result)
	}
	if result, err := translate("12+31/43%+42**+"); err != nil {
		t.Errorf("%v fails on 12+31/43%+42**+: %v", name, err)
	} else if result != "((1+2)+(((3/1)+(4%3))*(4*2)))" {
		t.Errorf("%v fails on 12+31/43%+42**+ with result %v", name, result)
	}
}
