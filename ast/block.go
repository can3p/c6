package ast

type Block struct {
	Stmts *StmtList
}

type BlockNode interface {
	MergeStmts(stmts *StmtList)
}

func NewBlock() *Block {
	return &Block{
		Stmts: &StmtList{},
	}
}

// Override the statements
func (self *Block) SetStmts(stms *StmtList) {
	self.Stmts = stms
}

func (self *Block) MergeBlock(block *Block) {
	self.Stmts.Stmts = append(self.Stmts.Stmts, block.Stmts.Stmts...)
}

func (self *Block) MergeStmts(stmts *StmtList) {
	for _, stm := range stmts.Stmts {
		self.Stmts.Append(stm)
	}
}
