all : parser.go

parser.go : parser.go.y
	goyacc -o $@ parser.go.y
