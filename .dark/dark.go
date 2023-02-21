package dark

// "private struct that implements a public interface"
// Conceptually, a value of an interface type, or interface value, has 2 components:
// a concrete type (type descriptor) and a value of that type.
// The descriptor is a pointer to virtual table and the interface value is the pointer
// to the instance of the concrete type that implements the interface.
//
// interface struct {
//		*vtable of functions, listed in the Darker interface declaration
//		*dark struct, as created by NewDark()
//	}

type Darker interface {
	ListVariantGroups() []string
	ListVariants(string) []string
	NewBaize(string) (*Baize, error)
	LoadBaize(string) (*Baize, error)
}

// dark holds the state for the current game/baize in play. It is NOT exported
// from this package, making it opaque to the client.
// All access to this struct is through the Darker interface
type dark struct {
	baize *Baize
}

func NewDark() Darker {
	return &dark{}
}
