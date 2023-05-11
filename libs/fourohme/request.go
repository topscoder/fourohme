package fourohme

type Request struct {
	Verb    string
	Url     string
	Headers []Header
}
