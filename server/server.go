package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type iptu struct {
	NumeroContribuinte     string
	DocContribuinte1       string
	NomeContribuinte1      string
	DocContribuinte2       string
	NomeContribuinte2      string
	NomeLogradouroImovel   string
	NumeroImovel           string
	ComplementoImovel      string
	BairroImovel           string
	ReferenciaImovel       string
	CepImovel              string
	FracaoIdeal            float32
	AreaTerreno            int
	AreaConstruida         int
	AreaOcupada            int
	ValorM2Terreno         float32
	ValorM2Construcao      float32
	AnoConstrucaoCorrigido string
	QuantidadePavimentos   int
	TestadaCalculo         float32
	TipoUsoImovel          string
	TipoPadraoConstrucao   string
	TipoTerreno            string
	FatorObsolescencia     float32
}

func (i iptu) String() string {
	return fmt.Sprintf("numeroContribuinte: %s, docContribuinte1: %s, nomeContribuinte1: %s, docContribuinte2: %s, nomeContribuinte2: %s, nomeLogradouroImovel: %s, numeroImovel: %s, complementoImovel: %s, bairroImovel: %s, referenciaImovel: %s, cepImovel: %s, fracaoIdeal: %f, areaTerreno: %d, areaConstruida: %d, areaOcupada: %d, valorM2Terreno: %f, valorM2Construcao: %f, anoConstrucaoCorrigido: %s, quantidadePavimentos: %d, testadaCalculo: %f, tipoUsoImovel: %s, tipoPadraoConstrucao: %s, tipoTerreno: %s, fatorObsolescencia: %f",
		i.NumeroContribuinte,
		i.DocContribuinte1,
		i.NomeContribuinte1,
		i.DocContribuinte2,
		i.NomeContribuinte2,
		i.NomeLogradouroImovel,
		i.NumeroImovel,
		i.ComplementoImovel,
		i.BairroImovel,
		i.ReferenciaImovel,
		i.CepImovel,
		i.FracaoIdeal,
		i.AreaTerreno,
		i.AreaConstruida,
		i.AreaOcupada,
		i.ValorM2Terreno,
		i.ValorM2Construcao,
		i.AnoConstrucaoCorrigido,
		i.QuantidadePavimentos,
		i.TestadaCalculo,
		i.TipoUsoImovel,
		i.TipoPadraoConstrucao,
		i.TipoTerreno,
		i.FatorObsolescencia)
}

var (
	db *sql.DB
)

// Server starts the webserver to handle the requests from the UI
func Server(innerDb *sql.DB) {
	db = innerDb

	defer db.Close()

	handleStatic()
	handleAPI()

	log.Println("WebServer started...")
	http.ListenAndServe(":8080", nil)
}

func handleStatic() {
	fs := http.FileServer(http.Dir("web"))
	http.Handle("/", fs)
}

func handleAPI() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing: %s", err.Error())
		}
		termos := r.FormValue("termos")
		termos = strings.Replace(termos, ".", " ", -1)
		termos = strings.Replace(termos, ",", " ", -1)
		termos = strings.Replace(termos, "-", " ", -1)
		termos = strings.Replace(termos, "_", " ", -1)
		termos = strings.TrimSpace(termos)
		termosSlice := strings.Split(termos, " ")

		termosFT := ""
		for _, value := range termosSlice {
			// Only search for strings with at least 3 chars
			if len(value) >= 3 {
				termosFT += "+" + value + " "
			}
		}
		if len(termosFT) == 0 {
			http.Error(w, "Os Termos devem conter ao menos 3 caracteres cada.", http.StatusBadRequest)
		}
		rows, err := db.Query("SELECT numero_contribuinte,doc_contribuinte_1,nome_contribuinte_1,doc_contribuinte_2,nome_contribuinte_2,nome_logradouro_imovel,numero_imovel,complemento_imovel,bairro_imovel,referencia_imovel,cep_imovel,fracao_ideal,area_terreno,area_construida,area_ocupada,valor_m2_terreno,valor_m2_construcao,ano_construcao_corrigido,quantidade_pavimentos,testada_calculo,tipo_uso_imovel,tipo_padrao_construcao,tipo_terreno,fator_obsolescencia FROM `iptu` WHERE (MATCH(`nome_contribuinte_1`,`nome_contribuinte_2`,`nome_logradouro_imovel`,`referencia_imovel`) AGAINST (? IN BOOLEAN MODE)) UNION SELECT numero_contribuinte,doc_contribuinte_1,nome_contribuinte_1,doc_contribuinte_2,nome_contribuinte_2,nome_logradouro_imovel,numero_imovel,complemento_imovel,bairro_imovel,referencia_imovel,cep_imovel,fracao_ideal,area_terreno,area_construida,area_ocupada,valor_m2_terreno,valor_m2_construcao,ano_construcao_corrigido,quantidade_pavimentos,testada_calculo,tipo_uso_imovel,tipo_padrao_construcao,tipo_terreno,fator_obsolescencia FROM `iptu` WHERE `doc_contribuinte_1` = ? OR `doc_contribuinte_2` = ? LIMIT 100", termosFT, termos, termos)
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		var iptus []iptu
		for rows.Next() {
			row := iptu{}
			rows.Scan(&row.NumeroContribuinte,
				&row.DocContribuinte1,
				&row.NomeContribuinte1,
				&row.DocContribuinte2,
				&row.NomeContribuinte2,
				&row.NomeLogradouroImovel,
				&row.NumeroImovel,
				&row.ComplementoImovel,
				&row.BairroImovel,
				&row.ReferenciaImovel,
				&row.CepImovel,
				&row.FracaoIdeal,
				&row.AreaTerreno,
				&row.AreaConstruida,
				&row.AreaOcupada,
				&row.ValorM2Terreno,
				&row.ValorM2Construcao,
				&row.AnoConstrucaoCorrigido,
				&row.QuantidadePavimentos,
				&row.TestadaCalculo,
				&row.TipoUsoImovel,
				&row.TipoPadraoConstrucao,
				&row.TipoTerreno,
				&row.FatorObsolescencia)
			iptus = append(iptus, row)
		}
		if iptus == nil {
			w.Write([]byte("[]"))
		} else {
			bytes, err := json.Marshal(&iptus)
			if err != nil {
				panic(err.Error())
			}
			println(string(bytes))
			w.Write(bytes)
		}
	})
}
