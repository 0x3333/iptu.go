package api

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/0x3333/iptu.go/db"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// IPTU Representa um Contribuinte de IPTU
type IPTU struct {
	NumeroContribuinte     string
	TipoContribuinte1      string
	DocContribuinte1       string
	NomeContribuinte1      string
	TipoContribuinte2      string
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
	URLMaps                string
}

func (i IPTU) String() string {
	return fmt.Sprintf("numeroContribuinte: %s, tipoContribuinte1: %s, docContribuinte1: %s, nomeContribuinte1: %s, tipoContribuinte2: %s, docContribuinte2: %s, nomeContribuinte2: %s, nomeLogradouroImovel: %s, numeroImovel: %s, complementoImovel: %s, bairroImovel: %s, referenciaImovel: %s, cepImovel: %s, fracaoIdeal: %f, areaTerreno: %d, areaConstruida: %d, areaOcupada: %d, valorM2Terreno: %f, valorM2Construcao: %f, anoConstrucaoCorrigido: %s, quantidadePavimentos: %d, testadaCalculo: %f, tipoUsoImovel: %s, tipoPadraoConstrucao: %s, tipoTerreno: %s, fatorObsolescencia: %f", i.NumeroContribuinte, i.TipoContribuinte1, i.DocContribuinte1, i.NomeContribuinte1, i.TipoContribuinte2, i.DocContribuinte2, i.NomeContribuinte2, i.NomeLogradouroImovel, i.NumeroImovel, i.ComplementoImovel, i.BairroImovel, i.ReferenciaImovel, i.CepImovel, i.FracaoIdeal, i.AreaTerreno, i.AreaConstruida, i.AreaOcupada, i.ValorM2Terreno, i.ValorM2Construcao, i.AnoConstrucaoCorrigido, i.QuantidadePavimentos, i.TestadaCalculo, i.TipoUsoImovel, i.TipoPadraoConstrucao, i.TipoTerreno, i.FatorObsolescencia)
}

// LimitSize is the size used in the LIMIT SQL query
const LimitSize = 150

var (
	regex1    = regexp.MustCompile(`(\d{1})[,\.\-\/ ]+(\d{1})`)
	regex2    = regexp.MustCompile(`[,\.\-\/]+`)
	regex3    = regexp.MustCompile(`\b[^ ]{1,2}\b`)
	regex4    = regexp.MustCompile(`[ ]+`)
	regex5    = regexp.MustCompile(`\s([^ ]+)`)
	regexCNPJ = regexp.MustCompile(`(..)(...)(...)(....)(..)`)
	regexCPF  = regexp.MustCompile(`(...)(...)(...)(..)`)
)

// RequestError represents a Request error
type RequestError struct {
	Invalid  bool
	HasError bool
	Message  string
}

// HandleRequest trata uma requisição de consulta na base de IPTU
func HandleRequest(termos string) (*[]IPTU, *RequestError) {
	if len(termos) == 0 {
		return nil, &RequestError{
			Invalid: true,
			Message: "Os Termos devem conter ao menos 3 caracteres cada.",
		}
	}
	// Remove os acentos
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	termos, _, _ = transform.String(t, termos)

	// Remove os .-/ que estão entre números
	termos = regex1.ReplaceAllString(termos, "$1$2")
	termos = regex2.ReplaceAllString(termos, " ")
	termos = regex3.ReplaceAllString(termos, " ")
	termos = regex4.ReplaceAllString(termos, " ")
	termos = strings.TrimSpace(termos)
	termosFT := "+" + regex5.ReplaceAllString(termos, " +$1")

	// log.Info.Printf("Termos: '%s' - TermosFT: '%s'", termos, termosFT)

	if len(termos) == 0 {
		return nil, &RequestError{
			Invalid: true,
			Message: "Os Termos devem conter ao menos 3 caracteres cada.",
		}
	}

	rows, err := db.Instance.Query(fmt.Sprintf(`
	(
		 SELECT numero_contribuinte,
            tipo_contribuinte_1,
            doc_contribuinte_1,
            nome_contribuinte_1,
            tipo_contribuinte_2,
            doc_contribuinte_2,
            nome_contribuinte_2,
            nome_logradouro_imovel,
            numero_imovel,
            complemento_imovel,
            bairro_imovel,
            referencia_imovel,
            cep_imovel,
            fracao_ideal,
            area_terreno,
            area_construida,
            area_ocupada,
            valor_m2_terreno,
            valor_m2_construcao,
            ano_construcao_corrigido,
            quantidade_pavimentos,
            testada_calculo,
            tipo_uso_imovel,
            tipo_padrao_construcao,
            tipo_terreno,
            fator_obsolescencia
   FROM iptu
   WHERE visivel = 1 AND (
			 MATCH(nome_contribuinte_1,nome_contribuinte_2,nome_logradouro_imovel,referencia_imovel) AGAINST (? IN BOOLEAN MODE)
		) LIMIT %d
	)
UNION
  (SELECT numero_contribuinte,
          tipo_contribuinte_1,
          doc_contribuinte_1,
          nome_contribuinte_1,
          tipo_contribuinte_2,
          doc_contribuinte_2,
          nome_contribuinte_2,
          nome_logradouro_imovel,
          numero_imovel,
          complemento_imovel,
          bairro_imovel,
          referencia_imovel,
          cep_imovel,
          fracao_ideal,
          area_terreno,
          area_construida,
          area_ocupada,
          valor_m2_terreno,
          valor_m2_construcao,
          ano_construcao_corrigido,
          quantidade_pavimentos,
          testada_calculo,
          tipo_uso_imovel,
          tipo_padrao_construcao,
          tipo_terreno,
          fator_obsolescencia
   FROM iptu
   WHERE visivel = 1
     AND doc_contribuinte_1 = ?
     OR doc_contribuinte_2 = ? LIMIT %d)
ORDER BY nome_contribuinte_1,
         nome_contribuinte_2
	`, LimitSize, LimitSize), termosFT, termos, termos)

	if err != nil {
		return nil, &RequestError{
			HasError: true,
			Message:  err.Error(),
		}
	}
	defer rows.Close()
	var result []IPTU
	for rows.Next() {
		row := IPTU{}
		rows.Scan(&row.NumeroContribuinte,
			&row.TipoContribuinte1,
			&row.DocContribuinte1,
			&row.NomeContribuinte1,
			&row.TipoContribuinte2,
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
		if len(row.DocContribuinte1) == 14 {
			row.DocContribuinte1 = regexCNPJ.ReplaceAllString(row.DocContribuinte1, "$1.$2.$3/$4-$5")
		} else if len(row.DocContribuinte1) == 11 {
			row.DocContribuinte1 = regexCPF.ReplaceAllString(row.DocContribuinte1, "$1.$2.$3-$4")
		}
		if len(row.DocContribuinte2) == 14 {
			row.DocContribuinte2 = regexCNPJ.ReplaceAllString(row.DocContribuinte2, "$1.$2.$3/$4-$5")
		} else if len(row.DocContribuinte2) == 11 {
			row.DocContribuinte2 = regexCPF.ReplaceAllString(row.DocContribuinte2, "$1.$2.$3-$4")
		}
		row.URLMaps = fmt.Sprintf("%s, %s, São Paulo - SP", row.NomeLogradouroImovel, row.NumeroImovel)
		result = append(result, row)
	}
	if result == nil {
		return nil, nil
	}
	return &result, nil
}
