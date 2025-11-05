package hf

// RepositoryType represents the kind of Hugging Face repository.
type RepositoryType int

const (
	RepoTypeModel RepositoryType = iota
	RepoTypeDataset
)
