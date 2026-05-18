package irnd

type Random interface {
	PublicID() string
	SessionToken() string
}
