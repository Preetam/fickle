package main

type Command struct {
	Type Operation
	Var1 string
	Var2 string
}

func GenerateCommand(op byte, var1 []byte, var2 []byte) Command {
	return Command{
		Type: Operation(op),
		Var1: string(var1),
		Var2: string(var2),
	}
}

func RunCommand(i *Instance, op byte, var1 []byte, var2 []byte) []byte {
	c := GenerateCommand(op, var1, var2)
	errByte, response := i.Execute(c)

	return append([]byte{errByte}, response...)
}
