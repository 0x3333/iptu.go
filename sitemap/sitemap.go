package sitemap

import (
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"
	"time"

	"bitbucket.org/terciofilho/iptu.go/db"
	"bitbucket.org/terciofilho/iptu.go/log"
)

const (
	maxRecordCount    = 2500000
	maxRecordsPerPage = 49999

	tplSitemaps = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">{{ $LastModify := .LastModify }}{{range $index, $element := .Docs}}
  <sitemap>
    <loc>http://consultaiptu.com.br/sitemaps/docs_{{ $index }}.xml.gz</loc>
    <lastmod>{{ $LastModify }}</lastmod>
  </sitemap>{{ else }}{{ end }}{{range $index, $element := .Nomes}}
  <sitemap>
    <loc>http://consultaiptu.com.br/sitemaps/nomes_{{ $index }}.xml.gz</loc>
    <lastmod>{{ $LastModify }}</lastmod>
  </sitemap>{{ else }}{{ end }}
</sitemapindex>
`

	tplDocs = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">{{ $LastModify := .LastModify }}{{ range .Docs }}
  <url>
    <loc>http://consultaiptu.com.br/s/{{ . }}</loc>
    <lastmod>{{ $LastModify }}</lastmod>
    <changefreq>monthly</changefreq>
  </url>{{ end }}
</urlset>
`

	tplNomes = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">{{ $LastModify := .LastModify }}{{ range .Nomes }}
  <url>
    <loc>http://consultaiptu.com.br/s/{{ . }}</loc>
    <lastmod>{{ $LastModify }}</lastmod>
    <changefreq>monthly</changefreq>
  </url>{{end}}
</urlset>
`
)

//Sitemap represents a sitemap struct to template
type Sitemap struct {
	LastModify string
	Docs       *[][]string
	Nomes      *[][]string
}

//Docs represents a Doc struct to template
type Docs struct {
	LastModify string
	Docs       *[]string
}

//Nomes represents a Nome struct to template
type Nomes struct {
	LastModify string
	Nomes      *[]string
}

//Generate generates the sitemap files
func Generate() {
	// Limpa o diretório
	absPath, _ := filepath.Abs("web/sitemaps/")
	os.RemoveAll(absPath)
	os.Mkdir(absPath, os.ModePerm)

	t := time.Now()
	lastModify := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	var wr *os.File

	//
	// Documentos
	//

	log.Info.Println("Starting Documentos...")
	rows, err := db.Instance.Query("SELECT `doc_contribuinte_1` as doc_contribuinte FROM iptu WHERE `doc_contribuinte_1` != \"\" GROUP BY `doc_contribuinte_1` UNION DISTINCT SELECT `doc_contribuinte_2` as doc_contribuinte FROM iptu WHERE `visivel` = 1 AND `doc_contribuinte_2` != \"\" GROUP BY `doc_contribuinte_2`")
	if err != nil {
		log.Error.Printf("Failed to query database: %s", err.Error())
	}
	defer rows.Close()
	docs := make([][]string, 0, maxRecordCount)
	innerDocs := make([]string, 0, maxRecordsPerPage)
	index := 0
	for rows.Next() {
		index++
		var doc string
		rows.Scan(&doc)
		innerDocs = append(innerDocs, doc)
		if index%maxRecordsPerPage == 0 {
			log.Info.Printf("index: %d\n", index)
			docs = append(docs, innerDocs)
			innerDocs = make([]string, 0, maxRecordsPerPage)
		}
	}
	templateDocs, err := template.New("tplDocs").Parse(tplDocs)
	if err != nil {
		log.Error.Println(err.Error())
	}
	for index, innerDoc := range docs {
		absPath, _ = filepath.Abs("web/sitemaps/docs_" + strconv.Itoa(index) + ".xml.gz")
		wr, err = os.Create(absPath)
		wrGzip := gzip.NewWriter(wr)
		if err != nil {
			log.Error.Println(err.Error())
			return
		}
		err = templateDocs.Execute(wrGzip, Docs{
			LastModify: lastModify,
			Docs:       &innerDoc,
		})
		if err != nil {
			log.Error.Println(err.Error())
		}
		wrGzip.Close()
		wr.Close()
	}

	//
	// Nomes
	//

	log.Info.Println("Starting Nomes...")
	rows, err = db.Instance.Query("SELECT `nome_contribuinte_1` as nome_contribuinte FROM iptu WHERE `nome_contribuinte_1` != \"\" GROUP BY `nome_contribuinte_1` UNION DISTINCT SELECT `nome_contribuinte_2` as nome_contribuinte FROM iptu WHERE `visivel` = 1 AND `nome_contribuinte_2` != \"\" GROUP BY `nome_contribuinte_2`")
	if err != nil {
		log.Error.Printf("Failed to query database: %s", err.Error())
	}
	defer rows.Close()
	nomes := make([][]string, 0, maxRecordCount)
	innerNomes := make([]string, 0, maxRecordsPerPage)
	index = 0
	regexWords := regexp.MustCompile(`[^\w ]+`)
	regexWords2 := regexp.MustCompile(`[_]+`)
	regexSpaces := regexp.MustCompile(` +`)
	for rows.Next() {
		index++
		var nome string
		rows.Scan(&nome)
		nome = regexWords.ReplaceAllString(nome, "")
		nome = regexWords2.ReplaceAllString(nome, " ")
		nome = regexSpaces.ReplaceAllString(nome, "-")
		innerNomes = append(innerNomes, nome)
		if index%maxRecordsPerPage == 0 {
			log.Info.Printf("index: %d\n", index)
			nomes = append(nomes, innerNomes)
			innerNomes = make([]string, 0, maxRecordsPerPage)
		}
	}
	templateNomes, err := template.New("tplNomes").Parse(tplNomes)
	if err != nil {
		log.Error.Println(err.Error())
	}
	for index, innerNomes := range nomes {
		absPath, _ = filepath.Abs("web/sitemaps/nomes_" + strconv.Itoa(index) + ".xml.gz")
		wr, err = os.Create(absPath)
		wrGzip := gzip.NewWriter(wr)
		if err != nil {
			log.Error.Println(err.Error())
			return
		}
		err = templateNomes.Execute(wrGzip, Nomes{
			LastModify: lastModify,
			Nomes:      &innerNomes,
		})
		if err != nil {
			log.Error.Println(err.Error())
		}
		wrGzip.Close()
		wr.Close()
	}

	//
	// Sitemap
	//

	log.Info.Println("Starting Sitemap...")
	templateSitemap, err := template.New("tplSitemaps").Parse(tplSitemaps)
	if err != nil {
		log.Error.Println(err.Error())
	}
	absPath, _ = filepath.Abs("web/sitemaps/sitemap.xml")
	wr, err = os.Create(absPath)
	if err != nil {
		log.Error.Println(err.Error())
		return
	}
	defer wr.Close()
	err = templateSitemap.Execute(wr, Sitemap{
		LastModify: lastModify,
		Docs:       &docs,
		Nomes:      &nomes,
	})
	if err != nil {
		log.Error.Println(err.Error())
	}
	log.Info.Println("Done!")
}
