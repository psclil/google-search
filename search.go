package googlesearch

import (
	"context"
	"github.com/gocolly/colly/v2"
	neturl "net/url"
	"strconv"
	"strings"
)

// Result represents a single result from Google Search.
type Result struct {

	// Rank is the order number of the search result.
	Rank int `json:"rank"`

	// URL of result.
	URL string `json:"url"`

	// Title of result.
	Title string `json:"title"`

	// Description of the result.
	Description string `json:"description"`
}

// GoogleDomains represents localized Google homepages. The 2 letter country code is based on ISO 3166-1 alpha-2.
//
// See: https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2
var GoogleDomains = map[string]string{
	"us":  "https://www.google.com/search?",
	"ac":  "https://www.google.ac/search?",
	"ad":  "https://www.google.ad/search?",
	"ae":  "https://www.google.ae/search?",
	"af":  "https://www.google.com.af/search?",
	"ag":  "https://www.google.com.ag/search?",
	"ai":  "https://www.google.com.ai/search?",
	"al":  "https://www.google.al/search?",
	"am":  "https://www.google.am/search?",
	"ao":  "https://www.google.co.ao/search?",
	"ar":  "https://www.google.com.ar/search?",
	"as":  "https://www.google.as/search?",
	"at":  "https://www.google.at/search?",
	"au":  "https://www.google.com.au/search?",
	"az":  "https://www.google.az/search?",
	"ba":  "https://www.google.ba/search?",
	"bd":  "https://www.google.com.bd/search?",
	"be":  "https://www.google.be/search?",
	"bf":  "https://www.google.bf/search?",
	"bg":  "https://www.google.bg/search?",
	"bh":  "https://www.google.com.bh/search?",
	"bi":  "https://www.google.bi/search?",
	"bj":  "https://www.google.bj/search?",
	"bn":  "https://www.google.com.bn/search?",
	"bo":  "https://www.google.com.bo/search?",
	"br":  "https://www.google.com.br/search?",
	"bs":  "https://www.google.bs/search?",
	"bt":  "https://www.google.bt/search?",
	"bw":  "https://www.google.co.bw/search?",
	"by":  "https://www.google.by/search?",
	"bz":  "https://www.google.com.bz/search?",
	"ca":  "https://www.google.ca/search?",
	"kh":  "https://www.google.com.kh/search?",
	"cc":  "https://www.google.cc/search?",
	"cd":  "https://www.google.cd/search?",
	"cf":  "https://www.google.cf/search?",
	"cat": "https://www.google.cat/search?",
	"cg":  "https://www.google.cg/search?",
	"ch":  "https://www.google.ch/search?",
	"ci":  "https://www.google.ci/search?",
	"ck":  "https://www.google.co.ck/search?",
	"cl":  "https://www.google.cl/search?",
	"cm":  "https://www.google.cm/search?",
	"cn":  "https://www.google.cn/search?",
	"co":  "https://www.google.com.co/search?",
	"cr":  "https://www.google.co.cr/search?",
	"cu":  "https://www.google.com.cu/search?",
	"cv":  "https://www.google.cv/search?",
	"cy":  "https://www.google.com.cy/search?",
	"cz":  "https://www.google.cz/search?",
	"de":  "https://www.google.de/search?",
	"dj":  "https://www.google.dj/search?",
	"dk":  "https://www.google.dk/search?",
	"dm":  "https://www.google.dm/search?",
	"do":  "https://www.google.com.do/search?",
	"dz":  "https://www.google.dz/search?",
	"ec":  "https://www.google.com.ec/search?",
	"ee":  "https://www.google.ee/search?",
	"eg":  "https://www.google.com.eg/search?",
	"es":  "https://www.google.es/search?",
	"et":  "https://www.google.com.et/search?",
	"fi":  "https://www.google.fi/search?",
	"fj":  "https://www.google.com.fj/search?",
	"fm":  "https://www.google.fm/search?",
	"fr":  "https://www.google.fr/search?",
	"ga":  "https://www.google.ga/search?",
	"gb":  "https://www.google.co.uk/search?",
	"ge":  "https://www.google.ge/search?",
	"gf":  "https://www.google.gf/search?",
	"gg":  "https://www.google.gg/search?",
	"gh":  "https://www.google.com.gh/search?",
	"gi":  "https://www.google.com.gi/search?",
	"gl":  "https://www.google.gl/search?",
	"gm":  "https://www.google.gm/search?",
	"gp":  "https://www.google.gp/search?",
	"gr":  "https://www.google.gr/search?",
	"gt":  "https://www.google.com.gt/search?",
	"gy":  "https://www.google.gy/search?",
	"hk":  "https://www.google.com.hk/search?",
	"hn":  "https://www.google.hn/search?",
	"hr":  "https://www.google.hr/search?",
	"ht":  "https://www.google.ht/search?",
	"hu":  "https://www.google.hu/search?",
	"id":  "https://www.google.co.id/search?",
	"iq":  "https://www.google.iq/search?",
	"ie":  "https://www.google.ie/search?",
	"il":  "https://www.google.co.il/search?",
	"im":  "https://www.google.im/search?",
	"in":  "https://www.google.co.in/search?",
	"io":  "https://www.google.io/search?",
	"is":  "https://www.google.is/search?",
	"it":  "https://www.google.it/search?",
	"je":  "https://www.google.je/search?",
	"jm":  "https://www.google.com.jm/search?",
	"jo":  "https://www.google.jo/search?",
	"jp":  "https://www.google.co.jp/search?",
	"ke":  "https://www.google.co.ke/search?",
	"ki":  "https://www.google.ki/search?",
	"kg":  "https://www.google.kg/search?",
	"kr":  "https://www.google.co.kr/search?",
	"kw":  "https://www.google.com.kw/search?",
	"kz":  "https://www.google.kz/search?",
	"la":  "https://www.google.la/search?",
	"lb":  "https://www.google.com.lb/search?",
	"lc":  "https://www.google.com.lc/search?",
	"li":  "https://www.google.li/search?",
	"lk":  "https://www.google.lk/search?",
	"ls":  "https://www.google.co.ls/search?",
	"lt":  "https://www.google.lt/search?",
	"lu":  "https://www.google.lu/search?",
	"lv":  "https://www.google.lv/search?",
	"ly":  "https://www.google.com.ly/search?",
	"ma":  "https://www.google.co.ma/search?",
	"md":  "https://www.google.md/search?",
	"me":  "https://www.google.me/search?",
	"mg":  "https://www.google.mg/search?",
	"mk":  "https://www.google.mk/search?",
	"ml":  "https://www.google.ml/search?",
	"mm":  "https://www.google.com.mm/search?",
	"mn":  "https://www.google.mn/search?",
	"ms":  "https://www.google.ms/search?",
	"mt":  "https://www.google.com.mt/search?",
	"mu":  "https://www.google.mu/search?",
	"mv":  "https://www.google.mv/search?",
	"mw":  "https://www.google.mw/search?",
	"mx":  "https://www.google.com.mx/search?",
	"my":  "https://www.google.com.my/search?",
	"mz":  "https://www.google.co.mz/search?",
	"na":  "https://www.google.com.na/search?",
	"ne":  "https://www.google.ne/search?",
	"nf":  "https://www.google.com.nf/search?",
	"ng":  "https://www.google.com.ng/search?",
	"ni":  "https://www.google.com.ni/search?",
	"nl":  "https://www.google.nl/search?",
	"no":  "https://www.google.no/search?",
	"np":  "https://www.google.com.np/search?",
	"nr":  "https://www.google.nr/search?",
	"nu":  "https://www.google.nu/search?",
	"nz":  "https://www.google.co.nz/search?",
	"om":  "https://www.google.com.om/search?",
	"pa":  "https://www.google.com.pa/search?",
	"pe":  "https://www.google.com.pe/search?",
	"ph":  "https://www.google.com.ph/search?",
	"pk":  "https://www.google.com.pk/search?",
	"pl":  "https://www.google.pl/search?",
	"pg":  "https://www.google.com.pg/search?",
	"pn":  "https://www.google.pn/search?",
	"pr":  "https://www.google.com.pr/search?",
	"ps":  "https://www.google.ps/search?",
	"pt":  "https://www.google.pt/search?",
	"py":  "https://www.google.com.py/search?",
	"qa":  "https://www.google.com.qa/search?",
	"ro":  "https://www.google.ro/search?",
	"rs":  "https://www.google.rs/search?",
	"ru":  "https://www.google.ru/search?",
	"rw":  "https://www.google.rw/search?",
	"sa":  "https://www.google.com.sa/search?",
	"sb":  "https://www.google.com.sb/search?",
	"sc":  "https://www.google.sc/search?",
	"se":  "https://www.google.se/search?",
	"sg":  "https://www.google.com.sg/search?",
	"sh":  "https://www.google.sh/search?",
	"si":  "https://www.google.si/search?",
	"sk":  "https://www.google.sk/search?",
	"sl":  "https://www.google.com.sl/search?",
	"sn":  "https://www.google.sn/search?",
	"sm":  "https://www.google.sm/search?",
	"so":  "https://www.google.so/search?",
	"st":  "https://www.google.st/search?",
	"sv":  "https://www.google.com.sv/search?",
	"td":  "https://www.google.td/search?",
	"tg":  "https://www.google.tg/search?",
	"th":  "https://www.google.co.th/search?",
	"tj":  "https://www.google.com.tj/search?",
	"tk":  "https://www.google.tk/search?",
	"tl":  "https://www.google.tl/search?",
	"tm":  "https://www.google.tm/search?",
	"to":  "https://www.google.to/search?",
	"tn":  "https://www.google.tn/search?",
	"tr":  "https://www.google.com.tr/search?",
	"tt":  "https://www.google.tt/search?",
	"tw":  "https://www.google.com.tw/search?",
	"tz":  "https://www.google.co.tz/search?",
	"ua":  "https://www.google.com.ua/search?",
	"ug":  "https://www.google.co.ug/search?",
	"uk":  "https://www.google.co.uk/search?",
	"uy":  "https://www.google.com.uy/search?",
	"uz":  "https://www.google.co.uz/search?",
	"vc":  "https://www.google.com.vc/search?",
	"ve":  "https://www.google.co.ve/search?",
	"vg":  "https://www.google.vg/search?",
	"vi":  "https://www.google.co.vi/search?",
	"vn":  "https://www.google.com.vn/search?",
	"vu":  "https://www.google.vu/search?",
	"ws":  "https://www.google.ws/search?",
	"za":  "https://www.google.co.za/search?",
	"zm":  "https://www.google.co.zm/search?",
	"zw":  "https://www.google.co.zw/search?",
}

// SearchOptions modifies how the Search function behaves.
type SearchOptions struct {

	// CountryCode sets the ISO 3166-1 alpha-2 code of the localized Google Search homepage to use.
	// The default is "us", which will return results from https://www.google.com.
	CountryCode string

	// LanguageCode sets the language code.
	// Default: en
	LanguageCode string

	// Limit sets how many results to fetch (at maximum).
	Limit int

	// Start sets from what rank the new result set should return.
	Start int

	// UserAgent sets the UserAgent of the http request.
	// Default: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"
	UserAgent string

	// OverLimit searches for more results than that specified by Limit.
	// It then reduces the returned results to match Limit.
	OverLimit bool
}

// Search returns a list of search results from Google.
func Search(ctx context.Context, searchTerm string, opts ...SearchOptions) ([]Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := RateLimit.Wait(ctx); err != nil {
		return nil, err
	}

	c := colly.NewCollector(colly.MaxDepth(1))
	if len(opts) == 0 {
		opts = append(opts, SearchOptions{})
	}

	if opts[0].UserAgent == "" {
		c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"
	} else {
		c.UserAgent = opts[0].UserAgent
	}

	var lc string
	if opts[0].LanguageCode == "" {
		lc = "en"
	} else {
		lc = opts[0].LanguageCode
	}

	results := []Result{}
	var rErr error
	rank := 1

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		rErr = err
	})

	c.OnHTML("div.g", func(e *colly.HTMLElement) {

		sel := e.DOM

		for i := range sel.Nodes {
			if err := ctx.Err(); err != nil {
				rErr = err
				return
			}

			item := sel.Eq(i)

			rDiv := item.Find("div.rc")

			linkHref, _ := rDiv.Find("a").Attr("href")
			linkText := strings.TrimSpace(linkHref)
			titleText := strings.TrimSpace(rDiv.Find("h3").Text())

			descText := strings.TrimSpace(rDiv.Find("div > div > span > span").Text())

			if linkText != "" && linkText != "#" {
				result := Result{
					Rank:        rank,
					URL:         linkText,
					Title:       titleText,
					Description: descText,
				}
				results = append(results, result)
				rank += 1
			}
		}
	})

	limit := opts[0].Limit
	if opts[0].OverLimit {
		limit = int(float64(opts[0].Limit) * 1.5)
	}

	url := url(searchTerm, opts[0].CountryCode, lc, limit, opts[0].Start)
	c.Visit(url)

	if rErr != nil {
		if strings.Contains(rErr.Error(), "Too Many Requests") {
			return nil, ErrBlocked
		}
		return nil, rErr
	}

	// Reduce results to max limit
	if opts[0].Limit != 0 && len(results) > opts[0].Limit {
		return results[:opts[0].Limit], nil
	}

	return results, nil
}

func url(searchTerm string, countryCode string, languageCode string, limit int, start int) string {
	countryCode = strings.ToLower(countryCode)
	googleBase, ok := GoogleDomains[countryCode]
	if !ok {
		googleBase = GoogleDomains["us"]
	}
	queryString := neturl.Values{}
	queryString.Set("q", strings.TrimSpace(searchTerm))
	queryString.Set("hl", languageCode)
	if start != 0 {
		queryString.Set("start", strconv.Itoa(start))
	}
	if limit != 0 {
		queryString.Set("num", strconv.Itoa(limit))
	}

	return googleBase + queryString.Encode()
}
