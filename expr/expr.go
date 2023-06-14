package expr

type Expr interface {
	Expr()
}

type DoC struct {
	Exprs []Expr
}

func (m DoC) Expr() {}

type MapC struct {
	Binds map[Expr]Expr
}

func (m MapC) Expr() {}

type ListC struct {
	Values []Expr
}

func (l ListC) Expr() {}

type StrC struct {
	Value string
}

func (s StrC) Expr() {}

type NumC struct {
	Value float64
}

func (num NumC) Expr() {}

type LamC struct {
	Params []Expr
	Body   Expr
}

func (num LamC) Expr() {}

type Call struct {
	Proc Expr
	Args []Expr
}

func (call Call) Expr() {}

type IdC struct {
	Value string
}

func (id IdC) Expr() {}

type Binop struct {
	Op    string
	Left  Expr
	Right Expr
}

func (bin Binop) Expr() {}

type Group struct {
	Body Expr
}

func (group Group) Expr() {}

type Let struct {
	Id    Expr
	Value Expr
}

func (let Let) Expr() {}

type Nil struct{}

func (nil Nil) Expr() {}

type If struct {
	Cond Expr
	Then Expr
	Else Expr
}

func (nil If) Expr() {}
