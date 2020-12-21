package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile)
	foods, err := ReadInput(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	allergens, err := FindAllergens(foods)
	if err != nil {
		log.Fatal(err)
	}
	safe := FindSafeIngredients(foods, allergens)
	fmt.Println("Number of uses of safe ingredients:", CountSafeIngredients(foods, safe))

	ingredients := make([]string, 0, len(allergens))
	for i := range allergens {
		ingredients = append(ingredients, i)
	}
	sort.Slice(ingredients, func(i, j int) bool {
		return allergens[ingredients[i]] < allergens[ingredients[j]]
	})
	fmt.Println(strings.Join(ingredients, ","))
}

type Food struct {
	ingredients StringSet
	allergens   StringSet
}

func ReadInput(r io.Reader) ([]Food, error) {
	var foods []Food

	s := bufio.NewScanner(r)
	for s.Scan() {
		l := s.Text()
		i := strings.Index(l, "(contains ")
		if i < 0 {
			i = len(l)
		}
		ingredients := strings.Split(strings.TrimSpace(l[:i]), " ")
		for i, s := range ingredients {
			ingredients[i] = strings.TrimSpace(s)
		}
		var allergens []string
		if i < len(l) {
			allergens = strings.Split(strings.TrimPrefix(l[i:len(l)-1], "(contains "), ", ")
			for i, s := range allergens {
				allergens[i] = strings.TrimSpace(s)
			}
		}
		f := Food{
			allergens:   make(StringSet),
			ingredients: make(StringSet),
		}
		f.allergens.Add(allergens...)
		f.ingredients.Add(ingredients...)
		foods = append(foods, f)
	}
	return foods, s.Err()
}

func CountSafeIngredients(fs []Food, safe StringSet) int {
	var n int
	for _, f := range fs {
		for i := range f.ingredients {
			if safe.Has(i) {
				n++
			}
		}
	}
	return n
}

func FindSafeIngredients(fs []Food, allergens map[string]string) StringSet {
	out := make(StringSet)
	for _, f := range fs {
		for i := range f.ingredients {
			if _, ok := allergens[i]; !ok {
				out.Add(i)
			}
		}
	}
	return out
}

func FindAllergens(fs []Food) (map[string]string, error) {
	candidates := make(map[string]StringSet)
	for _, f := range fs {
		for a := range f.allergens {
			if candidates[a].Empty() {
				candidates[a] = f.ingredients
			} else {
				candidates[a] = candidates[a].Intersect(f.ingredients)
			}
		}
	}
	allergens := make(map[string]string)
allergenLoop:
	for len(candidates) > 0 {
		for a, c := range candidates {
			i, ok := c.Singleton()
			if !ok {
				continue
			}
			for _, c := range candidates {
				c.Remove(i)
			}
			allergens[i] = a
			delete(candidates, a)
			continue allergenLoop
		}
		return nil, errors.New("can't determine ingredients for allergens")
	}
	return allergens, nil
}

type StringSet map[string]bool

func (s StringSet) Intersect(s2 StringSet) StringSet {
	out := make(StringSet)
	for e := range s {
		if s2[e] {
			out.Add(e)
		}
	}
	return out
}

func (s StringSet) Add(es ...string) {
	for _, e := range es {
		s[e] = true
	}
}

func (s StringSet) Remove(es ...string) {
	for _, e := range es {
		delete(s, e)
	}
}

func (s StringSet) Has(e string) bool {
	return s[e]
}

func (s StringSet) String() string {
	var es []string
	for e := range s {
		es = append(es, strconv.Quote(e))
	}
	return "{" + strings.Join(es, ", ") + "}"
}

// Singleton returns (e, true) if s represents the set {e}.
func (s StringSet) Singleton() (e string, ok bool) {
	for e := range s {
		return e, len(s) == 1
	}
	return "", false
}

func (s StringSet) Empty() bool {
	return len(s) == 0
}
