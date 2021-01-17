package game

type IFeature interface {
	Join(name string) (Character, error)
	List() ([]Character, error)
}
