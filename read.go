package factual

import(
	"fmt"
	"errors"
	"strings"
	"strconv"
	"encoding/json"
)

type Hours map[string][][]string

// Factual gives its hours as a seperatly encoded entity
func (hours *Hours) UnmarshalJSON(data []byte) (err error){
	// intermediary to avoid an infinite loop
	i := map[string][][]string{}
	s, err := strconv.Unquote(string(data))
	err = json.Unmarshal([]byte(s), &i)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid hours in JSON: %s (%s)", string(data), err))
	}
	*hours = Hours(i)
	return
}

func (q *ReadQuery) Params() map[string]string{
	filters, err := json.Marshal(q.filter)
	geo, err := json.Marshal(q.geo)

	if(err != nil){
		panic("Invalid filters")
	}

	sel := strings.Join(q._select, ",")
	includeCount := strconv.FormatBool(q.includeCount)
	sort := strings.Join(q.sort, ",")

	return map[string]string{
		"filters": string(filters),
		"q": q.q,
		"select": sel,
		"geo": string(geo),
		"threshold": q.threshold,
		"offset": string(q.offset),
		"limit": string(q.limit),
		"include_count": includeCount,
		"sort": sort,
	}
}

func (q *ReadQuery) Url() (u string){
	u = q.Table.url()
	if q.id != ""{
		u += "/" + q.id
	}
	return
}

func (q *ReadQuery) Iter() *Iter{
	return &Iter{q.Factual, q}
}

type ReadQuery struct{
	Factual Factual
	Table *Table
	id string
	q string
	filter F
	geo F
	_select []string
	threshold string
	offset int
	limit int
	includeCount bool
	// Incomplete: sort is "blending json" type
	sort []string
}

func (q *ReadQuery) Id(id string) *ReadQuery{
	q.id = id
	return q
}

func (q *ReadQuery) Filter(f F) *ReadQuery{
	q.filter = f
	return q
}

func (q *ReadQuery) Select(fields ...string) *ReadQuery{
	q._select = fields
	return q
}

func (q *ReadQuery) Threshold(t string) *ReadQuery{
	options := map[string]bool{
		"confident": true, 
		"default": true, 
		"comprehensive": true,
	}
	_, ok := options[t]
	if(ok == false){
		panic("Not a valid threshold")
	}
	q.threshold = t
	return q
}

func (q *ReadQuery) Offset(i int) *ReadQuery{
	q.offset = i
	return q
}

func (q *ReadQuery) Limit(i int) *ReadQuery{
	q.limit = i
	return q
}

func (q *ReadQuery) IncludeCount(i bool) *ReadQuery{
	q.includeCount = i
	return q
}

func (q *ReadQuery) Sort(fields ...string) *ReadQuery{
	q.sort = fields
	return q
}