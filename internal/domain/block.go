package domain

type Block struct {
	Number     string
	Hash       string
	ParentHash string
	Txns       []string
}
