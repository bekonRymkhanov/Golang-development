<<<<<<< HEAD
package recipes

import "errors"

var (
	NotFoundErr = errors.New("not found")
)

type MemStore struct {
	list map[string]Recipe
}

func NewMemStore() *MemStore {
	list := make(map[string]Recipe)
	newRecipe := Recipe{
		Name: "Spaghetti Bolognese",
		Ingredients: []Ingredient{
			{Name: "Spaghetti"},
			{Name: "Ground beef"},
			{Name: "Tomato sauce"},
			{Name: "Onions"},
			// Add more ingredients as needed
		},
	}

	// Adding the recipe to the 'list' map
	list["1"] = newRecipe
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(name string, recipe Recipe) error {
	m.list[name] = recipe
	return nil
}

func (m MemStore) Get(name string) (Recipe, error) {

	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return Recipe{}, NotFoundErr
}

func (m MemStore) List() (map[string]Recipe, error) {
	return m.list, nil
}

func (m MemStore) Update(name string, recipe Recipe) error {

	if _, ok := m.list[name]; ok {
		m.list[name] = recipe
		return nil
	}

	return NotFoundErr
}

func (m MemStore) Remove(name string) error {
	delete(m.list, name)
	return nil
}
=======
package recipes

import "errors"

var (
	NotFoundErr = errors.New("not found")
)

type MemStore struct {
	list map[string]Recipe
}

func NewMemStore() *MemStore {
	list := make(map[string]Recipe)
	newRecipe := Recipe{
		Name: "Spaghetti Bolognese",
		Ingredients: []Ingredient{
			{Name: "Spaghetti"},
			{Name: "Ground beef"},
			{Name: "Tomato sauce"},
			{Name: "Onions"},
			// Add more ingredients as needed
		},
	}

	// Adding the recipe to the 'list' map
	list["1"] = newRecipe
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(name string, recipe Recipe) error {
	m.list[name] = recipe
	return nil
}

func (m MemStore) Get(name string) (Recipe, error) {

	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return Recipe{}, NotFoundErr
}

func (m MemStore) List() (map[string]Recipe, error) {
	return m.list, nil
}

func (m MemStore) Update(name string, recipe Recipe) error {

	if _, ok := m.list[name]; ok {
		m.list[name] = recipe
		return nil
	}

	return NotFoundErr
}

func (m MemStore) Remove(name string) error {
	delete(m.list, name)
	return nil
}
>>>>>>> 84013f4f47cab449aae9ee53ab1ed2ab42c22d66