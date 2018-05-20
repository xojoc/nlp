/*  Copyright (C) 2018 Alexandru Cojocaru

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. */

package wikidata

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/binary"
	"github.com/davecgh/go-spew/spew"
	jsoniter "github.com/json-iterator/go"

	"xojoc.pw/must"
)

// type Value interface{}

type Coordinate struct {
	Latitude   float64
	Longitude  float64
	Globe      string
	Prescision float64
}

// fixme: use geo uri
func (c *Coordinate) String() string {
	if c == nil {
		return ""
	}
	if *c == (Coordinate{}) {
		return ""
	}
	return fmt.Sprintf("%v,%v", c.Latitude, c.Longitude)
}

type Quantity struct {
	Amount float64
	// fixme: 1 (unit-less) eq Q0
	Unit       ID
	LowerBound float64
	UpperBound float64
}

type Value struct {
	EntityID   ID
	Coordinate Coordinate
	Text       string
	Time       time.Time
	Quantity   Quantity
}

// type Text string

// type EntityID = ID

/*
type ValueType int

func init() {
	var t Text
	gob.Register(t)
	var id ID
	gob.Register(id)
	var tim time.Time
	gob.Register(tim)
	var coor Coordinate
	gob.Register(coor)
}

const (
	ValueTypeEntityID = iota
	ValueTypeCoordinate
	ValueTypeText
	ValueTypeTime
)


func GetValueType(v Value) ValueType {
	switch v.(type) {
	case EntityID:
		return ValueTypeEntityID
	case Coordinate:
		return ValueTypeCoordinate
	case Text:
		return ValueTypeText
	case time.Time:
		return ValueTypeTime
	default:
		panic("GetValueType: unknown value type")
	}
}
*/

type Claim struct {
	Value Value
}

type Entity struct {
	ID ID
	//	Labels       map[string]string
	Label string
	//	Aliases      map[string][]string
	Aliases []string
	//	Descriptions map[string]string
	Description string
	Claims      map[ID][]Claim
	SiteLinks   map[string]string
}

func (e *Entity) HasIstance(value string) bool {
	return e.HasStatement("P31", value)
}

func (e *Entity) HasStatement(statement, value string) bool {
	if e == nil {
		return false
	}
	for _, v := range e.Claims[NewID(statement)] {
		if v.Value.EntityID == NewID(value) {
			return true
		}
	}
	return false
}

type jsonEntity struct {
	ID           string
	Labels       map[string]struct{ Value string }
	Aliases      map[string][]struct{ Value string }
	Descriptions map[string]struct{ Value string }
	Claims       map[string]jsoniter.RawMessage
	SiteLinks    map[string]struct{ Title string }
}

type ID uint64

func NewID(s string) ID {
	u, err := strconv.ParseUint(s[1:], 10, 64)
	must.OK(err)
	if s[0] == 'P' {
		u |= 1 << 63
	}
	return ID(u)
}

func (i ID) Property() bool {
	if i&(1<<63) == 1<<63 {
		return true
	}
	return false
}

func (i ID) Item() bool {
	return !i.Property()
}

func (i ID) String() string {
	if i.Item() {
		return fmt.Sprintf("Q%d", i)
	}
	return fmt.Sprintf("P%d", i&^(1<<63))
}

func parseTime(t string) time.Time {
	// t := time.Time{}
	// e.g. +1994-01-01T00:00:00Z
	sign := +1
	if t[0] == '-' {
		sign = -1
	}
	year, err := strconv.Atoi(strings.Split(t[1:], "-")[0])
	must.OK(err)
	month, err := strconv.Atoi(strings.Split(t[1:], "-")[1])
	must.OK(err)
	day, err := strconv.Atoi(strings.Split(t[1:], "-")[2][:2])
	must.OK(err)
	return time.Date(sign*year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

/*
func init() {
	fmt.Println(parseTime("+1994-01-02T00:00:00Z"))
	fmt.Println(parseTime("-13798000000-00-00T00:00:00Z"))
}
*/

func jsonToEntity(e *jsonEntity) (*Entity, error) {
	ent := &Entity{}
	/*
		ent.Labels = make(map[string]string)
		for lang, label := range e.Labels {
			ent.Labels[lang] = label.Value
		}
	*/
	if l, ok := e.Labels["en"]; ok {
		ent.Label = l.Value
	}

	/*
		ent.Aliases = make(map[string][]string)
		for lang, aliases := range e.Aliases {
			for _, alias := range aliases {
				ent.Aliases[lang] = append(ent.Aliases[lang], alias.Value)
			}
		}
	*/
	ent.Aliases = make([]string, 0, len(e.Aliases["en"]))
	for _, alias := range e.Aliases["en"] {
		ent.Aliases = append(ent.Aliases, alias.Value)
	}

	/*
		ent.Descriptions = make(map[string]string)
		for lang, description := range e.Descriptions {
			ent.Descriptions[lang] = description.Value
		}
	*/
	if d, ok := e.Descriptions["en"]; ok {
		ent.Description = d.Value
	}

	ent.Claims = make(map[ID][]Claim, len(e.Claims))

	/*
		ent.SiteLinks = make(map[string]string)
		dbe.SiteLinks["wikipedia"] = e.SiteLinks.EnWiki.Title
		dbe.SiteLinks["wikiquote"] = e.SiteLinks.EnWikiquote.Title
	*/

	for pid, v := range e.Claims {
		type claimValue struct {
			MainSnak struct {
				DataValue struct {
					Value jsoniter.RawMessage
					Type  string
				}
			}
		}
		var cvs []claimValue
		err := jsoniter.Unmarshal(v, &cvs)
		if err != nil {
			return nil, err
		}
		for _, cv := range cvs {
			var c Claim
			v := cv.MainSnak.DataValue.Value
			switch cv.MainSnak.DataValue.Type {
			case "string":
				var text string
				err := jsoniter.Unmarshal(v, &text)
				if err != nil {
					return nil, err
				}
				//				c.Value = Text(text)
				c.Value.Text = text
			case "wikibase-entityid":
				var eid struct{ ID string }
				err := jsoniter.Unmarshal(v, &eid)
				if err != nil {
					return nil, err
				}
				// c.Value = NewID(eid.ID)
				c.Value.EntityID = NewID(eid.ID)
			case "globecoordinate":
				var coordinate Coordinate
				err := jsoniter.Unmarshal(v, &coordinate)
				if err != nil {
					return nil, err
				}
				c.Value.Coordinate = coordinate
			case "time":
				var t struct{ Time string }
				err := jsoniter.Unmarshal(v, &t)
				if err != nil {
					return nil, err
				}
				c.Value.Time = parseTime(t.Time)
			case "quantity":
				var q struct {
					Amount     string
					Unit       string
					LowerBound string
					UpperBound string
				}
				err := json.Unmarshal(v, &q)
				if err != nil {
					return nil, err
				}
				i := strings.LastIndex(q.Unit, "Q")
				var u ID
				if i != -1 {
					u = NewID(q.Unit[i:])
				}
				a, err := strconv.ParseFloat(q.Amount, 64)
				if err != nil {
					return nil, err
				}
				var low float64
				if q.LowerBound != "" {
					low, err = strconv.ParseFloat(q.LowerBound, 64)
					if err != nil {
						return nil, err
					}
				}
				var upp float64
				if q.UpperBound != "" {
					upp, err = strconv.ParseFloat(q.UpperBound, 64)
					if err != nil {
						return nil, err
					}
				}
				c.Value.Quantity = Quantity{
					Amount:     a,
					Unit:       u,
					LowerBound: low,
					UpperBound: upp,
				}
			default:
				//	log.Println("unknown type: ", cv.MainSnak.DataValue.Type)
				continue
			}
			ent.Claims[NewID(pid)] = append(ent.Claims[NewID(pid)], c)
		}
	}

	ent.SiteLinks = make(map[string]string, len(e.SiteLinks))
	for site, title := range e.SiteLinks {
		ent.SiteLinks[site] = title.Title
	}

	ent.ID = NewID(e.ID)

	return ent, nil
}

func Entities(r io.Reader, minLength int) (chan *Entity, chan error) {
	entities := make(chan *Entity)
	cerr := make(chan error, 1)
	lastLine := []byte("]\n")
	buf := bufio.NewReaderSize(r, 100*1024*1024)
	buf.ReadSlice('\n')

	go func() {
		for {
			line, err := buf.ReadSlice('\n')
			if err != nil {
				if err != io.EOF {
					cerr <- err
				}
				break
			}

			if bytes.Equal(line, lastLine) {
				break
			}

			if line[len(line)-2] == ',' {
				line = line[:len(line)-2]
			}

			if len(line) < minLength {
				continue
			}

			t := jsonEntity{}
			err = jsoniter.Unmarshal(line, &t)
			if err != nil {
				cerr <- err
				break
			}

			e, err := jsonToEntity(&t)
			if err != nil {
				cerr <- err
				break
			}
			if e != nil {
				entities <- e
			}
		}

		close(entities)
		close(cerr)
	}()

	return entities, cerr
}

func BuildDB(wiki io.Reader, db io.Writer, filter func(*Entity) *Entity, minLength int) error {
	w := gob.NewEncoder(db)
	entities, cerr := Entities(wiki, minLength)

	for e := range entities {
		e = filter(e)
		if e == nil {
			continue
		}
		err := w.Encode(e)
		if err != io.EOF {
			must.OK(err)
		}
	}

	return <-cerr
}

func OpenDB(r io.Reader) (chan *Entity, chan error) {
	g := gob.NewDecoder(bufio.NewReader(r))
	entities := make(chan *Entity)
	cerr := make(chan error, 1)

	go func() {
		for {
			var e Entity
			err := g.Decode(&e)
			if err != nil {
				if err != io.EOF {
					log.Println("opendb: ", err)
					//					cerr <- err
				}
				break
			}
			entities <- &e
		}
		close(entities)
		close(cerr)
	}()

	return entities, cerr
}
func BuildDB2(wiki io.Reader, db io.Writer, index io.Writer, filter func(*Entity) *Entity, minLength int) error {
	dbb := bufio.NewWriter(db)
	indexb := bufio.NewWriter(index)

	i := uint64(0)

	entities, cerr := Entities(wiki, minLength)

	for e := range entities {
		e = filter(e)
		if e == nil {
			continue
		}
		buf, err := binary.Marshal(e)
		must.OK(err)

		l := uint64(len(buf))
		n, err := dbb.Write(buf)
		must.OK(err)
		if uint64(n) != l {
			panic("n != l: " + fmt.Sprint(n) + ", " + fmt.Sprint(l))
		}

		_, err = indexb.WriteString(fmt.Sprintf("%d %d\n", e.ID, toIndex(i, l)))
		must.OK(err)
		i += l
	}

	must.OK(<-cerr)
	must.OK(dbb.Flush())
	must.OK(indexb.Flush())
	return nil
}

type DB struct {
	Index        map[ID]uint64
	IndexOrdered []ID
	Entities     io.ReaderAt
}

func toIndex(off, l uint64) uint64 {
	return (off & 0xffffffffffff) | (l << 48)
}
func fromIndex(i uint64) (uint64, uint64) {
	off := i & 0xffffffffffff
	l := i >> 48
	return off, l
}
func toInt(a []byte) uint64 {
	i, err := strconv.ParseUint(string(a), 10, 64)
	must.OK(err)
	return i
}
func OpenDB2(r io.ReaderAt, index io.Reader) (*DB, error) {
	d := &DB{}
	d.Index = make(map[ID]uint64)
	d.Entities = r

	buf := bufio.NewReader(index)
	for {
		l, err := buf.ReadSlice('\n')
		if err == io.EOF {
			break
		}
		must.OK(err)
		fs := bytes.Fields(l)

		d.Index[ID(toInt(fs[0]))] = toInt(fs[1])
		d.IndexOrdered = append(d.IndexOrdered, ID(toInt(fs[0])))
	}
	return d, nil
}

func (db *DB) Get(id string) *Entity {
	return db.GetID(NewID(id))
}

// fixme: lock
func (db *DB) GetID(id ID) *Entity {
	if _, ok := db.Index[id]; !ok {
		return nil
	}
	off, l := fromIndex(db.Index[id])
	b := make([]byte, l)
	n, err := db.Entities.ReadAt(b, int64(off))
	if err != io.EOF {
		must.OK(err)
	}
	if uint64(n) != l {
		panic("n != l: " + fmt.Sprint(n) + ", " + fmt.Sprint(l))
	}
	if err == io.EOF {
		log.Printf("%v: eof: %v\n", id, err)
	}
	e := Entity{}
	//	dec := gob.NewDecoder(bytes.NewReader(b))
	//	err = dec.Decode(&e)
	err = binary.Unmarshal(b, &e)
	if err != nil {
		if err == io.EOF {
			if e.ID == NewID("Q0") {
				fmt.Println("Q0: ", e)
				return nil
			}
		} else {
			log.Printf("%v: binary decode: %v\n", id, err)
			spew.Println(e)
			return nil
		}
	}

	return &e
}
func (db *DB) MustGetID(id ID) *Entity {
	return db.GetID(id)
}

// doc: https://www.wikidata.org/wiki/Help:Basic_membership_properties
func (db *DB) HasClass(id, class ID) bool {
	e := db.MustGetID(id)
	return db.hasClass(e, class)
}
func (db *DB) hasClass(e *Entity, class ID) bool {
	if e == nil {
		return false
	}
	for _, c := range append(e.Claims[NewID("P31")], e.Claims[NewID("P279")]...) {
		if c.Value.EntityID == class {
			return true
		}
		cve := db.MustGetID(c.Value.EntityID)
		ok := db.hasClass(cve, class)
		if ok {
			return true
		}
	}
	return false
}
