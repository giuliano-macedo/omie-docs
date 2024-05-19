package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"testing"

	"github.com/giuliano-macedo/omie-docs/internal/core"
	"github.com/giuliano-macedo/omie-docs/test/mocks"
	"github.com/stretchr/testify/assert"
)

const seed = 42

func TestGetPages(t *testing.T) {
	home, totalEntities := generateRandomHomePage()

	var fetcherMock = new(mocks.FetcherMock)
	for i := 0; i < totalEntities; i++ {
		fetcherMock.On("GetPage", pageUrl(i)).Return(generatePageHtml(fmt.Sprint("page", i)), nil)
	}

	pages, err := GetPages(fetcherMock, home, 16)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, len(pages), totalEntities)
	for i, page := range pages {
		if !assert.Equal(t, page.Name, fmt.Sprint("page", i)) {
			return
		}
		fetcherMock.MethodCalled("GetPage", pageUrl(i))
	}
}

func TestGetPagesError(t *testing.T) {
	home, totalEntities := generateRandomHomePage()

	getPageError := errors.New("get page error")
	var fetcherMock = new(mocks.FetcherMock)
	for i := 0; i < totalEntities; i++ {
		var resp *mocks.ReadCloseMock
		var err error
		if i == totalEntities/2 {
			resp, err = nil, getPageError
		} else {
			resp, err = generatePageHtml(fmt.Sprint("page", i)), nil
		}
		fetcherMock.On("GetPage", pageUrl(i)).Return(resp, err)
	}

	_, err := GetPages(fetcherMock, home, 16)
	if !assert.Error(t, err) {
		return
	}

	assert.EqualError(t, err, fmt.Sprintf("failed getting %v: %v", pageUrl(totalEntities/2), getPageError))

}

func generatePageHtml(pageName string) *mocks.ReadCloseMock {
	templateString := `
<html>
<body>
<div class="jumbotron banner bg-azul">
	<div class="col-12">
	<h1>{{.PageName}}</h1>
	<div>
		<div>
		<h4></h4>
		<a
			class="highlight"
			href="https://example.com/base/version1/directory/to/file/"
		>
			<pre><code>https://example.com/base/version1/directory/to/file/</code></pre>
		</a>
		</div>
	</div>
	</div>
</div>
</body>
</html>
	`
	tmpl, err := template.New("PageHtml").Parse(templateString)
	if err != nil {
		panic(err)
	}
	var buff bytes.Buffer
	ctx := struct{ PageName string }{pageName}

	if err := tmpl.Execute(&buff, ctx); err != nil {
		panic(err)
	}

	return mocks.NewReadCloseMock(buff.String())
}

func pageUrl(index int) string {
	return fmt.Sprint("https://example.com/page", index)
}

func rndNumber(rng *rand.Rand) int {
	return rng.Intn(8) + 2
}

func generateRandomEntities(rng *rand.Rand, initialIndex int) []core.Entity {
	entities := make([]core.Entity, 0, rndNumber(rng))
	i := initialIndex
	for j := 0; j < cap(entities); j++ {
		var entity core.Entity
		entity.Url = pageUrl(i)
		entities = append(entities, entity)
		i++
	}
	return entities
}
func generateRandomHomePage() (core.HomePage, int) {
	rng := rand.New(rand.NewSource(seed))

	home := core.HomePage{}
	totalEntities := 0
	home.Features = make([]core.Feature, 0, rndNumber(rng))
	for i := 0; i < cap(home.Features); i++ {
		var feature core.Feature

		feature.MainEntities = generateRandomEntities(rng, totalEntities)
		totalEntities += len(feature.MainEntities)
		feature.AuxiliaryEntites = generateRandomEntities(rng, totalEntities)
		totalEntities += len(feature.AuxiliaryEntites)

		home.Features = append(home.Features, feature)
	}
	return home, totalEntities
}
