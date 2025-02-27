package contacts

import (
	"encoding/json"
	"fmt"
	"os"
)

type Contact struct {
	Id    int    `json:"id"`
	First string `json:"first" validate:"required"`
	Last  string `json:"last" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type Contacts []Contact

type Store struct {
	c Contacts
}

func NewStore(file string) (*Store, error) {
	db, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	s := &Store{}
	if err := json.Unmarshal(db, &s.c); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) Page(p int) Contacts {
	fmt.Println("p", p)
	if p <= 1 {
		return s.c[:10]
	}

	l := len(s.c)
	start := p * 10
	if start > l {
		return s.c[l-10:]
	}

	end := (p + 1) * 10

	fmt.Println(min(end, l))

	return s.c[start:min(end, l)]
}

func (s *Store) Create(c Contact) error {
	if err := s.Validate(c); err != nil {
		return err
	}
	c.Id = len(s.c) + 1
	s.c = append(s.c, c)
	return nil
}

func (s *Store) Get(id int) Contact {
	for _, c := range s.c {
		if id == c.Id {
			return c
		}
	}

	return Contact{}
}

func (s *Store) Update(c Contact) error {
	if err := s.Validate(c); err != nil {
		return err
	}

	for i, v := range s.c {
		if c.Id != v.Id {
			continue
		}
		s.c[i] = c
		break
	}

	return nil
}

func (s *Store) Delete(id int) bool {
	for i, c := range s.c {
		if c.Id != int(id) {
			continue
		}
		copy(s.c[i:], s.c[i+1:])
		s.c[len(s.c)-1] = Contact{}
		s.c = s.c[:len(s.c)-1]
		return true
	}

	return false
}

func (s *Store) Validate(c Contact) error {
	for _, v := range s.c {
		if c.Email == v.Email {
			return fmt.Errorf("Email must be unique")
		}
	}

	return nil
}
