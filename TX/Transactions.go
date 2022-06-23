package TX

type TxsData struct {
	Txs []*Tx
}

func (txs *TxsData) AddTx(t *Tx) {
	txs.Txs = append(txs.Txs, t)
}

func NewTxs() *TxsData {
	return &TxsData{}
}

/*
func (txs *TxsData) FindTx(id []byte) *Tx {
	for _, v := range txs.Txs {
		if v.EqualHash(id) {
			return v
		}
	}
	return nil
}

func (txs *TxsData) TxsString() string {
	s := ""
	for _, tx := range txs.Txs {
		s = s + string(tx.From) + string(tx.To) + string(tx.Sid) + strconv.Itoa(tx.Price) + "/n"
	}
	return s
}
*/
