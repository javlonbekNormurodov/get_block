package eth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"eth_tracker/internal/domain"
)

func TestEthClient_GetBlockByNumber(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "jsonrpc":"2.0",
            "result":{
                "number":"0x10d4f",
                "hash":"0x123",
                "parentHash":"0x456",
                "transactions":["0x789"]
            },
            "id":1
        }`))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	ethClient := &Client{rpcURL: server.URL}
	block, err := ethClient.GetBlockByNumber(69007)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	expected := &domain.Block{
		Number:     "0x10d4f",
		Hash:       "0x123",
		ParentHash: "0x456",
		Txns:       []string{"0x789"},
	}
	if !reflect.DeepEqual(block, expected) {
		t.Errorf("expected block %+v, got %+v", expected, block)
	}
}
