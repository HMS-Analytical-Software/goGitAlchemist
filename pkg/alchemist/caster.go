package alchemist

// caster defines the ability to cast a spell (i.e. execute a command)
// from the Formula. It also allows checking if the spell is valid.
type caster interface {
	validate() error
	cast(assistant, Options) error
}
