package main

func replaceBody(p *post) string {
	html, err := p.RawBody().Html()
	if err != nil {
		panic(err)
	}
	return html
}
