package main

import (
	"fmt"
	"log"
	"sync"
)

type Setter interface {
	AddDetail(d *Detail) error
	String(d *Detail)
	SetColor()
}

type Detail struct {
	Color        string
	Saccharinity int
	Size         uint8
}

type Banana struct {
	Name string
}

type Apple struct {
	Name string
}

func (Banana) AddDetail(d *Detail) error {
	return d.SetSize(100)
}

func (Apple) AddDetail(d *Detail) error {
	return d.SetSize(200)
}

func (d *Detail) SetSize(size uint8) error {
	d.Size = size
	return nil
}

func (b *Banana) String(d *Detail) {
	fmt.Println("Name: ", b.Name)
	fmt.Println("Color: ", d.Color)
	fmt.Println("Saccharinity: ", d.Saccharinity)
	fmt.Println("Size: ", d.Size)
}

func (a *Apple) String(d *Detail) {
	fmt.Println("Name: ", a.Name)
	fmt.Println("Color: ", d.Color)
	fmt.Println("Saccharinity: ", d.Saccharinity)
	fmt.Println("Size: ", d.Size)
}

func (b *Banana) SetColor() { b.Name = "Yellow" }
func (a *Apple) SetColor()  { a.Name = "Red" }

func Color(s Setter) {
	s.SetColor()
}

func main() {
	banana := new(Banana)
	apple := new(Apple)

	Color(banana)
	Color(apple)

	var b Setter = banana
	var a Setter = apple

	_, err := Build(b, a)
	if err != nil {
		log.Fatal(err)
	}
}

func Build(s ...Setter) (*Detail, error) {
	d := new(Detail)
	d.Saccharinity = 20
	return d, d.build(s...)
}

func (d *Detail) build(s ...Setter) error {
	var once sync.Once
	for _, v := range s {
		if err := v.AddDetail(d); err != nil {
			return err
		}
		v.String(d)
		once.Do(func() { fmt.Println() })
	}

	return nil
}
