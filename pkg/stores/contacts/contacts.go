package contacts

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
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

func (s *Store) Count() int {
	return len(s.c)
}

func (s *Store) Page(p int) Contacts {
	l := len(s.c)
	start := (p - 1) * 10
	end := start + 10

	return s.c[start:min(end, l)]
}

func (s *Store) Search(text string) Contacts {
	cs := make(Contacts, 0, 10)
	for _, c := range s.c {
		switch {
		case strings.Contains(c.First, text):
			cs = append(cs, c)
		case strings.Contains(c.Last, text):
			cs = append(cs, c)
		case strings.Contains(c.Email, text):
			cs = append(cs, c)
		case strings.Contains(c.Phone, text):
			cs = append(cs, c)
		}
	}
	return cs
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
			return errors.New("Email must be unique")
		}
	}

	return nil
}

func (s *Store) Save(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")

	if err := e.Encode(s.c); err != nil {
		return err
	}

	return nil
}
