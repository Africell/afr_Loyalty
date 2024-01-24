package Lendme

import (
	"encoding/json"
	"net/http"
	"time"
)

func (Uc *UserControl) HTTP_API_Standard_response(w http.ResponseWriter, r *http.Request, transaction API_Standard_response, KeepDataInDB bool) {
	//CollectPrometheusHTTPMetrics(r.URL.Path, transaction.Status, "", transaction.Elapsedtime)
	transaction.StatusDate = time.Now()
	transaction.Elapsedtime = (time.Since(transaction.ReceiveDate).Nanoseconds()) / 1000000
	//Uc.Write_StandardResponse_log(transaction, "", KeepDataInDB)
	w.Header().Set("Content-Type", "application/json")
	if transaction.Status == "successful" {
		w.WriteHeader(transaction.StatusCode)
	} else if transaction.Status == "failed" {
		w.WriteHeader(transaction.StatusCode)
	}
	json.NewEncoder(w).Encode(transaction)
}
