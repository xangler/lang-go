package main

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type Obj struct {
	Age int
}

func newObj(L *lua.LState) int {
	person := &Obj{L.CheckInt(1)}
	ud := L.NewUserData()
	ud.Value = person
	L.SetMetatable(ud, L.GetTypeMetatable("Obj"))
	L.Push(ud)
	return 1
}

func checkObj(L *lua.LState) *Obj {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Obj); ok {
		return v
	}
	L.ArgError(1, "Obj expect")
	return nil
}

func objSetGetAge(L *lua.LState) int {
	p := checkObj(L)
	if L.GetTop() == 2 {
		p.Age = L.CheckInt(2)
		return 0
	}
	L.Push(lua.LNumber(p.Age))
	return 1
}

var objMethods = map[string]lua.LGFunction{
	"age": objSetGetAge,
}

func main() {
	fmt.Println("lua go test start..")
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("gfib", load)
	if err := L.DoFile("script/compose.lua"); err != nil {
		panic(err)
	}

	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("fib"),
		NRet:    1,
		Protect: true,
	}, lua.LNumber(10))
	if err != nil {
		panic(err)
	}

	ret := L.Get(-1)
	L.Pop(1)
	res, ok := ret.(lua.LNumber)
	if ok {
		fmt.Println(res)
	} else {
		fmt.Printf("something wrong, %v\n", res)
	}
	fmt.Println("lua go test end..")
}

func load(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.SetField(mod, "name", lua.LString("glib"))
	L.Push(mod)

	mt := L.NewTypeMetatable("Obj")
	L.SetGlobal("Obj", mt)
	L.SetField(mt, "new", L.NewFunction(newObj))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), objMethods))
	L.Push(mt)
	return 2
}

var exports = map[string]lua.LGFunction{
	"fib": gfib,
}

func gfib(L *lua.LState) int {
	lv := L.ToInt(1)
	fmt.Printf("go fib start %v\n", lv)
	out := fib(lv)
	L.Push(lua.LNumber(out))
	fmt.Printf("go fib end %v\n", out)
	return 1
}

func fib(n int) int {
	if n < 2 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
