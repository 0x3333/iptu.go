package render

import (
	"html/template"
	"io"

	"bitbucket.org/terciofilho/iptu.go/api"
	"bitbucket.org/terciofilho/iptu.go/log"
)

const tpl = `
<!DOCTYPE html>
<html lang="pt">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta http-equiv="x-ua-compatible" content="ie=edge">
    <meta name="description" content="Consulta de contribuintes do IPTU da Cidade de São Paulo - SP">

    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.0.0-alpha.2/css/bootstrap.min.css">
    <title>Consulta de Contribuintes do IPTU - São Paulo - SP</title>

    <style>
				body {
						margin-bottom: 60px;
				}
				.footer {
					  position: fixed;
					  bottom: 0;
					  width: 100%;
					  background-color: #f5f5f5;
				}
				.row.message {
            display: none;
        }
    </style>
</head>

<body>
    <div class="container">
        <div class="row header m-t-3">
            <div class="col-sm-12">
                <div class="jumbotron text-xs-center m-b-0">
                    <h1 class="display-3">Consulta IPTU</h1>
                    <p class="lead">Aqui você pode consultar pelos Contribuintes de IPTU da Cidade de São Paulo</p>
                </div>
            </div>
        </div>
        <div class="row header m-y-1">
            <div class="col-sm-offset-2 col-sm-8">
								<h5><span class="label label-default hidden-lg-up">Nome/CNPJ/CPF/Logradouro:</span></h5>
                <div class="input-group">
                    <span class="input-group-addon hidden-md-down">Nome/CNPJ/CPF/Logradouro:</span>
                    <input id="termos" type="text" class="form-control" placeholder="Digite os termos da busca">
                    <span class="input-group-btn"><button id="pesquisar" class="btn btn-success" type="button">Pesquisar!</button></span>
                </div>
            </div>
        </div>
        <div class="row error {{if ne .HasError true}}message {{end}}header m-t-3">
            <div class="col-sm-offset-3 col-sm-6">
                <div class="alert alert-danger" role="alert">
                    <strong>Ops!</strong> Aconteceu um erro ao processar a sua solicitação.
                </div>
            </div>
        </div>{{if eq .IsLimited true}}
        <div class="row invalid header m-t-1">
            <div class="col-sm-offset-3 col-sm-6">
                <div class="alert alert-warning" role="alert">
                    <strong>Ops!</strong> Os resultados foram limitados, seja mais específico na pesquisa.
                </div>
            </div>
        </div>{{end}}
        <div class="row invalid {{if ne .InvalidRequest true}}message{{end}} header m-t-3">
            <div class="col-sm-offset-3 col-sm-6">
                <div class="alert alert-warning" role="alert">
                    <strong>Ops!</strong> Digite termos com mais de 3 caracteres.
                </div>
            </div>
        </div>
    {{if ne .Index true}}
				{{if .IPTUs}}
        {{range .IPTUs}}
        <div class="row result header">
            <div class="col-sm-12">
                <div class="card" itemscope itemtype="http://schema.org/Person">
                    <div class="card-header"><strong itemprop="name">{{.NomeContribuinte1}}</strong> (<strong>{{.TipoContribuinte1}}</strong> <span itemprop="taxID">{{.DocContribuinte1}}</span>){{if .NomeContribuinte2}} - <strong itemprop="name">{{.NomeContribuinte2}}</strong> {{if .TipoContribuinte2}}(<strong>{{.TipoContribuinte2}}</strong> <span itemprop="taxID">{{.DocContribuinte2}}</span>){{end}}{{end}}</div>
                    <div class="card-block">
                        <div class="row">
														<div class="col-sm-4"><strong>N. Contribuinte:</strong> {{.NumeroContribuinte}}</div>
                            <div class="col-sm-8"><strong>Endereço: </strong><span itemprop="address">{{.NomeLogradouroImovel}}, {{.NumeroImovel}}{{if .ComplementoImovel}} - {{.ComplementoImovel}}{{end}}, {{.BairroImovel}}</span> - <strong>CEP:</strong> {{.CepImovel}}
														<a href="http://maps.google.com/?q={{ .URLMaps }}" target="_blank"><img src="https://upload.wikimedia.org/wikipedia/en/1/19/Google_Maps_Icon.png" width="32" height="32"></a>
														</div>
                        </div>
                        <div class="row">
                            <div class="col-sm-4"><strong>Ref.:</strong> {{.ReferenciaImovel}}</div>
                            <div class="col-sm-4"><strong>Fração Ideal:</strong> {{.FracaoIdeal}}</div>
                            <div class="col-sm-4"><strong>Ano Construção:</strong> {{.AnoConstrucaoCorrigido}}</div>
                        </div>
                        <div class="row">
                            <div class="col-sm-4"><strong>Área Terreno:</strong> {{.AreaTerreno}}</div>
                            <div class="col-sm-4"><strong>Área Construida:</strong> {{.AreaConstruida}}</div>
                            <div class="col-sm-4"><strong>Área Ocupada:</strong> {{.AreaOcupada}}</div>
                        </div>
                        <div class="row">
                            <div class="col-sm-4"><strong>Valor m<sup>2</sup> Terreno:</strong> {{.ValorM2Terreno}}</div>
                            <div class="col-sm-4"><strong>Valor m<sup>2</sup> Construção:</strong> {{.ValorM2Construcao}}</div>
														<div class="col-sm-4"><strong>Pavimentos:</strong> {{.QuantidadePavimentos}}</div>
                        </div>
                        <div class="row">
                            <div class="col-sm-12">{{.TipoUsoImovel}} - {{.TipoPadraoConstrucao}} - {{.TipoTerreno}}</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
				{{end}}
        {{else}}
        <div class="row noresults header m-t-3">
            <div class="col-sm-offset-3 col-sm-6">
                <div class="alert alert-warning" role="alert">
                    <strong>Ops!</strong> Nenhum resultado encontrado!
                </div>
            </div>
        </div>
        {{end}}
    {{end}}
    </div>

    <footer class="footer">
        <div class="container">
            <small><span class="text-muted">Dados disponibilizados pela Prefeitura de São Paulo conforme Decreto Nº 56.932, de 13 de abril de 2016 - <a href="http://diariooficial.imprensaoficial.com.br/doflash/prototipo/2016/Abril/14/cidade/pdf/pg_0001.pdf" target="_blank">Diário Oficial</a></span></small>
        </div>
    </footer>

    <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/2.2.4/jquery.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/tether/1.3.2/js/tether.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.0.0-alpha.2/js/bootstrap.min.js"></script>
    <script>
    (function(i, s, o, g, r, a, m) {
        i['GoogleAnalyticsObject'] = r;
        i[r] = i[r] || function() {
            (i[r].q = i[r].q || []).push(arguments)
        }, i[r].l = 1 * new Date();
        a = s.createElement(o),
            m = s.getElementsByTagName(o)[0];
        a.async = 1;
        a.src = g;
        m.parentNode.insertBefore(a, m)
    })(window, document, 'script', 'https://www.google-analytics.com/analytics.js', 'ga');
    ga('create', 'UA-80338230-1', 'auto');
    ga('send', 'pageview');

		function convertToURL(text) {
				return text.toLowerCase().replace(/[^\w ]+/g, '').replace(/ +/g, '-');
		}

    $(function() {
        $("#termos").keypress(function(e) {
            if (e.which == 13) {
                $("#pesquisar").click();
            }
        });
        $("#pesquisar").click(function() {
            window.location.href = "/s/" + convertToURL($("#termos").val());
        });

        $("#termos").focus();
    });
    </script>
</body>

</html>
`

type tplRender struct {
	IPTUs          *[]api.IPTU
	Index          bool
	InvalidRequest bool
	IsLimited      bool
	HasError       bool
}

// Render renders a Template based on the IPTUs
func Render(IPTUs *[]api.IPTU, index bool, invalidRequest bool, hasError bool, wr io.Writer) {
	template, err := template.New("tpl").Parse(tpl)
	if err != nil {
		log.Error.Println(err.Error())
	}
	isLimited := false
	if IPTUs != nil {
		isLimited = len(*IPTUs) >= api.LimitSize
	}
	err = template.Execute(wr, tplRender{
		IPTUs:          IPTUs,
		Index:          index,
		InvalidRequest: invalidRequest,
		IsLimited:      isLimited,
		HasError:       hasError,
	})
	if err != nil {
		log.Error.Println(err.Error())
	}
}
