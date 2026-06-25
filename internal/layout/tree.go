package layout

import "tsugu-mcp/family"

// 描画用の家系ツリー(primary=血統側、spouse=配偶者、children=次世代)
type treeNode struct {
	primary  person
	spouse   *person
	children []*treeNode
}

// カード生成に必要な人物情報
type person struct {
	relationship string
	name         string
	address      string
	honseki      string // 本籍(被相続人のみ)
	birth        family.Date
	death        *family.Date
	applicant    bool
	outcome      family.Outcome
	isDecedent   bool
}

// Documentを1本の家系ツリーへ変換(尊属あれば尊属夫婦を根とし被相続人・兄弟姉妹を子)
func toTree(doc family.Document) *treeNode {
	dec := &treeNode{
		primary:  decedentPerson(doc.Decedent),
		children: convertNodes(doc.Children),
	}
	if doc.Spouse != nil {
		sp := toPerson(*doc.Spouse)
		dec.spouse = &sp
	}
	if len(doc.Ascendants) == 0 {
		return dec
	}

	root := &treeNode{primary: toPerson(doc.Ascendants[0])}
	if len(doc.Ascendants) >= 2 {
		sp := toPerson(doc.Ascendants[1])
		root.spouse = &sp
	}
	root.children = append([]*treeNode{dec}, convertNodes(doc.Siblings)...)
	return root
}

func convertNodes(nodes []*family.Node) []*treeNode {
	var out []*treeNode
	for _, n := range nodes {
		tn := &treeNode{primary: toPerson(n.Person), children: convertNodes(n.Descendants)}
		if n.Spouse != nil {
			sp := toPerson(*n.Spouse)
			tn.spouse = &sp
		}
		out = append(out, tn)
	}
	return out
}

func toPerson(p family.Person) person {
	return person{
		relationship: p.Relationship,
		name:         p.Name,
		address:      p.Address,
		birth:        p.BirthDate,
		death:        p.DeathDate,
		applicant:    p.Applicant,
		outcome:      p.Outcome,
	}
}

func decedentPerson(d family.Decedent) person {
	p := person{
		relationship: "被",
		name:         d.Name,
		address:      d.LastAddress,
		honseki:      d.RegisteredDomicile,
		birth:        d.BirthDate,
		isDecedent:   true,
	}
	if !d.DeathDate.IsZero() {
		dd := d.DeathDate
		p.death = &dd
	}
	return p
}
