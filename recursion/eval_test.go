package recursion

//import "fmt"
import "testing"

func TestPrefixEval(t *testing.T) {
	testPrefixEvalFunction(t, EvalPrefixRecursive, "prefix recursive")
	testPrefixEvalFunction(t, EvalPrefixStack, "prefix stack")
}

func testPrefixEvalFunction(t *testing.T, eval func(string) (int, error), name string) {
	if val, err := eval(""); err == nil {
		t.Errorf("%v fails on empty string with value %v", name, val)
	}
	if val, err := eval("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 5 {
		t.Errorf("%v fails on argument 5 with value %v", name, val)
	}
	if val, err := eval("56"); err == nil {
		t.Errorf("%v fails on 56 with value %v", name, val)
	}
	if val, err := eval("+56"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 11 {
		t.Errorf("%v fails on +56 with value %v", name, val)
	}
	if val, err := eval("5+6"); err == nil {
		t.Errorf("%v fails on 5+6 with value %v", name, val)
	}
	if val, err := eval("*+56-72"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 55 {
		t.Errorf("%v fails on *+56-72 with value %v", name, val)
	}
	if val, err := eval("-*+7229"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 9 {
		t.Errorf("%v fails on -*+7229 with value %v", name, val)
	}
	if val, err := eval("+-8*72%-643"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != -4 {
		t.Errorf("%v fails on +-8*72%-643 with value %v", name, val)
	}
}

func TestInfixEval(t *testing.T) {
	testInfixEvalFunction(t, EvalInfixRecursive, "infix recursive")
	testInfixEvalFunction(t, EvalInfixStack, "infix stack")
}

func testInfixEvalFunction(t *testing.T, eval func(string) (int, error), name string) {
	if val, err := eval(""); err == nil {
		t.Errorf("%v fails on empty string with value %v", name, val)
	}
	if val, err := eval("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 5 {
		t.Errorf("%v fails on argument 5 with value %v", name, val)
	}
	if val, err := eval("56"); err == nil {
		t.Errorf("%v fails on 56 with value %v", name, val)
	}
	if val, err := eval("5+6"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 11 {
		t.Errorf("%v fails on 5+6 with value %v", name, val)
	}
	if val, err := eval("+56"); err == nil {
		t.Errorf("%v fails on 5+6 with value %v", name, val)
	}
	if val, err := eval("(56)"); err == nil {
		t.Errorf("%v fails on (56) with value %v", name, val)
	}
	if val, err := eval("(5+6"); err == nil {
		t.Errorf("%v fails on (5+6 with value %v", name, val)
	}
	if val, err := eval("(5+6)*(7-2)"); err != nil {
		t.Errorf("%v fails on (5+6)*(7-2): %v", name, err)
	} else if val != 55 {
		t.Errorf("%v fails on (5+6)*(7-2) with value %v", name, val)
	}
	if val, err := eval("5+6*7-2"); err != nil {
		t.Errorf("%v fails on 5+6*7-2: %v", name, err)
	} else if val != 75 {
		t.Errorf("%v fails on 5+6*7-2 with value %v", name, val)
	}
}

func TestPostfixEval(t *testing.T) {
	testPostfixEvalFunction(t, EvalPostfixRecursive, "postfix recursive")
	testPostfixEvalFunction(t, EvalPostfixStack, "postfix stack")
}

func testPostfixEvalFunction(t *testing.T, eval func(string) (int, error), name string) {
	if val, err := eval(""); err == nil {
		t.Errorf("%v fails on empty string with value %v", name, val)
	}
	if val, err := eval("5"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 5 {
		t.Errorf("%v fails on argument 5 with value %v", name, val)
	}
	if val, err := eval("56"); err == nil {
		t.Errorf("%v fails on 56 with value %v", name, val)
	}
	if val, err := eval("56+"); err != nil {
		t.Errorf("%v fails: %v", name, err)
	} else if val != 11 {
		t.Errorf("%v fails on 56+ with value %v", name, val)
	}
	if val, err := eval("5+6"); err == nil {
		t.Errorf("%v fails on 5+6 with value %v", name, val)
	}
	if val, err := eval("567*"); err == nil {
		t.Errorf("%v fails on 567* with value %v", name, val)
	}
	if val, err := eval("+56"); err == nil {
		t.Errorf("%v fails on +56 with value %v", name, val)
	}
	if val, err := eval("567*+"); err != nil {
		t.Errorf("%v fails on 567*+: %v", name, err)
	} else if val != 47 {
		t.Errorf("%v fails on 567*+ with value %v", name, val)
	}
	if val, err := eval("56*71-3+*"); err != nil {
		t.Errorf("%v fails on 56*71-3+*: %v", name, err)
	} else if val != 270 {
		t.Errorf("%v fails on 56*71-3+* with value %v", name, val)
	}
	if val, err := eval("12+31/43%+42**+"); err != nil {
		t.Errorf("%v fails on 12+31/43%+42**+: %v", name, err)
	} else if val != 35 {
		t.Errorf("%v fails on 12+31/43%+42**+ with value %v", name, val)
	}
}
